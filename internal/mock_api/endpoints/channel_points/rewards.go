// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package channel_points

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/twitchdev/twitch-cli/internal/database"
	"github.com/twitchdev/twitch-cli/internal/mock_api/authentication"
	"github.com/twitchdev/twitch-cli/internal/mock_api/mock_errors"
	"github.com/twitchdev/twitch-cli/internal/models"
	"github.com/twitchdev/twitch-cli/internal/util"
)

var rewardMethodsSupported = map[string]bool{
	http.MethodGet:    true,
	http.MethodPost:   true,
	http.MethodDelete: true,
	http.MethodPatch:  true,
	http.MethodPut:    false,
}

var rewardScopesByMethod = map[string][]string{
	http.MethodGet:    {"channel:read:redemptions", "channel:manage:redemptions"},
	http.MethodPost:   {"channel:manage:redemptions"},
	http.MethodDelete: {"channel:manage:redemptions"},
	http.MethodPatch:  {"channel:manage:redemptions"},
	http.MethodPut:    {},
}

type Reward struct{}

type PatchAndPostRewardBody struct {
	Title                      string `json:"title"`
	Cost                       *int   `json:"cost"`
	RewardPrompt               string `json:"prompt"`
	IsEnabled                  *bool  `json:"is_enabled"`
	BackgroundColor            string `json:"background_color"`
	IsUserInputRequired        bool   `json:"is_user_input_requird"`
	StreamMaxEnabled           bool   `json:"is_max_per_stream_enabled"`
	StreamMaxCount             int    `json:"max_per_stream"`
	StreamUserMaxEnabled       bool   `json:"is_max_per_user_per_stream_enabled"`
	StreamUserMaxCount         int    `json:"max_per_user_per_stream"`
	GlobalCooldownEnabled      bool   `json:"is_global_cooldown_enabled"`
	GlobalCooldownSeconds      int    `json:"global_cooldown_seconds"`
	ShouldRedemptionsSkipQueue bool   `json:"should_redemptions_skip_request_queue"`
}

func (e Reward) Path() string { return "/channel_points/custom_rewards" }

func (e Reward) GetRequiredScopes(method string) []string {
	return rewardScopesByMethod[method]
}

func (e Reward) ValidMethod(method string) bool {
	return rewardMethodsSupported[method]
}

