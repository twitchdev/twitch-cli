// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package moderation

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/twitchdev/twitch-cli/internal/database"
	"github.com/twitchdev/twitch-cli/internal/mock_api/authentication"
	"github.com/twitchdev/twitch-cli/internal/mock_api/mock_errors"
	"github.com/twitchdev/twitch-cli/internal/models"
)

var bansMethodsSupported = map[string]bool{
	http.MethodGet:    false,
	http.MethodPost:   true,
	http.MethodDelete: true,
	http.MethodPatch:  false,
	http.MethodPut:    false,
}

var bansScopesByMethod = map[string][]string{
	http.MethodGet:    {},
	http.MethodPost:   {"moderator:manage:banned_users"},
	http.MethodDelete: {"moderator:manage:banned_users"},
	http.MethodPatch:  {},
	http.MethodPut:    {},
}

type PostBansRequestBodyData struct {
	UserID    string `json:"user_id"`
	Duration  int    `json:"duration"`
	Reason    string `json:"reason"`
	Duplicate bool   `json:"-"`
}

type PostBansRequestBody struct {
	Data PostBansRequestBodyData `json:"data"`
}

type PostBansResponseBodyData struct {
	BroadcasterID string  `json:"broadcaster_id"`
	ModeratorID   string  `json:"moderator_Id"`
	UserID        string  `json:"user_id"`
	CreatedAt     string  `json:"created_at"`
	EndTime       *string `json:"end_time"`
}

type Bans struct{}

func (e Bans) Path() string { return "/moderation/bans" }

func (e Bans) GetRequiredScopes(method string) []string {
	return bansScopesByMethod[method]
}

func (e Bans) ValidMethod(method string) bool {
	return bansMethodsSupported[method]
}

func (e Bans) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	db = r.Context().Value("db").(database.CLIDatabase)

	switch r.Method {
	case http.MethodPost:
		postBans(w, r)
		break
	case http.MethodDelete:
		deleteBans(w, r)
		break
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func postBans(w http.ResponseWriter, r *http.Request) {
	userCtx := r.Context().Value("auth").(authentication.UserAuthentication)
	if !userCtx.MatchesModeratorIDParam(r) {
		mock_errors.WriteUnauthorized(w, "Moderator ID does not match token.")
		return
	}

	broadcasterID := r.URL.Query().Get("broadcaster_id")
	if broadcasterID == "" {
		mock_errors.WriteBadRequest(w, "Missing required parameter broadcaster_id")
		return
	}

	moderatorID := r.URL.Query().Get("moderator_id")
	if moderatorID == "" {
		mock_errors.WriteBadRequest(w, "Missing required parameter moderator_id")
		return
	}

	// Check if broadcaster exists
	broadcaster, err := db.NewQuery(r, 100).GetUser(database.User{ID: broadcasterID})
	if err != nil {
		mock_errors.WriteServerError(w, "error fetching broadcaster")
		return
	}
	if broadcaster.ID == "" {
		mock_errors.WriteUnauthorized(w, "The user specified in parameter moderator_id is not one of the broadcaster's moderators")
		return
	}

	var moderatorList []database.Moderator

	// Verify user is a moderator or is the broadcaster
	isModerator := false
	if broadcasterID == moderatorID {
		isModerator = true
	} else {
		moderatorListDbr, err := db.NewQuery(r, 1000).GetModeratorsForBroadcaster(broadcasterID)
		if err != nil {
			mock_errors.WriteServerError(w, err.Error())
			return
		}

		moderatorList = moderatorListDbr.Data.([]database.Moderator)

		for _, mod := range moderatorList {
			if mod.UserID == moderatorID {
				isModerator = true
			}
		}
	}
	if !isModerator {
		mock_errors.WriteUnauthorized(w, "The user specified in parameter moderator_id is not one of the broadcaster's moderators")
		return
	}

	var body PostBansRequestBody
	err = json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		mock_errors.WriteBadRequest(w, "Body unable to be parsed")
		return
	}

	if body.Data.UserID == "" {
		mock_errors.WriteBadRequest(w, "Missing required field user_id in request body")
		return
	}

	if len(body.Data.Reason) > 500 {
		mock_errors.WriteBadRequest(w, "Ban reason must be less than 500 characters")
		return
	}

	if body.Data.Duration > 1209600 || body.Data.Duration < 0 {
		mock_errors.WriteBadRequest(w, "Ban duration must be non-zero and less than 1,209,600 seconds")
		return
	}

	// Check if user is broadcaster
	if body.Data.UserID == broadcasterID {
		mock_errors.WriteBadRequest(w, fmt.Sprintf("User %v cannot be banned or timed out", broadcasterID))
		return
	}

	// Check if user is a moderator. Bypass this check if the broadcaster is banning someone.
	if moderatorID != broadcasterID {
		isModerator = false
		for _, mod := range moderatorList {
			if mod.UserID == body.Data.UserID {
				isModerator = true
			}
		}
		if isModerator {
			mock_errors.WriteBadRequest(w, fmt.Sprintf("You may not ban or time out user %v", body.Data.UserID))
			return
		}
	}

	// Check if user exists
	foundUser, err := db.NewQuery(r, 100).GetUser(database.User{ID: body.Data.UserID})
	if foundUser.ID == "" {
		mock_errors.WriteBadRequest(w, fmt.Sprintf("User %v doesn't exist", body.Data.UserID))
		return
	}

	// Check if user is already banned
	dbr, err := db.NewQuery(r, 100).GetBans(database.UserRequestParams{BroadcasterID: broadcasterID, UserID: body.Data.UserID})
	if err != nil {
		mock_errors.WriteServerError(w, "error fetching bans")
		return
	}
	if len(dbr.Data.([]database.Ban)) != 0 {
		mock_errors.WriteBadRequest(w, fmt.Sprintf("User %v is already banned", body.Data.UserID))
		return
	}

	err = db.NewQuery(r, 100).InsertBan(database.UserRequestParams{
		UserID:        body.Data.UserID,
		BroadcasterID: broadcasterID,
	})
	if err != nil {
		mock_errors.WriteServerError(w, "error inserting ban")
		return
	}

	timeNow := time.Now().UTC().Format(time.RFC3339)
	var timeLater *string
	if body.Data.Duration != 0 {
		time := time.Now().UTC().Add(time.Duration(body.Data.Duration) * time.Second).Format(time.RFC3339)
		timeLater = &time
	}

	response := PostBansResponseBodyData{
		BroadcasterID: broadcasterID,
		ModeratorID:   moderatorID,
		UserID:        body.Data.UserID,
		CreatedAt:     timeNow,
		EndTime:       timeLater,
	}

	bytes, _ := json.Marshal(models.APIResponse{Data: []PostBansResponseBodyData{response}})
	w.Write(bytes)
}

