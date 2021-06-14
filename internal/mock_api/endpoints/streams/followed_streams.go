// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package streams

import (
	"encoding/json"
	"net/http"

	"github.com/twitchdev/twitch-cli/internal/database"
	"github.com/twitchdev/twitch-cli/internal/mock_api/authentication"
	"github.com/twitchdev/twitch-cli/internal/mock_api/mock_errors"
	"github.com/twitchdev/twitch-cli/internal/models"
)

var followedStreamsMethodsSupported = map[string]bool{
	http.MethodGet:    true,
	http.MethodPost:   false,
	http.MethodDelete: false,
	http.MethodPatch:  false,
	http.MethodPut:    false,
}

var followedStreamsScopesByMethod = map[string][]string{
	http.MethodGet:    {"user:read:follows", "user:edit:follows"},
	http.MethodPost:   {},
	http.MethodDelete: {},
	http.MethodPatch:  {},
	http.MethodPut:    {},
}

type FollowedStreams struct{}

func (e FollowedStreams) Path() string { return "/streams/followed" }

func (e FollowedStreams) GetRequiredScopes(method string) []string {
	return followedStreamsScopesByMethod[method]
}

func (e FollowedStreams) ValidMethod(method string) bool {
	return followedStreamsMethodsSupported[method]
}

func (e FollowedStreams) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	db = r.Context().Value("db").(database.CLIDatabase)

	switch r.Method {
	case http.MethodGet:
		getFollowedStreams(w, r)
		break
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func getFollowedStreams(w http.ResponseWriter, r *http.Request) {
	userCtx := r.Context().Value("auth").(authentication.UserAuthentication)

	if userCtx.UserID != r.URL.Query().Get("user_id") {
		mock_errors.WriteUnauthorized(w, "user_id must match the token user")
		return
	}

	dbr, err := db.NewQuery(r, 100).GetFollowedStreams(userCtx.UserID)
	if err != nil {
		mock_errors.WriteServerError(w, "error getting streams")
		return
	}

	streams := dbr.Data.([]database.Stream)

	apiResponse := models.APIResponse{Data: streams}
	if len(streams) == dbr.Limit {
		apiResponse.Pagination = &models.APIPagination{Cursor: dbr.Cursor}
	}

	bytes, _ := json.Marshal(apiResponse)
	w.Write(bytes)
}