func (e Reward) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	db = r.Context().Value("db").(database.CLIDatabase)

	switch r.Method {
	case http.MethodGet:
		getRewards(w, r)
		break
	case http.MethodPost:
		postRewards(w, r)
		break
	case http.MethodPatch:
		patchRewards(w, r)
		break
	case http.MethodDelete:
		deleteRewards(w, r)
		break
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func getRewards(w http.ResponseWriter, r *http.Request) {
	userCtx := r.Context().Value("auth").(authentication.UserAuthentication)

	id := r.URL.Query().Get("id")
	//onlyManageableRewards := r.URL.Query().Get("only_manageable_rewards")

	if !userCtx.MatchesBroadcasterIDParam(r) {
		mock_errors.WriteUnauthorized(w, "Broadcaster ID does not match token.")
		return
	}

	user, err := db.NewQuery(r, 100).GetUser(database.User{ID: userCtx.UserID})
	if err != nil {
		mock_errors.WriteServerError(w, err.Error())
		return
	}

	if user.BroadcasterType == "" {
		mock_errors.WriteForbidden(w, "User is not an affiliate or partner.")
		return
	}

	cpr := database.ChannelPointsReward{BroadcasterID: userCtx.UserID}
	if id != "" {
		cpr.ID = id
	}

	dbr, err := db.NewQuery(r, 100).GetChannelPointsReward(cpr)
	if err != nil {
		mock_errors.WriteServerError(w, err.Error())
		return
	}
	data := dbr.Data.([]database.ChannelPointsReward)
	response := models.APIResponse{
		Data: data,
	}
	if len(dbr.Data.([]database.ChannelPointsReward)) == 0 {
		response.Data = []database.ChannelPointsReward{}
	}
	if len(dbr.Data.([]database.ChannelPointsReward)) == dbr.Limit {
		pagination := &models.APIPagination{}
		pagination.Cursor = dbr.Cursor
		response.Pagination = pagination
	}

	bytes, _ := json.Marshal(response)
	w.Write(bytes)
}

func postRewards(w http.ResponseWriter, r *http.Request) {
	userCtx := r.Context().Value("auth").(authentication.UserAuthentication)

	if !userCtx.MatchesBroadcasterIDParam(r) {
		mock_errors.WriteUnauthorized(w, "Broadcaster ID does not match token.")
		return
	}

	user, err := db.NewQuery(r, 100).GetUser(database.User{ID: userCtx.UserID})
	if err != nil {
		mock_errors.WriteServerError(w, err.Error())
		return
	}

	if user.BroadcasterType == "" {
		mock_errors.WriteForbidden(w, "User is not an affiliate or partner.")
		return
	}

	var body PatchAndPostRewardBody
	err = json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		mock_errors.WriteBadRequest(w, "Body unable to be parsed")
		return
	}

	if body.Title == "" || body.Cost == nil || *body.Cost == 0 {
		mock_errors.WriteBadRequest(w, "Title or cost misssing")
		return
	}

	if body.IsEnabled == nil {
		t := true
		body.IsEnabled = &t
	}
	if body.StreamMaxEnabled && body.StreamMaxCount == 0 {
		mock_errors.WriteBadRequest(w, "max_per_stream required if is_max_per_stream_enabled is true")
		return
	}

	if body.StreamUserMaxEnabled && body.StreamUserMaxCount == 0 {
		mock_errors.WriteBadRequest(w, "max_per_user_per_stream if is_max_per_user_per_stream_enabled is true")
		return
	}

	if body.GlobalCooldownEnabled && body.GlobalCooldownSeconds == 0 {
		mock_errors.WriteBadRequest(w, "global_cooldown_seconds required if is_global_cooldown_enabled is true")
		return
	}

	create := database.ChannelPointsReward{
		ID:                  util.RandomGUID(),
		BroadcasterID:       userCtx.UserID,
		RewardImage:         sql.NullString{},
		BackgroundColor:     body.BackgroundColor,
		IsEnabled:           body.IsEnabled,
		RewardPrompt:        body.RewardPrompt,
		Cost:                *body.Cost,
		Title:               body.Title,
		IsUserInputRequired: body.IsUserInputRequired,
		MaxPerStream: database.MaxPerStream{
			StreamMaxEnabled: body.StreamMaxEnabled,
			StreamMaxCount:   body.StreamUserMaxCount,
		},
		MaxPerUserPerStream: database.MaxPerUserPerStream{
			StreamUserMaxEnabled: body.StreamUserMaxEnabled,
			StreamMUserMaxCount:  body.StreamUserMaxCount,
		},
		GlobalCooldown: database.GlobalCooldown{
			GlobalCooldownEnabled: body.GlobalCooldownEnabled,
			GlobalCooldownSeconds: body.GlobalCooldownSeconds,
		},
		ShouldRedemptionsSkipQueue: body.ShouldRedemptionsSkipQueue,
	}

	err = db.NewQuery(r, 100).InsertChannelPointsReward(create)
	if err != nil {
		mock_errors.WriteServerError(w, err.Error())
		return
	}

	dbr, err := db.NewQuery(r, 100).GetChannelPointsReward(database.ChannelPointsReward{ID: create.ID})
	if err != nil {
		mock_errors.WriteServerError(w, err.Error())
		return
	}
	bytes, err := json.Marshal(models.APIResponse{Data: dbr.Data})
	w.Write(bytes)
}

