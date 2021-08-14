// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package streams

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/mattn/go-sqlite3"
	"github.com/twitchdev/twitch-cli/internal/database"
	"github.com/twitchdev/twitch-cli/internal/mock_api/authentication"
	"github.com/twitchdev/twitch-cli/internal/mock_api/mock_errors"
	"github.com/twitchdev/twitch-cli/internal/models"
)

var streamTagsMethodsSupported = map[string]bool{
	http.MethodGet:    true,
	http.MethodPost:   false,
	http.MethodDelete: false,
	http.MethodPatch:  false,
	http.MethodPut:    true,
}

var streamTagsScopesByMethod = map[string][]string{
	http.MethodGet:    {},
	http.MethodPost:   {},
	http.MethodDelete: {},
	http.MethodPatch:  {},
	http.MethodPut:    {"channel:manage:broadcast"},
}

type StreamTags struct{}

type PutBodyStreamTags struct {
	TagIDs []string `json:"tag_ids"`
}

func (e StreamTags) Path() string { return "/streams/tags" }

func (e StreamTags) GetRequiredScopes(method string) []string {
	return streamTagsScopesByMethod[method]
}

func (e StreamTags) ValidMethod(method string) bool {
	return streamTagsMethodsSupported[method]
}

func (e StreamTags) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	db = r.Context().Value("db").(database.CLIDatabase)

	switch r.Method {
	case http.MethodGet:
		getStreamTags(w, r)
		break
	case http.MethodPut:
		putStreamTags(w, r)
		break
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func getStreamTags(w http.ResponseWriter, r *http.Request) {
	userCtx := r.Context().Value("auth").(authentication.UserAuthentication)

	if !userCtx.MatchesBroadcasterIDParam(r) {
		mock_errors.WriteUnauthorized(w, "broadcaster_id does not match token")
		return
	}

	dbr, err := db.NewQuery(r, 100).GetStreamTags(userCtx.UserID)
	if err != nil {
		log.Print(err)
		mock_errors.WriteServerError(w, "error fetching tags")
		return
	}
	tagResponse := dbr.Data.([]database.Tag)

	apiResponse := models.APIResponse{
		Data: convertTags(tagResponse),
	}

	if len(tagResponse) == dbr.Limit {
		apiResponse.Pagination = &models.APIPagination{Cursor: dbr.Cursor}
	}

	bytes, _ := json.Marshal(apiResponse)
	w.Write(bytes)
}

func putStreamTags(w http.ResponseWriter, r *http.Request) {
	userCtx := r.Context().Value("auth").(authentication.UserAuthentication)

	if !userCtx.MatchesBroadcasterIDParam(r) {
		mock_errors.WriteUnauthorized(w, "broadcaster_id does not match token")
		return
	}

	body := PutBodyStreamTags{}

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		mock_errors.WriteBadRequest(w, "error parsing body")
		return
	}

	err = db.NewQuery(r, 100).DeleteAllStreamTags(userCtx.UserID)
	if err != nil {
		log.Print(err)
		mock_errors.WriteServerError(w, err.Error())
		return
	}
	for _, tag := range body.TagIDs {
		err = db.NewQuery(r, 100).InsertStreamTag(database.StreamTag{UserID: userCtx.UserID, TagID: tag})
		if err != nil {
			if database.DatabaseErrorIs(err, sqlite3.ErrConstraintForeignKey) {
				mock_errors.WriteBadRequest(w, "invalid tag provided")
				return
			}
			mock_errors.WriteServerError(w, err.Error())
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}
