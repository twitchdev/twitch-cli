// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package channel_points

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/twitchdev/twitch-cli/internal/database"
	"github.com/twitchdev/twitch-cli/internal/mock_api/authentication"
	"github.com/twitchdev/twitch-cli/internal/mock_api/mock_errors"
	"github.com/twitchdev/twitch-cli/internal/models"
)

var redemptionMethodsSupported = map[string]bool{
	http.MethodGet:    true,
	http.MethodPost:   false,
	http.MethodDelete: false,
	http.MethodPatch:  true,
	http.MethodPut:    false,
}

var redemptionScopesByMethod = map[string][]string{
	http.MethodGet:    {"channel:read:redemptions", "channel:manage:redemptions"},
	http.MethodPost:   {},
	http.MethodDelete: {},
	http.MethodPatch:  {"channel:manage:redemptions"},
	http.MethodPut:    {},
}

type Redemption struct{}

type PatchRedemptionBody struct {
	Status string `json:"status"`
}

func (e Redemption) Path() string { return "/channel_points/custom_rewards/redemptions" }

func (e Redemption) GetRequiredScopes(method string) []string {
	return redemptionScopesByMethod[method]
}

func (e Redemption) ValidMethod(method string) bool {
	return redemptionMethodsSupported[method]
}

func (e Redemption) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	db = r.Context().Value("db").(database.CLIDatabase)

	switch r.Method {
	case http.MethodGet:
		getRedemptions(w, r)
		break
	case http.MethodPatch:
		patchRedemptions(w, r)
		break
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func getRedemptions(w http.ResponseWriter, r *http.Request) {
	userCtx := r.Context().Value("auth").(authentication.UserAuthentication)
	reward_id := r.URL.Query().Get("reward_id")

	id := r.URL.Query().Get("id")
	status := r.URL.Query().Get("status")
	sort := r.URL.Query().Get("sort")

	if id == "" && status == "" {
		mock_errors.WriteBadRequest(w, "The status query parameter is required if you don't specify the id query parameter.")
		return
	}

	if !userCtx.MatchesBroadcasterIDParam(r) {
		mock_errors.WriteUnauthorized(w, "Broadcaster ID mismatch")
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

	if reward_id == "" {
		mock_errors.WriteBadRequest(w, "reward_id is required")
		return
	}

	if sort != "" && sort != "OLDEST" && sort != "NEWEST" {
		mock_errors.WriteBadRequest(w, "Invalid sort requested")
		return
	}

	cpr := database.ChannelPointsRedemption{BroadcasterID: userCtx.UserID, ID: id, RedemptionStatus: status, RewardID: reward_id}

	dbr, err := db.NewQuery(r, 100).GetChannelPointsRedemption(cpr, sort)
	if err != nil {
		mock_errors.WriteServerError(w, err.Error())
		return
	}

	response := models.APIResponse{
		Data: &dbr.Data,
	}
	if len(dbr.Data.([]database.ChannelPointsRedemption)) == 0 {
		response.Data = []database.ChannelPointsRedemption{}
	}
	if dbr.Limit == len(dbr.Data.([]database.ChannelPointsRedemption)) {
		p := &models.APIPagination{}
		p.Cursor = dbr.Cursor
		response.Pagination = p
	}

	bytes, _ := json.Marshal(response)
	w.Write(bytes)
}

func patchRedemptions(w http.ResponseWriter, r *http.Request) {
	userCtx := r.Context().Value("auth").(authentication.UserAuthentication)
	reward_id := r.URL.Query().Get("reward_id")

	ids := r.URL.Query()["id"]

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

	if len(ids) == 0 || reward_id == "" {
		mock_errors.WriteBadRequest(w, "Missing required parameter")
		return
	}

	if len(ids) > 50 {
		mock_errors.WriteBadRequest(w, "This endpoint only supports up to 50 IDs at a time")
		return
	}

	var body PatchRedemptionBody

	err = json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		mock_errors.WriteServerError(w, err.Error())
		return
	}

	if body.Status != "FULFILLED" && body.Status != "CANCELED" {
		mock_errors.WriteBadRequest(w, "Invalid status provided")
		return
	}

	responseData := []database.ChannelPointsRedemption{}
	for _, id := range ids {
		reward := database.ChannelPointsRedemption{RewardID: reward_id, BroadcasterID: userCtx.UserID, ID: id}
		dbr, err := db.NewQuery(r, 100).GetChannelPointsRedemption(reward, "")
		if err != nil {
			mock_errors.WriteServerError(w, err.Error())
			return
		}

		data := dbr.Data.([]database.ChannelPointsRedemption)
		if len(data) == 0 {
			mock_errors.WriteNotFound(w, fmt.Sprintf("Redemption ID %v does not exist or is not owned by the broadcaster", id))
			return
		}

		if data[0].RedemptionStatus != "UNFULFILLED" {
			mock_errors.WriteBadRequest(w, fmt.Sprintf("ID %v is already fulfilled or cancelled", id))
			return
		}

		update := database.ChannelPointsRedemption{ID: id, RedemptionStatus: body.Status}
		err = db.NewQuery(r, 100).UpdateChannelPointsRedemption(update)
		if err != nil {
			mock_errors.WriteServerError(w, err.Error())
			return
		}
		data[0].RedemptionStatus = body.Status
		responseData = append(responseData, data[0])
	}

	bytes, _ := json.Marshal(models.APIResponse{Data: responseData})
	w.Write(bytes)
}
