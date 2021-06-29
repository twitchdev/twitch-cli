// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package teams

import (
	"encoding/json"
	"net/http"

	"github.com/twitchdev/twitch-cli/internal/database"
	"github.com/twitchdev/twitch-cli/internal/mock_api/mock_errors"
	"github.com/twitchdev/twitch-cli/internal/models"
)

var teamsMethodsSupported = map[string]bool{
	http.MethodGet:    true,
	http.MethodPost:   false,
	http.MethodDelete: false,
	http.MethodPatch:  false,
	http.MethodPut:    false,
}

var teamsScopesByMethod = map[string][]string{
	http.MethodGet:    {},
	http.MethodPost:   {},
	http.MethodDelete: {},
	http.MethodPatch:  {},
	http.MethodPut:    {},
}

type Teams struct{}

func (e Teams) Path() string { return "/teams" }

func (e Teams) GetRequiredScopes(method string) []string {
	return teamsScopesByMethod[method]
}

func (e Teams) ValidMethod(method string) bool {
	return teamsMethodsSupported[method]
}

func (e Teams) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	db = r.Context().Value("db").(database.CLIDatabase)

	switch r.Method {
	case http.MethodGet:
		getTeams(w, r)
		break
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func getTeams(w http.ResponseWriter, r *http.Request) {
	if r.URL.Query().Get("name") == "" && r.URL.Query().Get("id") == "" {
		mock_errors.WriteBadRequest(w, "one of name or id is required")
		return
	}
	if r.URL.Query().Get("name") != "" && r.URL.Query().Get("id") != "" {
		mock_errors.WriteBadRequest(w, "only one of name or id is required")
		return
	}
	dbr, err := db.NewQuery(r, 100).GetTeam(database.Team{ID: r.URL.Query().Get("id"), TeamName: r.URL.Query().Get("name")})
	if err != nil {
		mock_errors.WriteServerError(w, "error fetching teams")
		return
	}

	teams := dbr.Data.([]database.Team)

	bytes, _ := json.Marshal(models.APIResponse{Data: teams})
	w.Write(bytes)
}
