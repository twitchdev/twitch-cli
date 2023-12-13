// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package charity

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"

	"github.com/twitchdev/twitch-cli/internal/database"
	"github.com/twitchdev/twitch-cli/internal/mock_api/authentication"
	"github.com/twitchdev/twitch-cli/internal/mock_api/mock_errors"
	"github.com/twitchdev/twitch-cli/internal/models"
	"github.com/twitchdev/twitch-cli/internal/util"
)

var campaignsMethodsSupported = map[string]bool{
	http.MethodGet:    true,
	http.MethodPost:   false,
	http.MethodDelete: false,
	http.MethodPatch:  false,
	http.MethodPut:    false,
}

var campaignsScopesByMethod = map[string][]string{
	http.MethodGet:    {"channel:read:charity"},
	http.MethodPost:   {},
	http.MethodDelete: {},
	http.MethodPatch:  {},
	http.MethodPut:    {},
}

type CharityCampaign struct{}

type GetCharityCampaignResponse struct {
	ID                 string        `json:"id"`
	BroadcasterID      string        `json:"broadcaster_id"`
	BroadcasterName    string        `json:"broadcaster_name"`
	BroadcasterLogin   string        `json:"broadcaster_login"`
	CharityName        string        `json:"charity_name"`
	CharityDescription string        `json:"charity_description"`
	CharityLogo        string        `json:"charity_logo"`
	CharityWebsite     string        `json:"charity_website"`
	CurrentAmount      CharityAmount `json:"current_amount"`
	TargetAmount       CharityAmount `json:"target_amount"`
}

func (e CharityCampaign) Path() string { return "/charity/campaigns" }

func (e CharityCampaign) GetRequiredScopes(method string) []string {
	return campaignsScopesByMethod[method]
}

func (e CharityCampaign) ValidMethod(method string) bool {
	return campaignsMethodsSupported[method]
}

func (e CharityCampaign) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	db = r.Context().Value("db").(database.CLIDatabase)
	switch r.Method {
	case http.MethodGet:
		getCharityCampaign(w, r)
		break
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func getCharityCampaign(w http.ResponseWriter, r *http.Request) {
	userCtx := r.Context().Value("auth").(authentication.UserAuthentication)
	if !userCtx.MatchesBroadcasterIDParam(r) {
		mock_errors.WriteUnauthorized(w, "Broadcaster ID does not match token.")
		return
	}

	broadcasterID := r.URL.Query().Get("broadcaster_id")

	if broadcasterID == "" {
		mock_errors.WriteBadRequest(w, "one of broadcaster_id, game_id, or id is required")
		return
	}

	user, err := db.NewQuery(r, 100).GetUser(database.User{ID: userCtx.UserID})
	if err != nil {
		mock_errors.WriteServerError(w, err.Error())
		return
	}

	// Generate current and target amounts
	target := rand.Intn((10000*100)-(1000*100)) + (1000 * 100) // Between $1,000 and $10,000
	current := rand.Intn(target-500) + 500 + 500               // Between $500 and target+$500

	// Currently there's no interactivity between the Mock API and mock EventSub, and no API to add additional value, so this will return random values much like the Hype Train API
	charityCampaign := GetCharityCampaignResponse{
		ID:                 util.RandomGUID(),
		BroadcasterID:      userCtx.UserID,
		BroadcasterName:    user.DisplayName,
		BroadcasterLogin:   user.UserLogin,
		CharityName:        fmt.Sprintf("%v's Amazing Charity #%v!", user.DisplayName, util.RandomInt(1000)),
		CharityDescription: "Yet another amazing charity! PogChamp",
		CharityLogo:        "https://abc.cloudfront.net/ppgf/1000/100.png",
		CharityWebsite:     "https://www.example.com",
		CurrentAmount: CharityAmount{
			Value:         current,
			DecimalPlaces: 2,
			Currency:      "USD",
		},
		TargetAmount: CharityAmount{
			Value:         target,
			DecimalPlaces: 2,
			Currency:      "USD",
		},
	}

	campaigns := []GetCharityCampaignResponse{charityCampaign}

	apiResponse := models.APIResponse{
		Data: campaigns,
	}

	bytes, _ := json.Marshal(apiResponse)
	w.Write(bytes)
}
