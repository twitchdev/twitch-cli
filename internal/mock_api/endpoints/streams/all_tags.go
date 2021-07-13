// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package streams

import (
	"encoding/json"
	"net/http"

	"github.com/twitchdev/twitch-cli/internal/database"
	"github.com/twitchdev/twitch-cli/internal/mock_api/mock_errors"
	"github.com/twitchdev/twitch-cli/internal/models"
)

var allTagsMethodsSupported = map[string]bool{
	http.MethodGet:    true,
	http.MethodPost:   false,
	http.MethodDelete: false,
	http.MethodPatch:  false,
	http.MethodPut:    false,
}

var allTagsScopesByMethod = map[string][]string{
	http.MethodGet:    {},
	http.MethodPost:   {},
	http.MethodDelete: {},
	http.MethodPatch:  {},
	http.MethodPut:    {},
}

type AllTags struct{}

func (e AllTags) Path() string { return "/tags/streams" }

func (e AllTags) GetRequiredScopes(method string) []string {
	return allTagsScopesByMethod[method]
}

func (e AllTags) ValidMethod(method string) bool {
	return allTagsMethodsSupported[method]
}

func (e AllTags) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	db = r.Context().Value("db").(database.CLIDatabase)

	switch r.Method {
	case http.MethodGet:
		getAllTags(w, r)
		break
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func getAllTags(w http.ResponseWriter, r *http.Request) {
	tagIDs := r.URL.Query()["tag_id"]
	dbResponse := database.DBResponse{}
	tags := []database.Tag{}

	if len(tagIDs) > 100 {
		mock_errors.WriteBadRequest(w, "only 100 tag_ids can be provided at a time")
		return
	}

	if len(tagIDs) > 0 {
		for _, id := range tagIDs {
			println(id)
			t := database.Tag{ID: id}
			dbr, err := db.NewQuery(r, 100).GetTags(t)
			if err != nil {
				mock_errors.WriteServerError(w, "error fetching tags")
				return
			}

			tagResponse := dbr.Data.([]database.Tag)
			tags = append(tags, tagResponse...)
		}
	} else {
		t := database.Tag{}
		dbr, err := db.NewQuery(r, 100).GetTags(t)
		if err != nil {
			mock_errors.WriteServerError(w, "error fetching tags")
			return
		}
		dbResponse = *dbr
		tagResponse := dbr.Data.([]database.Tag)
		tags = append(tags, tagResponse...)
	}

	apiResponse := models.APIResponse{
		Data: convertTags(tags),
	}

	if len(tagIDs) == 0 && len(tags) == dbResponse.Limit {
		apiResponse.Pagination = &models.APIPagination{Cursor: dbResponse.Cursor}
	}

	bytes, _ := json.Marshal(apiResponse)
	w.Write(bytes)
}
