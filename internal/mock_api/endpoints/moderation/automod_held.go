// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package moderation

import (
	"encoding/json"
	"net/http"

	"github.com/twitchdev/twitch-cli/internal/database"
	"github.com/twitchdev/twitch-cli/internal/mock_api/authentication"
	"github.com/twitchdev/twitch-cli/internal/mock_api/mock_errors"
)

var automodHeldMethodsSupported = map[string]bool{
	http.MethodGet:    false,
	http.MethodPost:   true,
	http.MethodDelete: false,
	http.MethodPatch:  false,
	http.MethodPut:    false,
}

var automodHeldScopesByMethod = map[string][]string{
	http.MethodGet:    {},
	http.MethodPost:   {"moderator:manage:automod"},
	http.MethodDelete: {},
	http.MethodPatch:  {},
	http.MethodPut:    {},
}

type AutomodHeld struct{}

type PostAutomodHeldBody struct {
	UserID    string `json:"user_id"`
	MessageID string `json:"msg_id"`
	Action    string `json:"action"`
}

func (e AutomodHeld) Path() string { return "/moderation/automod/message" }

func (e AutomodHeld) GetRequiredScopes(method string) []string {
	return automodHeldScopesByMethod[method]
}

func (e AutomodHeld) ValidMethod(method string) bool {
	return automodHeldMethodsSupported[method]
}

func (e AutomodHeld) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	db = r.Context().Value("db").(database.CLIDatabase)

	switch r.Method {
	case http.MethodPost:
		getAutomodHeld(w, r)
		break
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func getAutomodHeld(w http.ResponseWriter, r *http.Request) {
	userCtx := r.Context().Value("auth").(authentication.UserAuthentication)
	body := PostAutomodHeldBody{}

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		mock_errors.WriteBadRequest(w, "error reading body")
		return
	}

	if userCtx.UserID != body.UserID {
		mock_errors.WriteUnauthorized(w, "user_id must match token")
		return
	}

	if body.Action != "ALLOW" && body.Action != "DENY" {
		mock_errors.WriteBadRequest(w, "action must be one of ALLOW or DENY")
		return
	}

	if body.MessageID == "" {
		mock_errors.WriteBadRequest(w, "msg_id is required")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
