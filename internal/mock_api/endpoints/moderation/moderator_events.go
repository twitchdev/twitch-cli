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

var moderatorEventsMethodsSupported = map[string]bool{
	http.MethodGet:    true,
	http.MethodPost:   false,
	http.MethodDelete: false,
	http.MethodPatch:  false,
	http.MethodPut:    false,
}

var moderatorEventsScopesByMethod = map[string][]string{
	http.MethodGet:    {"moderation:read"},
	http.MethodPost:   {},
	http.MethodDelete: {},
	http.MethodPatch:  {},
	http.MethodPut:    {},
}

type ModeratorEvents struct{}

func (e ModeratorEvents) Path() string { return "/moderation/moderators/events" }

func (e ModeratorEvents) GetRequiredScopes(method string) []string {
	return moderatorEventsScopesByMethod[method]
}

func (e ModeratorEvents) ValidMethod(method string) bool {
	return moderatorEventsMethodsSupported[method]
}

func (e ModeratorEvents) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	db = r.Context().Value("db").(database.CLIDatabase)

	switch r.Method {
	case http.MethodGet:
		getModeratorEvents(w, r)
		break
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func getModeratorEvents(w http.ResponseWriter, r *http.Request) {
	userCtx := r.Context().Value("auth").(authentication.UserAuthentication)
	dbr := &database.DBResponse{}
	var err error
	if !userCtx.MatchesBroadcasterIDParam(r) {
		mock_errors.WriteBadRequest(w, "broadcaster_id does not match token")
		return
	}
	bans := []database.ModeratorAction{}
	userIDs := r.URL.Query()["user_id"]
	if len(userIDs) > 0 {
		for _, id := range userIDs {
			dbr, err := db.NewQuery(r, 100).GetModeratorEvents(database.UserRequestParams{BroadcasterID: userCtx.UserID, UserID: id})
			if err != nil {
				mock_errors.WriteServerError(w, "error fetching bans")
				return
			}
			bans = append(bans, dbr.Data.([]database.ModeratorAction)...)
		}
	} else {
		dbr, err = db.NewQuery(r, 100).GetModeratorEvents(database.UserRequestParams{BroadcasterID: userCtx.UserID})
		if err != nil {
			log.Print(err)
			mock_errors.WriteServerError(w, "error fetching bans")
			return
		}
		bans = append(bans, dbr.Data.([]database.ModeratorAction)...)
	}
	apiResponse := models.APIResponse{Data: bans}
	if dbr.Cursor != "" {
		apiResponse.Pagination = &models.APIPagination{Cursor: dbr.Cursor}
	}

	bytes, _ := json.Marshal(apiResponse)
	w.Write(bytes)
}
