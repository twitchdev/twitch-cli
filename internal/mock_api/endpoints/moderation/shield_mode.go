// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package moderation

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/twitchdev/twitch-cli/internal/database"
	"github.com/twitchdev/twitch-cli/internal/mock_api/authentication"
	"github.com/twitchdev/twitch-cli/internal/mock_api/mock_errors"
	"github.com/twitchdev/twitch-cli/internal/models"
	"github.com/twitchdev/twitch-cli/internal/util"
)

var shieldModeMethodsSupported = map[string]bool{
	http.MethodGet:    true,
	http.MethodPost:   false,
	http.MethodDelete: false,
	http.MethodPatch:  false,
	http.MethodPut:    true,
}

var shieldModeScopesByMethod = map[string][]string{
	http.MethodGet:    {"moderator:manage:shield_mode", "moderator:read:shield_mode"},
	http.MethodPost:   {},
	http.MethodDelete: {},
	http.MethodPatch:  {},
	http.MethodPut:    {"moderator:manage:shield_mode"},
}

type GetShieldModeStatusResponseBody struct {
	IsActive        bool   `json:"is_active"`
	ModeratorID     string `json:"moderator_id"`
	ModeratorName   string `json:"moderator_name"`
	ModeratorLogin  string `json:"moderator_login"`
	LastActivatedAt string `json:"last_activated_at"`
}

type PutShieldModeStatusRequestBody struct {
	IsActive bool `json:"is_active"`
}

type ShieldMode struct{}

func (e ShieldMode) Path() string { return "/moderation/shield_mode" }

func (e ShieldMode) GetRequiredScopes(method string) []string {
	return shieldModeScopesByMethod[method]
}

func (e ShieldMode) ValidMethod(method string) bool {
	return shieldModeMethodsSupported[method]
}

func (e ShieldMode) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	db = r.Context().Value("db").(database.CLIDatabase)

	switch r.Method {
	case http.MethodGet:
		getShieldModeStatus(w, r)
		break
	case http.MethodPut:
		putShieldModeStatus(w, r)
		break
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func getShieldModeStatus(w http.ResponseWriter, r *http.Request) {
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

	dbr, err := db.NewQuery(r, 100).GetChatSettingsByBroadcaster(broadcasterID)
	if err != nil {
		log.Print(err)
		mock_errors.WriteServerError(w, fmt.Sprintf("error fetching chat settings: %v", err.Error()))
		return
	}

	chatSettings := dbr.Data.([]database.ChatSettings)[0]

	shieldModeSettings := GetShieldModeStatusResponseBody{
		IsActive:        chatSettings.ShieldModeIsActive,
		ModeratorID:     chatSettings.ShieldModeModeratorID,
		ModeratorName:   chatSettings.ShieldModeModeratorName,
		ModeratorLogin:  chatSettings.ShieldModeModeratorLogin,
		LastActivatedAt: chatSettings.ShieldModeLastActivated,
	}

	bytes, _ := json.Marshal(models.APIResponse{Data: []GetShieldModeStatusResponseBody{shieldModeSettings}})
	w.Write(bytes)
}

func putShieldModeStatus(w http.ResponseWriter, r *http.Request) {
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

	var body PutShieldModeStatusRequestBody
	err = json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		mock_errors.WriteBadRequest(w, "Body unable to be parsed")
		return
	}

	// Get moderator info
	moderator, err := db.NewQuery(r, 100).GetUser(database.User{ID: moderatorID})

	var newChatSettings database.ChatSettings

	if body.IsActive {
		newChatSettings = database.ChatSettings{
			BroadcasterID:            broadcasterID,
			ShieldModeIsActive:       body.IsActive,
			ShieldModeModeratorID:    moderatorID,
			ShieldModeModeratorLogin: moderator.UserLogin,
			ShieldModeModeratorName:  moderator.DisplayName,
			ShieldModeLastActivated:  util.GetTimestamp().Format(time.RFC3339Nano),
		}
	} else {
		newChatSettings = database.ChatSettings{
			BroadcasterID:      broadcasterID,
			ShieldModeIsActive: body.IsActive,
		}
	}

	err = db.NewQuery(r, 100).UpdateChatSettings(newChatSettings)

	dbr, err := db.NewQuery(r, 100).GetChatSettingsByBroadcaster(broadcasterID)
	if err != nil {
		log.Print(err)
		mock_errors.WriteServerError(w, fmt.Sprintf("error fetching chat settings: %v", err.Error()))
		return
	}

	updatedChatSettings := dbr.Data.([]database.ChatSettings)[0]
	updatedShieldModeSettings := GetShieldModeStatusResponseBody{
		IsActive:        updatedChatSettings.ShieldModeIsActive,
		ModeratorID:     updatedChatSettings.ShieldModeModeratorID,
		ModeratorName:   updatedChatSettings.ShieldModeModeratorName,
		ModeratorLogin:  updatedChatSettings.ShieldModeModeratorLogin,
		LastActivatedAt: updatedChatSettings.ShieldModeLastActivated,
	}

	bytes, _ := json.Marshal([]GetShieldModeStatusResponseBody{updatedShieldModeSettings})
	w.Write(bytes)
}
