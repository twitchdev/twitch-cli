// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package channels

import (
	"encoding/json"
	"net/http"

	"github.com/twitchdev/twitch-cli/internal/database"
	"github.com/twitchdev/twitch-cli/internal/mock_api/authentication"
	"github.com/twitchdev/twitch-cli/internal/mock_api/mock_errors"
	"github.com/twitchdev/twitch-cli/internal/models"
)

var editorMethodsSupported = map[string]bool{
	http.MethodGet:    true,
	http.MethodPost:   false,
	http.MethodDelete: false,
	http.MethodPatch:  false,
	http.MethodPut:    false,
}

var editorScopesByMethod = map[string][]string{
	http.MethodGet:    {"channel:read:editors"},
	http.MethodPost:   {},
	http.MethodDelete: {},
	http.MethodPatch:  {},
	http.MethodPut:    {},
}

type Editors struct{}

func (e Editors) Path() string { return "/channels/editors" }

func (e Editors) GetRequiredScopes(method string) []string {
	return editorScopesByMethod[method]
}

func (e Editors) ValidMethod(method string) bool {
	return editorMethodsSupported[method]
}

func (e Editors) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	db = r.Context().Value("db").(database.CLIDatabase)
	switch r.Method {
	case http.MethodGet:
		getEditors(w, r)
		break
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		break
	}

}

func getEditors(w http.ResponseWriter, r *http.Request) {
	broadcasterID := r.URL.Query().Get("broadcaster_id")
	userCtx := r.Context().Value("auth").(authentication.UserAuthentication)

	if broadcasterID == "" {
		mock_errors.WriteBadRequest(w, "Broacaster ID is required")
		return
	}

	if broadcasterID != userCtx.UserID {
		mock_errors.WriteUnauthorized(w, "Broadcaster ID does not match token")
		return
	}

	dbr, err := db.NewQuery(r, 100).GetEditors(database.User{ID: broadcasterID})
	if err != nil {
		mock_errors.WriteServerError(w, err.Error())
		return
	}
	if len(dbr.Data.([]database.Editor)) == 0 {
		dbr.Data = []database.Editor{}
	}
	response := models.APIResponse{
		Data: &dbr.Data,
	}

	bytes, _ := json.Marshal(response)
	w.Write(bytes)
}
