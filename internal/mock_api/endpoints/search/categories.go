// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package search

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/twitchdev/twitch-cli/internal/database"
	"github.com/twitchdev/twitch-cli/internal/mock_api/mock_errors"
	"github.com/twitchdev/twitch-cli/internal/models"
)

var searchCategoriesMethodsSupported = map[string]bool{
	http.MethodGet:    true,
	http.MethodPost:   false,
	http.MethodDelete: false,
	http.MethodPatch:  false,
	http.MethodPut:    false,
}

var searchCategoriesScopesByMethod = map[string][]string{
	http.MethodGet:    {},
	http.MethodPost:   {},
	http.MethodDelete: {},
	http.MethodPatch:  {},
	http.MethodPut:    {},
}

type SearchCategories struct{}

func (e SearchCategories) Path() string { return "/search/categories" }

func (e SearchCategories) GetRequiredScopes(method string) []string {
	return searchCategoriesScopesByMethod[method]
}

func (e SearchCategories) ValidMethod(method string) bool {
	return searchCategoriesMethodsSupported[method]
}

func (e SearchCategories) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	db = r.Context().Value("db").(database.CLIDatabase)

	switch r.Method {
	case http.MethodGet:
		searchCategories(w, r)
		break
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func searchCategories(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")

	if query == "" {
		mock_errors.WriteBadRequest(w, "query is required")
		return
	}

	dbr, err := db.NewQuery(r, 100).SearchCategories(query)
	if err != nil {
		mock_errors.WriteServerError(w, "error searching categories")
		return
	}

	categories := dbr.Data.([]database.Category)
	for i, c := range categories {
		categories[i].BoxartURL = fmt.Sprintf("https://static-cdn.jtvnw.net/ttv-boxart/%v-{width}x{height}.jpg", url.PathEscape(c.Name))
	}
	apiResponse := models.APIResponse{Data: categories}

	if len(categories) == dbr.Limit {
		apiResponse.Pagination = &models.APIPagination{Cursor: dbr.Cursor}
	}

	bytes, _ := json.Marshal(apiResponse)
	w.Write(bytes)
}
