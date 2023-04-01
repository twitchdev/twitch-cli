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

var defaultGlobalSets = []BadgesResponse{
	{
		SetID: "clip-champ",
		Versions: []BadgesVersion{
			{
				ID:          "1",
				Title:       "Power Clipper",
				Description: "Power Clipper",
				ClickAction: ptr("visit_url"),
				ClickURL:    ptr("https://help.twitch.tv/customer/portal/articles/2918323-clip-champs-guide"),
			},
		},
	},
	{
		SetID: "extension",
		Versions: []BadgesVersion{
			{
				ID:          "1",
				Title:       "Extension",
				Description: "Extension",
				ClickAction: nil,
				ClickURL:    nil,
			},
		},
	},
	{
		SetID: "hype-train",
		Versions: []BadgesVersion{
			{
				ID:          "1",
				Title:       "Current Hype Train Conductor",
				Description: "Top supporter during the most recent hype train",
				ClickAction: ptr("visit_url"),
				ClickURL:    ptr("https://help.twitch.tv/s/article/hype-train-guide"),
			},
			{
				ID:          "2",
				Title:       "Former Hype Train Conductor",
				Description: "Top supporter during prior hype trains",
				ClickAction: ptr("visit_url"),
				ClickURL:    ptr("https://help.twitch.tv/s/article/hype-train-guide"),
			},
		},
	},
	{
		SetID: "moderator",
		Versions: []BadgesVersion{
			{
				ID:          "1",
				Title:       "Moderator",
				Description: "Moderator",
				ClickAction: nil,
				ClickURL:    nil,
			},
		},
	},
	{
		SetID: "premium",
		Versions: []BadgesVersion{
			{
				ID:          "1",
				Title:       "Prime Gaming",
				Description: "Prime Gaming",
				ClickAction: ptr("visit_url"),
				ClickURL:    ptr("https://gaming.amazon.com"),
			},
		},
	},
	{
		SetID: "vip",
		Versions: []BadgesVersion{
			{
				ID:          "1",
				Title:       "VIP",
				Description: "VIP",
				ClickAction: ptr("visit_url"),
				ClickURL:    ptr("https://help.twitch.tv/customer/en/portal/articles/659115-twitch-chat-badges-guide"),
			},
		},
	},
}

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
		versions := set.Versions
		for i := range versions {
			uuid := util.RandomGUID()
			versions[i].ImageURL1X = fmt.Sprintf("https://static-cdn.jtvnw.net/badges/v1/%v/1", uuid)
			versions[i].ImageURL2X = fmt.Sprintf("https://static-cdn.jtvnw.net/badges/v1/%v/2", uuid)
			versions[i].ImageURL4X = fmt.Sprintf("https://static-cdn.jtvnw.net/badges/v1/%v/3", uuid)
		}

		badges = append(badges, BadgesResponse{
			SetID:    set.SetID,
			Versions: versions,
		})
	}

	bytes, _ := json.Marshal(models.APIResponse{Data: badges})
	w.Write(bytes)
}
