// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package moderation

import (
	"net/http"

	"github.com/twitchdev/twitch-cli/internal/database"
	"github.com/twitchdev/twitch-cli/internal/mock_api/authentication"
	"github.com/twitchdev/twitch-cli/internal/mock_api/mock_errors"
)

var chatMethodsSupported = map[string]bool{
	http.MethodGet:    false,
	http.MethodPost:   false,
	http.MethodDelete: true,
	http.MethodPatch:  false,
	http.MethodPut:    false,
}

var chatScopesByMethod = map[string][]string{
	http.MethodGet:    {},
	http.MethodPost:   {},
	http.MethodDelete: {"moderator:manage:chat_messages"},
	http.MethodPatch:  {},
	http.MethodPut:    {},
}

type Chat struct{}

func (e Chat) Path() string { return "/moderation/chat" }

func (e Chat) GetRequiredScopes(method string) []string {
	return chatScopesByMethod[method]
}

func (e Chat) ValidMethod(method string) bool {
	return chatMethodsSupported[method]
}

func (e Chat) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	db = r.Context().Value("db").(database.CLIDatabase)

	switch r.Method {
	case http.MethodDelete:
		deleteChat(w, r)
		break
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func deleteChat(w http.ResponseWriter, r *http.Request) {
	userCtx := r.Context().Value("auth").(authentication.UserAuthentication)
	if !userCtx.MatchesModeratorIDParam(r) {
		mock_errors.WriteUnauthorized(w, "Moderator ID does not match token.")
		return
	}

	broadcasterID := r.URL.Query().Get("broadcaster_id")
	if broadcasterID == "" {
		mock_errors.WriteBadRequest(w, "Missing required parameter broadcaster_id")
		return
	}

	moderatorID := r.URL.Query().Get("moderator_id")
	if moderatorID == "" {
		mock_errors.WriteBadRequest(w, "Missing required parameter moderator_id")
		return
	}

	broadcaster, err := db.NewQuery(r, 100).GetUser(database.User{ID: broadcasterID})
	if err != nil {
		mock_errors.WriteServerError(w, "error fetching broadcaster")
		return
	}
	if broadcaster.ID == "" {
		mock_errors.WriteUnauthorized(w, "The user specified in parameter moderator_id is not one of the broadcaster's moderators")
		return
	}

	// Verify user is a moderator or is the broadcaster
	isModerator := false
	if broadcasterID == moderatorID {
		isModerator = true
	} else {
		moderatorListDbr, err := db.NewQuery(r, 1000).GetModeratorsForBroadcaster(broadcasterID)
		if err != nil {
			mock_errors.WriteServerError(w, err.Error())
			return
		}
		for _, mod := range moderatorListDbr.Data.([]database.Moderator) {
			if mod.UserID == moderatorID {
				isModerator = true
			}
		}
	}
	if !isModerator {
		mock_errors.WriteUnauthorized(w, "The user specified in parameter moderator_id is not one of the broadcaster's moderators")
		return
	}

	messageID := r.URL.Query().Get("message_id")
	if messageID == "" {
		mock_errors.WriteBadRequest(w, "Missing required parameter message_id")
		return
	}

	w.WriteHeader(http.StatusNoContent)

	// No connection to chat on here, so we're just gonna pretend it worked (HTTP 204)
}
