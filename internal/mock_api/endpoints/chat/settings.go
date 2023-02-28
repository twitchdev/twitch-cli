// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package chat

import (
	"encoding/json"
	"net/http"

	"github.com/twitchdev/twitch-cli/internal/database"
	"github.com/twitchdev/twitch-cli/internal/mock_api/authentication"
	"github.com/twitchdev/twitch-cli/internal/mock_api/mock_errors"
	"github.com/twitchdev/twitch-cli/internal/models"
)

var settingsMethodsSupported = map[string]bool{
	http.MethodGet:    true,
	http.MethodPost:   false,
	http.MethodDelete: false,
	http.MethodPatch:  true,
	http.MethodPut:    false,
}

var settingsScopesByMethod = map[string][]string{
	http.MethodGet:    {},
	http.MethodPost:   {},
	http.MethodDelete: {},
	http.MethodPatch:  {"moderator:manage:chat_settings"},
	http.MethodPut:    {},
}

// Only used when the user isn't a moderator
type GetSettingsResponseUnprivileged struct {
	BroadcasterID                 string `json:"broadcaster_id"`
	SlowMode                      bool   `json:"slow_mode"`
	SlowModeWaitTime              int    `json:"slow_mode_wait_time"`
	FollowerMode                  bool   `json:"follower_mode"`
	FollowerModeDuration          int    `json:"follower_mode_duration"`
	SubscriberMode                bool   `json:"subscriber_mode"`
	EmoteMode                     bool   `json:"emote_mode"`
	UniqueChatMode                bool   `json:"unique_chat_mode"`
	NonModeratorChatDelay         bool   `json:"-"`
	NonModeratorChatDelayDuration int    `json:"-"`
}