func patchRewards(w http.ResponseWriter, r *http.Request) {
	userCtx := r.Context().Value("auth").(authentication.UserAuthentication)
	id := r.URL.Query().Get("id")

	if !userCtx.MatchesBroadcasterIDParam(r) {
		mock_errors.WriteUnauthorized(w, "Broadcaster ID does not match token.")
		return
	}
	user, err := db.NewQuery(r, 100).GetUser(database.User{ID: userCtx.UserID})
	if err != nil {
		mock_errors.WriteServerError(w, err.Error())
		return
	}

	if user.BroadcasterType == "" {
		mock_errors.WriteForbidden(w, "User is not an affiliate or partner.")
		return
	}

	if id == "" {
		mock_errors.WriteBadRequest(w, "ID is required")
		return
	}
	var body PatchAndPostRewardBody
	err = json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		mock_errors.WriteBadRequest(w, "Body unable to be parsed")
		return
	}

	if *body.Cost == 0 {
		mock_errors.WriteBadRequest(w, "Cost must be greater than 0")
		return
	}

	if body.StreamMaxEnabled && body.StreamMaxCount == 0 {
		mock_errors.WriteBadRequest(w, "max_per_stream required if is_max_per_stream_enabled is true")
		return
	}

	if body.StreamUserMaxEnabled && body.StreamUserMaxCount == 0 {
		mock_errors.WriteBadRequest(w, "max_per_user_per_stream if is_max_per_user_per_stream_enabled is true")
		return
	}

	if body.GlobalCooldownEnabled && body.GlobalCooldownSeconds == 0 {
		mock_errors.WriteBadRequest(w, "global_cooldown_seconds required if is_global_cooldown_enabled is true")
		return
	}

	update := database.ChannelPointsReward{
		ID:                  id,
		BroadcasterID:       userCtx.UserID,
		RewardImage:         sql.NullString{},
		BackgroundColor:     body.BackgroundColor,
		IsEnabled:           body.IsEnabled,
		RewardPrompt:        body.RewardPrompt,
		Cost:                *body.Cost,
		Title:               body.Title,
		IsUserInputRequired: body.IsUserInputRequired,
		MaxPerStream: database.MaxPerStream{
			StreamMaxEnabled: body.StreamMaxEnabled,
			StreamMaxCount:   body.StreamUserMaxCount,
		},
		MaxPerUserPerStream: database.MaxPerUserPerStream{
			StreamUserMaxEnabled: body.StreamUserMaxEnabled,
			StreamMUserMaxCount:  body.StreamUserMaxCount,
		},
		GlobalCooldown: database.GlobalCooldown{
			GlobalCooldownEnabled: body.GlobalCooldownEnabled,
			GlobalCooldownSeconds: body.GlobalCooldownSeconds,
		},
		ShouldRedemptionsSkipQueue: body.ShouldRedemptionsSkipQueue,
	}

	err = db.NewQuery(r, 100).UpdateChannelPointsReward(update)
	if err != nil {
		mock_errors.WriteServerError(w, err.Error())
		return
	}

	dbr, err := db.NewQuery(r, 100).GetChannelPointsReward(database.ChannelPointsReward{ID: id})
	if err != nil {
		mock_errors.WriteServerError(w, err.Error())
		return
	}
	bytes, err := json.Marshal(models.APIResponse{Data: dbr.Data})
	w.Write(bytes)
}

func deleteRewards(w http.ResponseWriter, r *http.Request) {
	userCtx := r.Context().Value("auth").(authentication.UserAuthentication)
	id := r.URL.Query().Get("id")

	if !userCtx.MatchesBroadcasterIDParam(r) {
		mock_errors.WriteUnauthorized(w, "Broadcaster ID does not match token.")
		return
	}
	user, err := db.NewQuery(r, 100).GetUser(database.User{ID: userCtx.UserID})
	if err != nil {
		mock_errors.WriteServerError(w, err.Error())
		return
	}

	if user.BroadcasterType == "" {
		mock_errors.WriteForbidden(w, "User is not an affiliate or partner.")
		return
	}

	if id == "" {
		mock_errors.WriteBadRequest(w, "ID is required")
		return
	}

	dbr, err := db.NewQuery(r, 100).GetChannelPointsReward(database.ChannelPointsReward{ID: id, BroadcasterID: userCtx.UserID})
	if err != nil {
		mock_errors.WriteServerError(w, err.Error())
		return
	}

	reward := dbr.Data.([]database.ChannelPointsReward)
	if len(reward) == 0 {
		mock_errors.WriteNotFound(w, "Custom reward not found for broadcaster")
		return
	}

	err = db.NewQuery(r, 100).DeleteChannelPointsReward(reward[0].ID)
	if err != nil {
		mock_errors.WriteServerError(w, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
