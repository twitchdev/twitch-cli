// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package chat

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/twitchdev/twitch-cli/internal/database"
	"github.com/twitchdev/twitch-cli/internal/mock_api/authentication"
	"github.com/twitchdev/twitch-cli/internal/mock_api/mock_errors"
)

var announcementsMethodsSupported = map[string]bool{
	http.MethodGet:    false,
	http.MethodPost:   true,
	http.MethodDelete: false,
	http.MethodPatch:  false,
	http.MethodPut:    false,
}

var announcementsScopesByMethod = map[string][]string{
	http.MethodGet:    {},
	http.MethodPost:   {"moderator:manage:announcements"},
	http.MethodDelete: {},
	http.MethodPatch:  {},
	http.MethodPut:    {},
}

type PostAnnouncementsRequestBody struct {
	Message string `json:"message"`
	Color   string `json:"color"`
}

type Announcements struct{}

func (e Announcements) Path() string { return "/chat/announcements" }

func (e Announcements) GetRequiredScopes(method string) []string {
	return announcementsScopesByMethod[method]
}

func (e Announcements) ValidMethod(method string) bool {
	return announcementsMethodsSupported[method]
}

func (e Announcements) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	db = r.Context().Value("db").(database.CLIDatabase)

	switch r.Method {
	case http.MethodPost:
		postAnnouncements(w, r)
		break
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func postAnnouncements(w http.ResponseWriter, r *http.Request) {
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

	var body PostAnnouncementsRequestBody
	err = json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		mock_errors.WriteBadRequest(w, "Body unable to be parsed")
		return
	}

	if body.Message == "" {
		mock_errors.WriteBadRequest(w, "The message field in the request's body is required.")
		return
	}

	colorLowerCase := strings.ToLower(body.Color)
	if colorLowerCase != "" && colorLowerCase != "blue" && colorLowerCase != "green" && colorLowerCase != "orange" && colorLowerCase != "purple" && colorLowerCase != "primary" {
		mock_errors.WriteBadRequest(w, "The specific color is not valid")
		return
	}

	// No connection to chat on here, and no way to GET or PATCH announcements via API
	// For the time being, we just ingest it and pretend it worked (HTTP 204)
	w.WriteHeader(http.StatusNoContent)
}