type PatchSettingsRequestBody struct {
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

type PatchSettingsResponseBody struct {
	BroadcasterID                 string `json:"broadcaster_id"`
	ModeratorID                   string `json:"moderator_id"`
	SlowMode                      bool   `json:"slow_mode"`
	SlowModeWaitTime              int    `json:"slow_mode_wait_time"`
	FollowerMode                  bool   `json:"follower_mode"`
	FollowerModeDuration          int    `json:"follower_mode_duration"`
	SubscriberMode                bool   `json:"subscriber_mode"`
	EmoteMode                     bool   `json:"emote_mode"`
	UniqueChatMode                bool   `json:"unique_chat_mode"`
	NonModeratorChatDelay         bool   `json:"non_moderator_chat_delay"`
	NonModeratorChatDelayDuration int    `json:"non_moderator_chat_delay_duration"`
}

type Settings struct{}

func (e Settings) Path() string { return "/chat/settings" }

func (e Settings) GetRequiredScopes(method string) []string {
	return settingsScopesByMethod[method]
}

func (e Settings) ValidMethod(method string) bool {
	return settingsMethodsSupported[method]
}

func (e Settings) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	db = r.Context().Value("db").(database.CLIDatabase)

	switch r.Method {
	case http.MethodGet:
		getSettings(w, r)
		break
	case http.MethodPatch:
		patchSettings(w, r)
		break
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func getSettings(w http.ResponseWriter, r *http.Request) {
	broadcasterID := r.URL.Query().Get("broadcaster_id")
	if broadcasterID == "" {
		mock_errors.WriteBadRequest(w, "Missing required parameter broadcaster_id")
		return
	}

	broadcaster, err := db.NewQuery(r, 100).GetUser(database.User{ID: broadcasterID})
	if err != nil {
		mock_errors.WriteServerError(w, "error fetching broadcaster")
		return
	}
	if broadcaster.ID == "" {
		mock_errors.WriteBadRequest(w, "No broadcaster by that ID exists")
		return
	}

	dbr, err := db.NewQuery(r, 100).GetChatSettingsByBroadcaster(broadcasterID)
	if err != nil {
		mock_errors.WriteServerError(w, err.Error())
		return
	}

	settings := dbr.Data.([]database.ChatSettings)

	// Moderator check
	isModerator := false
	moderatorID := r.URL.Query().Get("moderator_id")
	userCtx := r.Context().Value("auth").(authentication.UserAuthentication)
	if moderatorID == "" && userCtx.MatchesBroadcasterIDParam(r) {
		// No Moderator ID was given, and the Broadcaster ID matches the user access token.
		isModerator = true
	} else {
		isModerator = validateModerator(w, r, moderatorID, broadcasterID)
	}

	apiResponse := models.APIResponse{
		Data: settings,
	}

	// User is not a moderator. Remove moderator fields from response
	if !isModerator {
		apiResponse.Data = []GetSettingsResponseUnprivileged{
			{
				BroadcasterID:        settings[0].BroadcasterID,
				SlowMode:             *settings[0].SlowMode,
				SlowModeWaitTime:     *settings[0].SlowModeWaitTime,
				FollowerMode:         *settings[0].FollowerMode,
				FollowerModeDuration: *settings[0].FollowerModeDuration,
				SubscriberMode:       *settings[0].SubscriberMode,
				EmoteMode:            *settings[0].EmoteMode,
				UniqueChatMode:       *settings[0].UniqueChatMode,
			},
		}
	}

	bytes, _ := json.Marshal(apiResponse)
	w.Write(bytes)
}

func patchSettings(w http.ResponseWriter, r *http.Request) {
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

	// Check if broadcaster exists
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

	var body PatchSettingsRequestBody
	err = json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		mock_errors.WriteBadRequest(w, "Body unable to be parsed")
		return
	}

	if body.SlowModeWaitTime != nil {
		if *body.SlowModeWaitTime < 3 || *body.SlowModeWaitTime > 120 {
			mock_errors.WriteBadRequest(w, "slow_mode_wait_time must be greater than or equal to 3, and less than or equal to 120")
			return
		}
	}

	if body.FollowerModeDuration != nil {
		if *body.FollowerModeDuration < 0 || *body.SlowModeWaitTime > 129600 {
			mock_errors.WriteBadRequest(w, "follower_mode_duration must be greater than or equal to 0, and less than or equal to 129600")
			return
		}
	}

	if body.NonModeratorChatDelayDuration != nil {
		if *body.NonModeratorChatDelayDuration != 2 && *body.NonModeratorChatDelayDuration != 4 && *body.NonModeratorChatDelayDuration != 6 {
			mock_errors.WriteBadRequest(w, "non_moderator_chat_delay_duration must be one of the following values: 2, 4, 6")
			return
		}
	}

	update := database.ChatSettings{
		BroadcasterID:                 broadcasterID,
		SlowMode:                      body.SlowMode,
		SlowModeWaitTime:              body.SlowModeWaitTime,
		FollowerMode:                  body.FollowerMode,
		FollowerModeDuration:          body.FollowerModeDuration,
		SubscriberMode:                body.SubscriberMode,
		EmoteMode:                     body.EmoteMode,
		UniqueChatMode:                body.UniqueChatMode,
		NonModeratorChatDelay:         body.NonModeratorChatDelay,
		NonModeratorChatDelayDuration: body.NonModeratorChatDelayDuration,
	}

	err = db.NewQuery(r, 100).UpdateChatSettings(update)
	if err != nil {
		mock_errors.WriteServerError(w, err.Error())
		return
	}

	dbr, err := db.NewQuery(r, 100).GetChatSettingsByBroadcaster(broadcasterID)
	cs := dbr.Data.([]database.ChatSettings)[0]

	settings := []PatchSettingsResponseBody{
		{
			BroadcasterID:                 cs.BroadcasterID,
			ModeratorID:                   moderatorID,
			SlowMode:                      *cs.SlowMode,
			SlowModeWaitTime:              *cs.SlowModeWaitTime,
			FollowerMode:                  *cs.FollowerMode,
			FollowerModeDuration:          *cs.FollowerModeDuration,
			SubscriberMode:                *cs.SubscriberMode,
			EmoteMode:                     *cs.EmoteMode,
			UniqueChatMode:                *cs.UniqueChatMode,
			NonModeratorChatDelay:         *cs.NonModeratorChatDelay,
			NonModeratorChatDelayDuration: *cs.NonModeratorChatDelayDuration,
		},
	}

	bytes, _ := json.Marshal(models.APIResponse{
		Data: settings,
	})
	w.Write(bytes)
}

func validateModerator(w http.ResponseWriter, r *http.Request, moderatorId string, broadcasterId string) bool {
	// Check if Moderator ID matches user access token
	userCtx := r.Context().Value("auth").(authentication.UserAuthentication)
	if !userCtx.MatchesModeratorIDParam(r) {
		return false
	}

	// Check if Moderator ID is the broadcaster
	if moderatorId == broadcasterId {
		return true
	}

	// Check if Moderator ID is a moderator of this channel
	moderatorListDbr, err := db.NewQuery(r, 1000).GetModeratorsForBroadcaster(broadcasterId)
	if err != nil {
		mock_errors.WriteServerError(w, err.Error())
		return false
	}
	for _, mod := range moderatorListDbr.Data.([]database.Moderator) {
		if mod.UserID == moderatorId {
			return true // Moderator found
		}
	}

	// Not found in moderator list
	return false
}
