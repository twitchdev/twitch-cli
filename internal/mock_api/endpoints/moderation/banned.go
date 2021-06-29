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

var bansMethodsSupported = map[string]bool{
	http.MethodGet:    true,
	http.MethodPost:   false,
	http.MethodDelete: false,
	http.MethodPatch:  false,
	http.MethodPut:    false,
}

var bansScopesByMethod = map[string][]string{
	http.MethodGet:    {"moderation:read"},
	http.MethodPost:   {},
	http.MethodDelete: {},
	http.MethodPatch:  {},
	http.MethodPut:    {},
}

type Bans struct{}

func (e Bans) Path() string { return "/moderation/banned" }

func (e Bans) GetRequiredScopes(method string) []string {
	return bansScopesByMethod[method]
}

func (e Bans) ValidMethod(method string) bool {
	return bansMethodsSupported[method]
}

func (e Bans) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	db = r.Context().Value("db").(database.CLIDatabase)

	switch r.Method {
	case http.MethodGet:
		getBans(w, r)
		break
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func getBans(w http.ResponseWriter, r *http.Request) {
	userCtx := r.Context().Value("auth").(authentication.UserAuthentication)
	dbr := &database.DBResponse{}
	var err error
	if !userCtx.MatchesBroadcasterIDParam(r) {
		mock_errors.WriteUnauthorized(w, "broadcaster_id does not match token")
		return
	}
	bans := []database.Ban{}
	userIDs := r.URL.Query()["user_id"]
	if len(userIDs) > 0 {
		for _, id := range userIDs {
			dbr, err := db.NewQuery(r, 100).GetBans(database.UserRequestParams{BroadcasterID: userCtx.UserID, UserID: id})
			if err != nil {
				mock_errors.WriteServerError(w, "error fetching bans")
				return
			}
			bans = append(bans, dbr.Data.([]database.Ban)...)
		}
	} else {
		dbr, err = db.NewQuery(r, 100).GetBans(database.UserRequestParams{BroadcasterID: userCtx.UserID})
		if err != nil {
			log.Print(err)
			mock_errors.WriteServerError(w, "error fetching bans")
			return
		}
		bans = append(bans, dbr.Data.([]database.Ban)...)
	}
	apiResponse := models.APIResponse{Data: bans}
	if dbr.Cursor != "" {
		apiResponse.Pagination = &models.APIPagination{Cursor: dbr.Cursor}
	}

	bytes, _ := json.Marshal(apiResponse)
	w.Write(bytes)
}
