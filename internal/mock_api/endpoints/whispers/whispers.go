// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package whispers

import (
	"encoding/json"
	"net/http"

	"github.com/twitchdev/twitch-cli/internal/database"
	"github.com/twitchdev/twitch-cli/internal/mock_api/authentication"
	"github.com/twitchdev/twitch-cli/internal/mock_api/mock_errors"
)

var whispersMethodsSupported = map[string]bool{
	http.MethodGet:    false,
	http.MethodPost:   true,
	http.MethodDelete: false,
	http.MethodPatch:  false,
	http.MethodPut:    false,
}

var whispersScopesByMethod = map[string][]string{
	http.MethodGet:    {},
	http.MethodPost:   {"user:manage:whispers"},
	http.MethodDelete: {},
	http.MethodPatch:  {},
	http.MethodPut:    {},
}

type PostWhisperRequestBody struct {
	Message string `json:"message"`
}

type Whispers struct{}

func (e Whispers) Path() string { return "/whispers" }

func (e Whispers) GetRequiredScopes(method string) []string {
	return whispersScopesByMethod[method]
}

func (e Whispers) ValidMethod(method string) bool {
	return whispersMethodsSupported[method]
}

func (e Whispers) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	db = r.Context().Value("db").(database.CLIDatabase)

	switch r.Method {
	case http.MethodPost:
		postWhispers(w, r)
		break
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func postWhispers(w http.ResponseWriter, r *http.Request) {
	userCtx := r.Context().Value("auth").(authentication.UserAuthentication)
	if !userCtx.MatchesSpecifiedIDParam(r, "from_user_id") {
		mock_errors.WriteUnauthorized(w, "from_user_id does not match token")
		return
	}

	fromUserID := r.URL.Query().Get("from_user_id")
	if fromUserID == "" {
		mock_errors.WriteBadRequest(w, "Missing required parameter from_user_id")
		return
	}

	toUserID := r.URL.Query().Get("to_user_id")
	if toUserID == "" {
		mock_errors.WriteBadRequest(w, "Missing required parameter to_user_id")
		return
	}

	if fromUserID == toUserID {
		mock_errors.WriteBadRequest(w, "The IDs on from_user_id and to_user_id cannot be the same ID")
		return
	}

	// Check if user exists
	user, err := db.NewQuery(r, 100).GetUser(database.User{ID: toUserID})
	if err != nil {
		mock_errors.WriteServerError(w, "error pulling to_user_id from database: "+err.Error())
		return
	}
	if user.ID == "" {
		mock_errors.WriteNotFound(w, "User specified in to_user_id doesn't exist")
		return
	}

	var body PostWhisperRequestBody
	err = json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		mock_errors.WriteBadRequest(w, "Body unable to be parsed")
		return
	}

	if body.Message == "" {
		mock_errors.WriteBadRequest(w, "Message field must be present and not contain an empty string")
		return
	}

	if len(body.Message) > 10000 {
		mock_errors.WriteBadRequest(w, "Message must be less than 10,000 characters")
		return
	}

	// Chat is not supported in Mock API, so we're pretending this worked.
	// This implementation also has no support for suspended users, blocked users, or users with whispers disabled

	w.WriteHeader(http.StatusNoContent)
}
