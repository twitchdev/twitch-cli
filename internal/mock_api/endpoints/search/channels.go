// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package search

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/twitchdev/twitch-cli/internal/database"
	"github.com/twitchdev/twitch-cli/internal/mock_api/mock_errors"
	"github.com/twitchdev/twitch-cli/internal/models"
)

var searchChannelsMethodsSupported = map[string]bool{
	http.MethodGet:    true,
	http.MethodPost:   false,
	http.MethodDelete: false,
	http.MethodPatch:  false,
	http.MethodPut:    false,
}

var searchChannelsScopesByMethod = map[string][]string{
	http.MethodGet:    {},
	http.MethodPost:   {},
	http.MethodDelete: {},
	http.MethodPatch:  {},
	http.MethodPut:    {},
}

type SearchChannels struct{}

func (e SearchChannels) Path() string { return "/search/channels" }

func (e SearchChannels) GetRequiredScopes(method string) []string {
	return searchChannelsScopesByMethod[method]
}

func (e SearchChannels) ValidMethod(method string) bool {
	return searchChannelsMethodsSupported[method]
}

func (e SearchChannels) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	db = r.Context().Value("db").(database.CLIDatabase)

	switch r.Method {
	case http.MethodGet:
		searchChannels(w, r)
		break
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
func searchChannels(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")
	live_only := false

	if query == "" {
		mock_errors.WriteBadRequest(w, "query is required")
		return
	}

	if r.URL.Query().Get("live_only") != "" {
		live_only, _ = strconv.ParseBool(r.URL.Query().Get("live_only"))
	}
	println(live_only)
	dbr, err := db.NewQuery(r, 100).SearchChannels(query, live_only)
	if err != nil {
		log.Print(err)
		mock_errors.WriteServerError(w, "error searching categories")
		return
	}

	categories := dbr.Data.([]database.SearchChannel)

	apiResponse := models.APIResponse{Data: categories}

	if len(categories) == dbr.Limit {
		apiResponse.Pagination = &models.APIPagination{Cursor: dbr.Cursor}
	}

	bytes, _ := json.Marshal(apiResponse)
	w.Write(bytes)
}
