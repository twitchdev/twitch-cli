// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package channels

import (
	"encoding/json"
	"math"
	"net/http"
	"time"

	"github.com/twitchdev/twitch-cli/internal/database"
	"github.com/twitchdev/twitch-cli/internal/mock_api/authentication"
	"github.com/twitchdev/twitch-cli/internal/mock_api/mock_errors"
	"github.com/twitchdev/twitch-cli/internal/models"
)

var commercialMethodsSupported = map[string]bool{
	http.MethodGet:    false,
	http.MethodPost:   true,
	http.MethodDelete: false,
	http.MethodPatch:  false,
	http.MethodPut:    false,
}

var commercialScopesByMethod = map[string][]string{
	http.MethodGet:    {},
	http.MethodPost:   {"channel:edit:commercial"},
	http.MethodDelete: {},
	http.MethodPatch:  {},
	http.MethodPut:    {},
}

type CommercialEndpoint struct{}

type CommercialEndpointRequest struct {
	Length        *int   `json:"length"`
	BroadcasterID string `json:"broadcaster_id"`
}

type CommercialEndpointResponse struct {
	Length     int    `json:"length"`
	Message    string `json:"message"`
	RetryAfter int    `json:"retry_after"`
}

// maps broadcaster ID to a UTC time to determine any cooldowns.
var commericalCooldown = map[string]time.Time{}

func (e CommercialEndpoint) Path() string { return "/channels/commercial" }

func (e CommercialEndpoint) GetRequiredScopes(method string) []string {
	return commercialScopesByMethod[method]
}

func (e CommercialEndpoint) ValidMethod(method string) bool {
	return commercialMethodsSupported[method]
}

func (e CommercialEndpoint) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	db = r.Context().Value("db").(database.CLIDatabase)

	switch r.Method {
	case http.MethodPost:
		postCommercial(w, r)
		break
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		break
	}
}

func postCommercial(w http.ResponseWriter, r *http.Request) {
	var body CommercialEndpointRequest
	userCtx := r.Context().Value("auth").(authentication.UserAuthentication)

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		w.Write(mock_errors.GetErrorBytes(http.StatusInternalServerError, err, "Error unmarshaling body"))
		return
	}

	if body.BroadcasterID == "" || body.Length == nil {
		mock_errors.WriteBadRequest(w, "Body missing one of broadcaster_id or length")
		return
	}

	if body.BroadcasterID != userCtx.UserID {
		mock_errors.WriteUnauthorized(w, "broadcaster_id does not match token")
		return
	}

	if validLength(body.Length) == false {
		mock_errors.WriteBadRequest(w, "Length is an invalid value")
		return
	}

	defaultRetryLength := 480

	commericalResponse := []CommercialEndpointResponse{
		{
			Length:     *body.Length,
			Message:    "",
			RetryAfter: defaultRetryLength,
		},
	}

	s, err := db.NewQuery(r, 1).GetStream(database.Stream{UserID: body.BroadcasterID})
	if err != nil {
		mock_errors.WriteServerError(w, "Error fetching stream status")
		return
	}

	// if offline, return error
	if len(s.Data.([]database.Stream)) == 0 {
		mock_errors.WriteBadRequest(w, "User is not currently live and must be to run commercials")
		return
	}
	// if there's an entry, validate it and respond with an error if within cooldown window
	if commericalCooldown[body.BroadcasterID].IsZero() == false {
		retryAfter := math.Round(commericalCooldown[body.BroadcasterID].UTC().Sub(time.Now().UTC()).Seconds())

		// still in cooldown, so return message stating as such
		if retryAfter > 0 {
			commericalResponse[0] = CommercialEndpointResponse{
				Length:     0,
				Message:    "Please try again later",
				RetryAfter: int(retryAfter),
			}
		} else {
			// after cooldown, so remvoe entry and continue
			delete(commericalCooldown, body.BroadcasterID)
		}

	} else {
		// otherwise, add user to cooldown tracker
		commericalCooldown[body.BroadcasterID] = time.Now().UTC().Add(time.Second * time.Duration(defaultRetryLength))
	}

	apiResponse := models.APIResponse{
		Data: commericalResponse,
	}

	bytes, _ := json.Marshal(apiResponse)
	w.Write(bytes)
}

func validLength(length *int) bool {
	possibleLengths := []int{30, 60, 90, 120, 150, 180}
	for _, l := range possibleLengths {
		if *length == l {
			return true
		}
	}

	return false
}
