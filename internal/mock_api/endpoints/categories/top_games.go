// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package categories

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/twitchdev/twitch-cli/internal/database"
	"github.com/twitchdev/twitch-cli/internal/mock_api/mock_errors"
	"github.com/twitchdev/twitch-cli/internal/models"
)

var topGamesMethodsSupported = map[string]bool{
	http.MethodGet:    true,
	http.MethodPost:   false,
	http.MethodDelete: false,
	http.MethodPatch:  false,
	http.MethodPut:    false,
}

var topGamesScopesByMethod = map[string][]string{
	http.MethodGet:    {},
	http.MethodPost:   {},
	http.MethodDelete: {},
	http.MethodPatch:  {},
	http.MethodPut:    {},
}

type TopGames struct{}

func (e TopGames) Path() string { return "/games/top" }

func (e TopGames) GetRequiredScopes(method string) []string {
	return topGamesScopesByMethod[method]
}

func (e TopGames) ValidMethod(method string) bool {
	return topGamesMethodsSupported[method]
}

func (e TopGames) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	db = r.Context().Value("db").(database.CLIDatabase)

	switch r.Method {
	case http.MethodGet:
		getTopGames(w, r)
		break
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func getTopGames(w http.ResponseWriter, r *http.Request) {
	dbr, err := db.NewQuery(r, 100).GetTopGames()
	if err != nil {
		mock_errors.WriteServerError(w, "error fetching entitlements")
		return
	}
	games := dbr.Data.([]database.Category)
	if len(games) == 0 {
		games = []database.Category{}
	}

	for i, g := range games {
		games[i].BoxartURL = fmt.Sprintf("https://static-cdn.jtvnw.net/ttv-boxart/%v-{width}x{height}.jpg", url.PathEscape(g.Name))
	}
	apiResponse := models.APIResponse{
		Data: games,
	}
	if len(games) == dbr.Limit {
		apiResponse.Pagination = &models.APIPagination{
			Cursor: dbr.Cursor,
		}
	}

	bytes, err := json.Marshal(apiResponse)
	w.Write(bytes)
}
