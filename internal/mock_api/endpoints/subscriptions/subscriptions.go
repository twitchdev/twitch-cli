// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package subscriptions

import (
	"encoding/json"
	"net/http"

	"github.com/twitchdev/twitch-cli/internal/database"
	"github.com/twitchdev/twitch-cli/internal/mock_api/authentication"
	"github.com/twitchdev/twitch-cli/internal/mock_api/mock_errors"
	"github.com/twitchdev/twitch-cli/internal/models"
)

var broadcasterSubscriptionsMethodsSupported = map[string]bool{
	http.MethodGet:    true,
	http.MethodPost:   false,
	http.MethodDelete: false,
	http.MethodPatch:  false,
	http.MethodPut:    false,
}

var broadcasterSubscriptionsScopesByMethod = map[string][]string{
	http.MethodGet:    {"channel:read:subscriptions"},
	http.MethodPost:   {},
	http.MethodDelete: {},
	http.MethodPatch:  {},
	http.MethodPut:    {},
}

type BroadcasterSubscriptions struct{}

func (e BroadcasterSubscriptions) Path() string { return "/subscriptions" }

func (e BroadcasterSubscriptions) GetRequiredScopes(method string) []string {
	return broadcasterSubscriptionsScopesByMethod[method]
}

func (e BroadcasterSubscriptions) ValidMethod(method string) bool {
	return broadcasterSubscriptionsMethodsSupported[method]
}

func (e BroadcasterSubscriptions) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	db = r.Context().Value("db").(database.CLIDatabase)

	switch r.Method {
	case http.MethodGet:
		getBroadcasterSubscriptions(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

func getBroadcasterSubscriptions(w http.ResponseWriter, r *http.Request) {
	userCtx := r.Context().Value("auth").(authentication.UserAuthentication)

	if !userCtx.MatchesBroadcasterIDParam(r) {
		mock_errors.WriteUnauthorized(w, "broadcaster_id does not match token")
		return
	}

	s := database.Subscription{
		BroadcasterID: userCtx.UserID,
		UserID:        r.URL.Query().Get("user_id"),
	}

	dbr, err := db.NewQuery(r, 100).GetSubscriptions(s)
	if err != nil {
		mock_errors.WriteServerError(w, "error fetching subscriptions")
		return
	}

	if len(dbr.Data.([]database.Subscription)) == 0 {
		dbr.Data = []database.Subscription{}
	}

	body := models.APIResponse{
		Data:  dbr.Data,
		Total: &dbr.Total,
		// This would usually be something like tier 1 = 1 pt, tier 2 = 2 pts, tier 3 = 6 pts. For simplicity, return total instead
		Points: dbr.Total,
	}

	if dbr.Cursor != "" {
		body.Pagination = &models.APIPagination{
			Cursor: dbr.Cursor,
		}
	}

	j, _ := json.Marshal(body)

	w.Write(j)
}
