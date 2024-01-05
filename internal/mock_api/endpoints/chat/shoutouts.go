// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package chat

import (
	"net/http"

	"github.com/twitchdev/twitch-cli/internal/database"
	"github.com/twitchdev/twitch-cli/internal/mock_api/authentication"
	"github.com/twitchdev/twitch-cli/internal/mock_api/mock_errors"
)

var shoutoutsMethodsSupported = map[string]bool{
	http.MethodGet:    false,
	http.MethodPost:   true,
	http.MethodDelete: false,
	http.MethodPatch:  false,
	http.MethodPut:    false,
}

var shoutoutsScopesByMethod = map[string][]string{
	http.MethodGet:    {},
	http.MethodPost:   {"moderator:manage:shoutouts"},
	http.MethodDelete: {},
	http.MethodPatch:  {},
	http.MethodPut:    {},
}

type PostShoutoutsRequestBody struct {
	SlowMode                      *bool `json:"slow_mode"`
	SlowModeWaitTime              *int  `json:"slow_mode_wait_time"`
	FollowerMode                  *bool `json:"follower_mode"`
	FollowerModeDuration          *int  `json:"follower_mode_duration"`
	SubscriberMode                *bool `json:"subscriber_mode"`
	EmoteMode                     *bool `json:"emote_mode"`
	UniqueChatMode                *bool `json:"unique_chat_mode"`
	NonModeratorChatDelay         *bool `json:"non_moderator_chat_delay"`
	NonModeratorChatDelayDuration *int  `json:"non_moderator_chat_delay_duration"`
}
type Shoutouts struct{}

func (e Shoutouts) Path() string { return "/chat/shoutouts" }

func (e Shoutouts) GetRequiredScopes(method string) []string {
	return shoutoutsScopesByMethod[method]
}

func (e Shoutouts) ValidMethod(method string) bool {
	return shoutoutsMethodsSupported[method]
}

func (e Shoutouts) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	db = r.Context().Value("db").(database.CLIDatabase)

	switch r.Method {
	case http.MethodPost:
		postShoutouts(w, r)
		break
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func postShoutouts(w http.ResponseWriter, r *http.Request) {
	userCtx := r.Context().Value("auth").(authentication.UserAuthentication)
	if !userCtx.MatchesModeratorIDParam(r) {
		mock_errors.WriteUnauthorized(w, "Moderator ID does not match token.")
		return
	}

	fromBroadcasterId := r.URL.Query().Get("from_broadcaster_id")
	if fromBroadcasterId == "" {
		mock_errors.WriteBadRequest(w, "Missing required parameter from_broadcaster_id")
		return
	}

	toBroadcasterId := r.URL.Query().Get("to_broadcaster_id")
	if toBroadcasterId == "" {
		mock_errors.WriteBadRequest(w, "Missing required parameter to_broadcaster_id")
		return
	}

	moderatorID := r.URL.Query().Get("moderator_id")
	if moderatorID == "" {
		mock_errors.WriteBadRequest(w, "Missing required parameter moderator_id")
		return
	}

	fromBroadcaster, err := db.NewQuery(r, 100).GetUser(database.User{ID: fromBroadcasterId})
	if err != nil {
		mock_errors.WriteServerError(w, "error fetching fromBrodcasterId")
		return
	}
	if fromBroadcaster.ID == "" {
		mock_errors.WriteBadRequest(w, "Invalid from_broadcaser_id: No broadcaster by that ID exists")
		return
	}

	toBroadcaster, err := db.NewQuery(r, 100).GetUser(database.User{ID: toBroadcasterId})
	if err != nil {
		mock_errors.WriteServerError(w, "error fetching toBrodcasterId")
		return
	}
	if toBroadcaster.ID == "" {
		mock_errors.WriteBadRequest(w, "Invalid to_broadcaser_id: No broadcaster by that ID exists")
		return
	}

	// Verify user is a moderator or is the broadcaster
	isModerator := false
	if fromBroadcasterId == moderatorID {
		isModerator = true
	} else {
		moderatorListDbr, err := db.NewQuery(r, 1000).GetModeratorsForBroadcaster(fromBroadcasterId)
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

	// No connection to chat on here, and no way to GET or PATCH announcements via API
	// For the time being, we just ingest it and pretend it worked (HTTP 204)
	w.WriteHeader(http.StatusNoContent)
}
