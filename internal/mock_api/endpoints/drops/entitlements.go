// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package drops

import (
	"encoding/json"
	"net/http"

	"github.com/twitchdev/twitch-cli/internal/database"
	"github.com/twitchdev/twitch-cli/internal/mock_api/authentication"
	"github.com/twitchdev/twitch-cli/internal/mock_api/mock_errors"
	"github.com/twitchdev/twitch-cli/internal/models"
)

var dropsEntitlementsMethodsSupported = map[string]bool{
	http.MethodGet:    true,
	http.MethodPost:   false,
	http.MethodDelete: false,
	http.MethodPatch:  false,
	http.MethodPut:    false,
}

var dropsEntitlementsScopesByMethod = map[string][]string{
	http.MethodGet:    {},
	http.MethodPost:   {},
	http.MethodDelete: {},
	http.MethodPatch:  {},
	http.MethodPut:    {},
}

type DropsEntitlements struct{}

func (e DropsEntitlements) Path() string { return "/entitlements/drops" }

func (e DropsEntitlements) GetRequiredScopes(method string) []string {
	return dropsEntitlementsScopesByMethod[method]
}

func (e DropsEntitlements) ValidMethod(method string) bool {
	return dropsEntitlementsMethodsSupported[method]
}

func (e DropsEntitlements) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	db = r.Context().Value("db").(database.CLIDatabase)

	switch r.Method {
	case http.MethodGet:
		getEntitlements(w, r)
		break
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func getEntitlements(w http.ResponseWriter, r *http.Request) {
	userCtx := r.Context().Value("auth").(authentication.UserAuthentication)
	id := r.URL.Query().Get("id")
	userID := r.URL.Query().Get("user_id")
	gameID := r.URL.Query().Get("game_id")

	if userCtx.UserID != "" && userID != "" {
		mock_errors.WriteBadRequest(w, "user_id is invalid when using user access token")
		return
	}
	if userCtx.UserID != "" {
		userID = userCtx.UserID
	}
	e := database.DropsEntitlement{UserID: userID, GameID: gameID, ID: id}
	dbr, err := db.NewQuery(r, 1000).GetDropsEntitlements(e)
	if err != nil {
		mock_errors.WriteServerError(w, "error fetching entitlements")
		return
	}
	entitlements := dbr.Data.([]database.DropsEntitlement)
	if len(entitlements) == 0 {
		entitlements = []database.DropsEntitlement{}
	}
	apiResponse := models.APIResponse{
		Data: entitlements,
	}
	if len(entitlements) == dbr.Limit {
		apiResponse.Pagination = &models.APIPagination{
			Cursor: dbr.Cursor,
		}
	}
	bytes, err := json.Marshal(apiResponse)
	w.Write(bytes)
}
