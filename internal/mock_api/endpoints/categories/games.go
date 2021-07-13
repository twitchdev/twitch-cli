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

var gamesMethodsSupported = map[string]bool{
	http.MethodGet:    true,
	http.MethodPost:   false,
	http.MethodDelete: false,
	http.MethodPatch:  false,
	http.MethodPut:    false,
}

var gamesScopesByMethod = map[string][]string{
	http.MethodGet:    {},
	http.MethodPost:   {},
	http.MethodDelete: {},
	http.MethodPatch:  {},
	http.MethodPut:    {},
}

type Games struct{}

func (e Games) Path() string { return "/games" }

func (e Games) GetRequiredScopes(method string) []string {
	return gamesScopesByMethod[method]
}

func (e Games) ValidMethod(method string) bool {
	return gamesMethodsSupported[method]
}

func (e Games) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	db = r.Context().Value("db").(database.CLIDatabase)

	switch r.Method {
	case http.MethodGet:
		getGames(w, r)
		break
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func getGames(w http.ResponseWriter, r *http.Request) {
	games := []database.Category{}
	ids := r.URL.Query()["id"]
	names := r.URL.Query()["name"]

	if len(ids) == 0 && len(names) == 0 {
		mock_errors.WriteBadRequest(w, "at least one name or id is required")
		return
	}

	if len(ids)+len(names) > 100 {
		mock_errors.WriteBadRequest(w, "you may only pass up to 100 ids and names")
		return
	}
	for _, id := range ids {
		dbr, err := db.NewQuery(r, 100).GetCategories(database.Category{ID: id})
		if err != nil {
			mock_errors.WriteServerError(w, "error getting category")
			return
		}
		game := dbr.Data.([]database.Category)
		games = append(games, game...)
	}
	for _, name := range names {
		dbr, err := db.NewQuery(r, 100).GetCategories(database.Category{Name: name})
		if err != nil {
			mock_errors.WriteServerError(w, "error getting category")
			return
		}
		game := dbr.Data.([]database.Category)
		games = append(games, game...)
	}
	for i, g := range games {
		games[i].BoxartURL = fmt.Sprintf("https://static-cdn.jtvnw.net/ttv-boxart/%v-{width}x{height}.jpg", url.PathEscape(g.Name))
	}
	bytes, _ := json.Marshal(models.APIResponse{Data: games})
	w.Write(bytes)
}
