// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package goals

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/twitchdev/twitch-cli/internal/database"
	"github.com/twitchdev/twitch-cli/internal/mock_api/authentication"
	"github.com/twitchdev/twitch-cli/internal/mock_api/mock_errors"
	"github.com/twitchdev/twitch-cli/internal/models"
	"github.com/twitchdev/twitch-cli/internal/util"
)

var goalsMethodsSupported = map[string]bool{
	http.MethodGet:    true,
	http.MethodPost:   false,
	http.MethodDelete: false,
	http.MethodPatch:  false,
	http.MethodPut:    false,
}

var goalsScopesByMethod = map[string][]string{
	http.MethodGet:    {"channel:read:goals"},
	http.MethodPost:   {},
	http.MethodDelete: {},
	http.MethodPatch:  {},
	http.MethodPut:    {},
}

type Goals struct{}

type GetCreatorGoalsResponse struct {
	ID               string `json:"id"`
	BroadcasterID    string `json:"broadcaster_id"`
	BroadcasterName  string `json:"broadcaster_name"`
	BroadcasterLogin string `json:"broadcaster_login"`
	Type             string `json:"type"`
	Description      string `json:"description"`
	CurrentAmount    int    `json:"current_amount"`
	TargetAmount     int    `json:"target_amount"`
	CreatedAt        string `json:"created_at"`
}

func (e Goals) Path() string { return "/goals" }

func (e Goals) GetRequiredScopes(method string) []string {
	return goalsScopesByMethod[method]
}

func (e Goals) ValidMethod(method string) bool {
	return goalsMethodsSupported[method]
}

func (e Goals) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	db = r.Context().Value("db").(database.CLIDatabase)
	switch r.Method {
	case http.MethodGet:
		getGoals(w, r)
		break
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func getGoals(w http.ResponseWriter, r *http.Request) {
	userCtx := r.Context().Value("auth").(authentication.UserAuthentication)
	if !userCtx.MatchesBroadcasterIDParam(r) {
		mock_errors.WriteUnauthorized(w, "Broadcaster ID does not match token.")
		return
	}

	broadcasterID := r.URL.Query().Get("broadcaster_id")

	if broadcasterID == "" {
		mock_errors.WriteBadRequest(w, "query parameter broadcaster_id is required")
		return
	}

	user, err := db.NewQuery(r, 100).GetUser(database.User{ID: userCtx.UserID})
	if err != nil {
		mock_errors.WriteServerError(w, err.Error())
		return
	}

	targetAmount := rand.Intn(1500-300) + 300        // Between 300 and 1500
	currentAmount := rand.Intn(targetAmount-99) + 99 // Between 99 and targetAmount

	var lowerBound int64 = time.Now().Unix() - (43200 * 6 * 60) // 6 months ago -- (30 days in minutes) * 6 * 60 seconds
	randomTimestamp := time.Unix(rand.Int63n(time.Now().Unix()-lowerBound)+lowerBound, 0).Format(time.RFC3339)

	types := []string{
		"follower",
		"subscriber",
		"subscription_count",
		"new_subscription",
		"new_subscription_count",
	}
	rand.Seed(time.Now().UnixNano())
	goalType := types[rand.Intn(len(types))]

	description := fmt.Sprintf("Lets get to %v", targetAmount)
	if goalType == "follower" {
		description = fmt.Sprintf("%v followers!", description)
	} else {
		description = fmt.Sprintf("%v subs!", description)
	}

	goal := GetCreatorGoalsResponse{
		ID:               util.RandomGUID(),
		BroadcasterID:    userCtx.UserID,
		BroadcasterName:  user.DisplayName,
		BroadcasterLogin: user.UserLogin,
		Type:             goalType,
		Description:      description,
		CurrentAmount:    currentAmount,
		TargetAmount:     targetAmount,
		CreatedAt:        randomTimestamp,
	}

	goals := []GetCreatorGoalsResponse{goal}

	bytes, _ := json.Marshal(models.APIResponse{
		Data: goals,
	})
	w.Write(bytes)
}
