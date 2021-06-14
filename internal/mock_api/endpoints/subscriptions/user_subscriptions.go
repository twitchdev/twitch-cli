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

var userSubscriptionsMethodsSupported = map[string]bool{
	http.MethodGet:    true,
	http.MethodPost:   false,
	http.MethodDelete: false,
	http.MethodPatch:  false,
	http.MethodPut:    false,
}

var userSubscriptionsScopesByMethod = map[string][]string{
	http.MethodGet:    {"user:read:subscriptions"},
	http.MethodPost:   {},
	http.MethodDelete: {},
	http.MethodPatch:  {},
	http.MethodPut:    {},
}

type UserSubscriptions struct{}

func (e UserSubscriptions) Path() string { return "/subscriptions/user" }

func (e UserSubscriptions) GetRequiredScopes(method string) []string {
	return userSubscriptionsScopesByMethod[method]
}

func (e UserSubscriptions) ValidMethod(method string) bool {
	return userSubscriptionsMethodsSupported[method]
}

func (e UserSubscriptions) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	db = r.Context().Value("db").(database.CLIDatabase)

	switch r.Method {
	case http.MethodGet:
		getUserSubscriptions(w, r)
		break
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

func getUserSubscriptions(w http.ResponseWriter, r *http.Request) {
	userCtx := r.Context().Value("auth").(authentication.UserAuthentication)

	if userCtx.UserID != r.URL.Query().Get("user_id") {
		mock_errors.WriteUnauthorized(w, "user_id does not match token")
		return
	}

	if r.URL.Query().Get("broadcaster_id") == "" {
		mock_errors.WriteBadRequest(w, "broadcaster_id is required")
		return
	}
	s := database.Subscription{
		BroadcasterID: r.URL.Query().Get("broadcaster_id"),
		UserID:        userCtx.UserID,
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
		Data: dbr.Data,
	}

	j, _ := json.Marshal(body)

	w.Write(j)
}
