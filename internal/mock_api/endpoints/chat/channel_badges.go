// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package chat

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/twitchdev/twitch-cli/internal/database"
	"github.com/twitchdev/twitch-cli/internal/mock_api/mock_errors"
	"github.com/twitchdev/twitch-cli/internal/models"
	"github.com/twitchdev/twitch-cli/internal/util"
)

var channelBadgesMethodsSupported = map[string]bool{
	http.MethodGet:    true,
	http.MethodPost:   false,
	http.MethodDelete: false,
	http.MethodPatch:  false,
	http.MethodPut:    false,
}

var channelBadgesScopesByMethod = map[string][]string{
	http.MethodGet:    {},
	http.MethodPost:   {},
	http.MethodDelete: {},
	http.MethodPatch:  {},
	http.MethodPut:    {},
}

type ChannelBadges struct{}

func (e ChannelBadges) Path() string { return "/chat/badges" }

var defaultChannelSubscriberBadgeVersions = [][]string{
	{"0", "Subscriber"},            // Tier 1
	{"3", "3-Month Subscriber"},    // Tier 1 - 3 months
	{"6", "6-Month Subscriber"},    // Tier 1 - 6 months
	{"9", "9-Month Subscriber"},    // Tier 1 - 9 months
	{"12", "1-Year Subscriber"},    // Tier 1 - 12 months
	{"2000", "Subscriber"},         // Tier 2
	{"2003", "3-Month Subscriber"}, // Tier 2 - 3 months
	{"2006", "6-Month Subscriber"}, // Tier 2 - 6 months
	{"2009", "9-Month Subscriber"}, // Tier 2 - 9 months
	{"2012", "1-Year Subscriber"},  // Tier 2 - 12 months
	{"3000", "Subscriber"},         // Tier 3
	{"3003", "3-Month Subscriber"}, // Tier 3 - 3 months
	{"3006", "6-Month Subscriber"}, // Tier 3 - 6 months
	{"3009", "9-Month Subscriber"}, // Tier 3 - 9 months
	{"3012", "1-Year Subscriber"},  // Tier 3 - 12 months
}

func (e ChannelBadges) GetRequiredScopes(method string) []string {
	return channelBadgesScopesByMethod[method]
}

func (e ChannelBadges) ValidMethod(method string) bool {
	return channelBadgesMethodsSupported[method]
}

func (e ChannelBadges) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	db = r.Context().Value("db").(database.CLIDatabase)

	switch r.Method {
	case http.MethodGet:
		getChannelBadges(w, r)
		break
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
func getChannelBadges(w http.ResponseWriter, r *http.Request) {
	broadcasterID := r.URL.Query().Get("broadcaster_id")
	if broadcasterID == "" {
		mock_errors.WriteBadRequest(w, "Missing required parameter broadcaster_id")
		return
	}

	badges := []BadgesResponse{}

	b := BadgesResponse{
		SetID:    "subscriber",
		Versions: []BadgesVersion{},
	}

	for _, v := range defaultChannelSubscriberBadgeVersions {
		id := util.RandomGUID()
		b.Versions = append(b.Versions, BadgesVersion{
			ID:          v[0],
			ImageURL1X:  fmt.Sprintf("https://static-cdn.jtvnw.net/badges/v1/%v/1", id),
			ImageURL2X:  fmt.Sprintf("https://static-cdn.jtvnw.net/badges/v1/%v/2", id),
			ImageURL4X:  fmt.Sprintf("https://static-cdn.jtvnw.net/badges/v1/%v/3", id),
			Title:       v[1],
			Description: v[1],
			ClickAction: ptr("subscribe_to_channel"),
			ClickURL:    nil,
		})
	}
	badges = append(badges, b)

	bytes, _ := json.Marshal(models.APIResponse{Data: badges})
	w.Write(bytes)
}
