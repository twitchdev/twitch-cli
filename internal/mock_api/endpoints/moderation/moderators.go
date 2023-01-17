// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package moderation

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/twitchdev/twitch-cli/internal/database"
	"github.com/twitchdev/twitch-cli/internal/mock_api/authentication"
	"github.com/twitchdev/twitch-cli/internal/mock_api/mock_errors"
	"github.com/twitchdev/twitch-cli/internal/models"
)

var moderatorsMethodsSupported = map[string]bool{
	http.MethodGet:    true,
	http.MethodPost:   true,
	http.MethodDelete: true,
	http.MethodPatch:  false,
	http.MethodPut:    false,
}

var moderatorsScopesByMethod = map[string][]string{
	http.MethodGet:    {"moderation:read", "channel:manage:moderators"},
	http.MethodPost:   {"channel:manage:moderators"},
	http.MethodDelete: {"channel:manage:moderators"},
	http.MethodPatch:  {},
	http.MethodPut:    {},
}

type Moderators struct{}

func (e Moderators) Path() string { return "/moderation/moderators" }

func (e Moderators) GetRequiredScopes(method string) []string {
	return moderatorsScopesByMethod[method]
}

func (e Moderators) ValidMethod(method string) bool {
	return moderatorsMethodsSupported[method]
}

func (e Moderators) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	db = r.Context().Value("db").(database.CLIDatabase)

	switch r.Method {
	case http.MethodGet:
		getModerators(w, r)
		break
	case http.MethodPost:
		postModerators(w, r)
		break
	case http.MethodDelete:
		deleteModerators(w, r)
		break
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func getModerators(w http.ResponseWriter, r *http.Request) {
	userCtx := r.Context().Value("auth").(authentication.UserAuthentication)
	if !userCtx.MatchesBroadcasterIDParam(r) {
		mock_errors.WriteUnauthorized(w, "broadcaster_id does not match token")
		return
	}

	dbr := &database.DBResponse{}
	var err error

	bans := []database.Moderator{}
	userIDs := r.URL.Query()["user_id"]
	if len(userIDs) > 0 {
		for _, id := range userIDs {
			dbr, err := db.NewQuery(r, 100).GetModerators(database.UserRequestParams{BroadcasterID: userCtx.UserID, UserID: id})
			if err != nil {
				mock_errors.WriteServerError(w, "error fetching bans")
				return
			}
			bans = append(bans, dbr.Data.([]database.Moderator)...)
		}
	} else {
		dbr, err = db.NewQuery(r, 100).GetModerators(database.UserRequestParams{BroadcasterID: userCtx.UserID})
		if err != nil {
			log.Print(err)
			mock_errors.WriteServerError(w, "error fetching bans")
			return
		}
		bans = append(bans, dbr.Data.([]database.Moderator)...)
	}
	apiResponse := models.APIResponse{Data: bans}
	if dbr.Cursor != "" {
		apiResponse.Pagination = &models.APIPagination{Cursor: dbr.Cursor}
	}

	bytes, _ := json.Marshal(apiResponse)
	w.Write(bytes)
}

func postModerators(w http.ResponseWriter, r *http.Request) {
	userCtx := r.Context().Value("auth").(authentication.UserAuthentication)
	if !userCtx.MatchesBroadcasterIDParam(r) {
		mock_errors.WriteUnauthorized(w, "broadcaster_id does not match token")
		return
	}

	broadcasterID := r.URL.Query().Get("broadcaster_id")
	if broadcasterID == "" {
		mock_errors.WriteBadRequest(w, "Missing required parameter broadcaster_id")
		return
	}

	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		mock_errors.WriteBadRequest(w, "Missing required parameter user_id")
		return
	}

	// Check if user exists
	user, err := db.NewQuery(r, 100).GetUser(database.User{ID: userID})
	if err != nil {
		mock_errors.WriteServerError(w, "error pulling user: "+err.Error())
		return
	}
	if user.ID == "" {
		mock_errors.WriteBadRequest(w, "User specified in user_id doesn't exist")
		return
	}

	// Check if user is already a moderator on the channel, or is the broadcaster
	isModerator := false
	if broadcasterID == userID {
		isModerator = true
	} else {
		moderatorListDbr, err := db.NewQuery(r, 1000).GetModeratorsForBroadcaster(broadcasterID)
		if err != nil {
			mock_errors.WriteServerError(w, err.Error())
			return
		}
		for _, mod := range moderatorListDbr.Data.([]database.Moderator) {
			if mod.UserID == userID {
				isModerator = true
			}
		}
	}
	if isModerator {
		mock_errors.WriteUnauthorized(w, "The user specified in parameter moderator_id is already a moderator on this channel")
		return
	}

	// Check if the user is banned from the channel
	dbr, err := db.NewQuery(r, 100).GetBans(database.UserRequestParams{BroadcasterID: broadcasterID, UserID: userID})
	if err != nil {
		mock_errors.WriteServerError(w, "error fetching bans")
		return
	}
	if len(dbr.Data.([]database.Ban)) != 0 {
		mock_errors.WriteBadRequest(w, "User cannot become a moderator because they are banned.")
		return
	}

	// Check if the user is a VIP
	vipDbr, err := db.NewQuery(r, 100).GetVIPsByBroadcaster(broadcasterID)
	if err != nil {
		mock_errors.WriteServerError(w, "error fetching vips: "+err.Error())
		return
	}
	for _, vip := range vipDbr.Data.([]database.VIP) {
		if vip.UserID == userID {
			mock_errors.WriteUnprocessableEntity(w, "The user is currently a VIP. They must be removed as a VIP to become a moderator.")
			return
		}
	}

	err = db.NewQuery(r, 100).AddModerator(database.UserRequestParams{BroadcasterID: broadcasterID, UserID: userID})
	if err != nil {
		mock_errors.WriteServerError(w, "error adding moderator")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func deleteModerators(w http.ResponseWriter, r *http.Request) {
	userCtx := r.Context().Value("auth").(authentication.UserAuthentication)
	if !userCtx.MatchesBroadcasterIDParam(r) {
		mock_errors.WriteUnauthorized(w, "broadcaster_id does not match token")
		return
	}

	broadcasterID := r.URL.Query().Get("broadcaster_id")
	if broadcasterID == "" {
		mock_errors.WriteBadRequest(w, "Missing required parameter broadcaster_id")
		return
	}

	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		mock_errors.WriteBadRequest(w, "Missing required parameter user_id")
		return
	}

	// Check if user is a moderator on the channel
	isModerator := false
	moderatorListDbr, err := db.NewQuery(r, 1000).GetModeratorsForBroadcaster(broadcasterID)
	if err != nil {
		mock_errors.WriteServerError(w, err.Error())
		return
	}
	for _, mod := range moderatorListDbr.Data.([]database.Moderator) {
		if mod.UserID == userID {
			isModerator = true
		}
	}
	if !isModerator {
		mock_errors.WriteBadRequest(w, "The user specified in parameter moderator_id is not a moderator on this channel")
		return
	}

	// Check if the user is the broadcaster
	if userID == broadcasterID {
		mock_errors.WriteBadRequest(w, "The broadcaster cannot be removed as a moderator")
		return
	}

	err = db.NewQuery(r, 100).RemoveModerator(broadcasterID, userID)
	if err != nil {
		mock_errors.WriteServerError(w, "error removing moderator: "+err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
