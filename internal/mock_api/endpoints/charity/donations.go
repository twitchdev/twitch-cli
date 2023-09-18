// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package charity

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/twitchdev/twitch-cli/internal/database"
	"github.com/twitchdev/twitch-cli/internal/mock_api/authentication"
	"github.com/twitchdev/twitch-cli/internal/mock_api/mock_errors"
	"github.com/twitchdev/twitch-cli/internal/models"
	"github.com/twitchdev/twitch-cli/internal/util"
)

var donationsMethodsSupported = map[string]bool{
	http.MethodGet:    true,
	http.MethodPost:   false,
	http.MethodDelete: false,
	http.MethodPatch:  false,
	http.MethodPut:    false,
}

var donationsScopesByMethod = map[string][]string{
	http.MethodGet:    {"channel:read:charity"},
	http.MethodPost:   {},
	http.MethodDelete: {},
	http.MethodPatch:  {},
	http.MethodPut:    {},
}

type CharityDonations struct{}

type GetCharityDonationsResponse struct {
	ID        		string        `json:"campaign_id"`
	UserID    		string        `json:"user_id"`
	UserLogin 		string        `json:"user_login"`
	UserName  		string        `json:"user_name"`
	TargetAmount 	CharityAmount `json:"target_amount"`
}

func (e CharityDonations) Path() string { return "/charity/donations" }

func (e CharityDonations) GetRequiredScopes(method string) []string {
	return donationsScopesByMethod[method]
}

func (e CharityDonations) ValidMethod(method string) bool {
	return donationsMethodsSupported[method]
}

func (e CharityDonations) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	db = r.Context().Value("db").(database.CLIDatabase)
	switch r.Method {
	case http.MethodGet:
		getCharityDonations(w, r)
		break
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func getCharityDonations(w http.ResponseWriter, r *http.Request) {
	userCtx := r.Context().Value("auth").(authentication.UserAuthentication)
	first := 20
	if !userCtx.MatchesBroadcasterIDParam(r) {
		mock_errors.WriteUnauthorized(w, "Broadcaster ID does not match token.")
		return
	}

	broadcasterID := r.URL.Query().Get("broadcaster_id")

	if broadcasterID == "" {
		mock_errors.WriteBadRequest(w, "one of broadcaster_id, game_id, or id is required")
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

	user, err := db.NewQuery(r, 100).GetUser(database.User{ID: userCtx.UserID})
	if err != nil {
		mock_errors.WriteServerError(w, err.Error())
		return
	}

	// Currently there's no interactivity between the Mock API and mock EventSub, and no API to add additional value, so this will return random values much like the Hype Train API
	donations := []GetCharityDonationsResponse{}

	for i := 0; i < first; i++ {
		d := GetCharityDonationsResponse{
			ID:        util.RandomGUID(),
			UserID:    userCtx.UserID,
			UserName:  user.DisplayName,
			UserLogin: user.UserLogin,
			Amount: CharityAmount{
				Value:         rand.Intn(150000-300) + 300, // Between $3 and $1,500
				DecimalPlaces: 2,
				Currency:      "USD",
			},
		}

		donations = append(donations, d)
	}

	apiResponse := models.APIResponse{
		Data:       donations,
		Pagination: &models.APIPagination{}, // Since these are randomly generated, true pagination isn't implemented
	}

	bytes, _ := json.Marshal(apiResponse)
	w.Write(bytes)
}