func deleteBans(w http.ResponseWriter, r *http.Request) {
	userCtx := r.Context().Value("auth").(authentication.UserAuthentication)
	if !userCtx.MatchesModeratorIDParam(r) {
		mock_errors.WriteUnauthorized(w, "Moderator ID does not match token.")
		return
	}

	broadcasterID := r.URL.Query().Get("broadcaster_id")
	if broadcasterID == "" {
		mock_errors.WriteBadRequest(w, "Missing required parameter broadcaster_id")
		return
	}

	moderatorID := r.URL.Query().Get("moderator_id")
	if moderatorID == "" {
		mock_errors.WriteBadRequest(w, "Missing required parameter moderator_id")
		return
	}

	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		mock_errors.WriteBadRequest(w, "Missing required parameter user_id")
		return
	}

	// Check if broadcaster exists
	broadcaster, err := db.NewQuery(r, 100).GetUser(database.User{ID: broadcasterID})
	if err != nil {
		mock_errors.WriteServerError(w, "error fetching broadcaster")
		return
	}
	if broadcaster.ID == "" {
		mock_errors.WriteUnauthorized(w, "The user specified in parameter moderator_id is not one of the broadcaster's moderators")
		return
	}

	var moderatorList []database.Moderator

	// Verify user is a moderator or is the broadcaster
	isModerator := false
	if broadcasterID == moderatorID {
		isModerator = true
	} else {
		moderatorListDbr, err := db.NewQuery(r, 1000).GetModeratorsForBroadcaster(broadcasterID)
		if err != nil {
			mock_errors.WriteServerError(w, err.Error())
			return
		}

		moderatorList = moderatorListDbr.Data.([]database.Moderator)

		for _, mod := range moderatorList {
			if mod.UserID == moderatorID {
				isModerator = true
			}
		}
	}
	if !isModerator {
		mock_errors.WriteUnauthorized(w, "The user specified in parameter moderator_id is not one of the broadcaster's moderators")
		return
	}

	// Check if user exists
	bannedUser, err := db.NewQuery(r, 100).GetUser(database.User{ID: userID})
	if bannedUser.ID == "" {
		mock_errors.WriteBadRequest(w, fmt.Sprintf("User %v doesn't exist", userID))
		return
	}

	// Check if user is banned
	dbr, err := db.NewQuery(r, 100).GetBans(database.UserRequestParams{BroadcasterID: broadcasterID, UserID: userID})
	if err != nil {
		mock_errors.WriteServerError(w, "error fetching bans")
		return
	}
	if len(dbr.Data.([]database.Ban)) == 0 {
		mock_errors.WriteBadRequest(w, "User is not banned")
		return
	}

	err = db.NewQuery(r, 100).DeleteBan(database.UserRequestParams{BroadcasterID: broadcasterID, UserID: bannedUser.ID})
	if err != nil {
		mock_errors.WriteServerError(w, "Error deleting user's ban: "+err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
