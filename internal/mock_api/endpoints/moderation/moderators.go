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
	http.MethodPost:   false,
	http.MethodDelete: false,
	http.MethodPatch:  false,
	http.MethodPut:    false,
}

var moderatorsScopesByMethod = map[string][]string{
	http.MethodGet:    {"moderation:read"},
	http.MethodPost:   {},
	http.MethodDelete: {},
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
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func getModerators(w http.ResponseWriter, r *http.Request) {
	userCtx := r.Context().Value("auth").(authentication.UserAuthentication)
	dbr := &database.DBResponse{}
	var err error
	if !userCtx.MatchesBroadcasterIDParam(r) {
		mock_errors.WriteBadRequest(w, "broadcaster_id does not match token")
		return
	}
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
