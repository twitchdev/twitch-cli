// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package chat

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/twitchdev/twitch-cli/internal/database"
	"github.com/twitchdev/twitch-cli/internal/models"
	"github.com/twitchdev/twitch-cli/internal/util"
)

var globalBadgesMethodsSupported = map[string]bool{
	http.MethodGet:    true,
	http.MethodPost:   false,
	http.MethodDelete: false,
	http.MethodPatch:  false,
	http.MethodPut:    false,
}

var globalBadgesScopesByMethod = map[string][]string{
	http.MethodGet:    {},
	http.MethodPost:   {},
	http.MethodDelete: {},
	http.MethodPatch:  {},
	http.MethodPut:    {},
}

type GlobalBadges struct{}

var defaultGlobalSets = []string{"clip-champ", "extension", "hype-train", "moderator", "premium", "vip"}

func (e GlobalBadges) Path() string { return "/chat/badges/global" }

func (e GlobalBadges) GetRequiredScopes(method string) []string {
	return globalBadgesScopesByMethod[method]
}

func (e GlobalBadges) ValidMethod(method string) bool {
	return globalBadgesMethodsSupported[method]
}

func (e GlobalBadges) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	db = r.Context().Value("db").(database.CLIDatabase)

	switch r.Method {
	case http.MethodGet:
		getGlobalBadges(w, r)
		break
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func getGlobalBadges(w http.ResponseWriter, r *http.Request) {
	badges := []BadgesResponse{}

	for _, set := range defaultGlobalSets {
		id := util.RandomGUID()
		badges = append(badges, BadgesResponse{
			SetID: set,
			Versions: []BadgesVersion{
				{
					ID:         "1",
					ImageURL1X: fmt.Sprintf("https://static-cdn.jtvnw.net/badges/v1/%v/1", id),
					ImageURL2X: fmt.Sprintf("https://static-cdn.jtvnw.net/badges/v1/%v/2", id),
					ImageURL4X: fmt.Sprintf("https://static-cdn.jtvnw.net/badges/v1/%v/3", id),
				},
			},
		})
	}

	bytes, _ := json.Marshal(models.APIResponse{Data: badges})
	w.Write(bytes)
}
