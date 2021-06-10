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

var bannedEventsMethodsSupported = map[string]bool{
	http.MethodGet:    true,
	http.MethodPost:   false,
	http.MethodDelete: false,
	http.MethodPatch:  false,
	http.MethodPut:    false,
}

var bannedEventsScopesByMethod = map[string][]string{
	http.MethodGet:    {"moderation:read"},
	http.MethodPost:   {},
	http.MethodDelete: {},
	http.MethodPatch:  {},
	http.MethodPut:    {},
}

type BannedEvents struct{}

func (e BannedEvents) Path() string { return "/moderation/banned/events" }

func (e BannedEvents) GetRequiredScopes(method string) []string {
	return bannedEventsScopesByMethod[method]
}

func (e BannedEvents) ValidMethod(method string) bool {
	return bannedEventsMethodsSupported[method]
}

func (e BannedEvents) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	db = r.Context().Value("db").(database.CLIDatabase)

	switch r.Method {
	case http.MethodGet:
		getBanEvents(w, r)
		break
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
func getBanEvents(w http.ResponseWriter, r *http.Request) {
	userCtx := r.Context().Value("auth").(authentication.UserAuthentication)
	dbr := &database.DBResponse{}
	var err error
	if !userCtx.MatchesBroadcasterIDParam(r) {
		mock_errors.WriteBadRequest(w, "broadcaster_id does not match token")
		return
	}
	bans := []database.BanEvent{}
	userIDs := r.URL.Query()["user_id"]
	if len(userIDs) > 0 {
		for _, id := range userIDs {
			dbr, err := db.NewQuery(r, 100).GetBanEvents(database.UserRequestParams{BroadcasterID: userCtx.UserID, UserID: id})
			if err != nil {
				mock_errors.WriteServerError(w, "error fetching bans")
				return
			}
			bans = append(bans, dbr.Data.([]database.BanEvent)...)
		}
	} else {
		dbr, err = db.NewQuery(r, 100).GetBanEvents(database.UserRequestParams{BroadcasterID: userCtx.UserID})
		if err != nil {
			log.Print(err)
			mock_errors.WriteServerError(w, "error fetching bans")
			return
		}
		bans = append(bans, dbr.Data.([]database.BanEvent)...)
	}
	apiResponse := models.APIResponse{Data: bans}
	if dbr.Cursor != "" {
		apiResponse.Pagination = &models.APIPagination{Cursor: dbr.Cursor}
	}

	bytes, _ := json.Marshal(apiResponse)
	w.Write(bytes)
}
