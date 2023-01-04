// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package hype_train

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/twitchdev/twitch-cli/internal/database"
	"github.com/twitchdev/twitch-cli/internal/mock_api/authentication"
	"github.com/twitchdev/twitch-cli/internal/mock_api/mock_errors"
	"github.com/twitchdev/twitch-cli/internal/models"
	"github.com/twitchdev/twitch-cli/internal/util"
)

var hypeTrainEventsMethodsSupported = map[string]bool{
	http.MethodGet:    true,
	http.MethodPost:   false,
	http.MethodDelete: false,
	http.MethodPatch:  false,
	http.MethodPut:    false,
}

var hypeTrainEventsScopesByMethod = map[string][]string{
	http.MethodGet:    {"channel:read:hype_train"},
	http.MethodPost:   {},
	http.MethodDelete: {},
	http.MethodPatch:  {},
	http.MethodPut:    {},
}

type HypeTrainEvents struct{}

type HypeTrainEventsResponse struct {
	ID             string             `json:"id"`
	EventType      string             `json:"event_type"`
	EventTimestamp string             `json:"event_timestamp"`
	Version        string             `json:"version"`
	EventData      HypeTrainEventData `json:"event_data"`
}

type HypeTrainEventData struct {
	ID               string                  `json:"id"`
	BroadcasterID    string                  `json:"broadcaster_id"`
	CooldownEndTime  string                  `json:"cooldown_end_time"`
	ExpiresAt        string                  `json:"expires_at"`
	Goal             int                     `json:"goal"`
	LastContribution HypeTrainContribution   `json:"last_contribution"`
	Level            int                     `json:"level"`
	StartedAt        string                  `json:"started_at"`
	TopContributions []HypeTrainContribution `json:"top_contributions"`
	Total            int                     `json:"total"`
}

type HypeTrainContribution struct {
	Total int    `json:"total"`
	Type  string `json:"type"`
	User  string `json:"user"`
}

func (e HypeTrainEvents) Path() string { return "/hypetrain/events" }

func (e HypeTrainEvents) GetRequiredScopes(method string) []string {
	return hypeTrainEventsScopesByMethod[method]
}

func (e HypeTrainEvents) ValidMethod(method string) bool {
	return hypeTrainEventsMethodsSupported[method]
}

func (e HypeTrainEvents) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	db = r.Context().Value("db").(database.CLIDatabase)

	switch r.Method {
	case http.MethodGet:
		getHypeTrainEvents(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

func getHypeTrainEvents(w http.ResponseWriter, r *http.Request) {
	userCtx := r.Context().Value("auth").(authentication.UserAuthentication)
	first := 1
	events := []HypeTrainEventsResponse{}
	if !userCtx.MatchesBroadcasterIDParam(r) {
		mock_errors.WriteUnauthorized(w, "Broadcaster ID does not match token.")
		return
	}

	firstParam := r.URL.Query().Get("first")
	if firstParam != "" {
		first, _ = strconv.Atoi(firstParam)
	}

	if first == 0 || first > 100 {
		mock_errors.WriteBadRequest(w, "first must be greater than 0 and less than 100")
		return
	}

	goal := int(util.RandomInt(10 * 1000))

	for i := 0; i < first; i++ {
		c := HypeTrainContribution{
			Total: goal / (int(util.RandomInt(int64(i+2)) + 1)),
			Type:  "BITS",
			User:  userCtx.UserID,
		}
		h := HypeTrainEventsResponse{
			ID:             util.RandomGUID(),
			EventType:      "hypetrain.progression",
			EventTimestamp: util.GetTimestamp().Format(time.RFC3339),
			Version:        "1.0",
			EventData: HypeTrainEventData{
				ID:               util.RandomGUID(),
				BroadcasterID:    userCtx.UserID,
				CooldownEndTime:  util.GetTimestamp().Add(60 * time.Minute).Format(time.RFC3339),
				ExpiresAt:        util.GetTimestamp().Add(10 * time.Minute).Format(time.RFC3339),
				Level:            int(util.RandomInt(5)) + 1,
				Goal:             goal,
				StartedAt:        util.GetTimestamp().Add(-10 * time.Minute).Format(time.RFC3339),
				LastContribution: c,
				TopContributions: []HypeTrainContribution{c},
				Total:            c.Total,
			},
		}
		events = append(events, h)
	}

	bytes, _ := json.Marshal(
		models.APIResponse{
			Data:       events,
			Pagination: &models.APIPagination{}, // Since these are randomly generated, true pagination isn't implemented
		},
	)
	w.Write(bytes)
}
