// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package chat

import "github.com/twitchdev/twitch-cli/internal/database"

var db database.CLIDatabase
var defaultEmoteTypes = []string{"subscription", "bitstier", "follower"}

const templateEmoteURL = "https://static-cdn.jtvnw.net/emoticons/v2/{{id}}/{{format}}/{{theme_mode}}/{{scale}}"

type BadgesResponse struct {
	SetID    string          `json:"set_id"`
	Versions []BadgesVersion `json:"versions"`
}

type BadgesVersion struct {
	ID          string  `json:"id"`
	ImageURL1X  string  `json:"image_url_1x"`
	ImageURL2X  string  `json:"image_url_2x"`
	ImageURL4X  string  `json:"image_url_4x"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	ClickAction *string `json:"click_action"`
	ClickURL    *string `json:"click_url"`
}

type EmotesResponse struct {
	ID         string       `json:"id"`
	Name       string       `json:"name"`
	Images     EmotesImages `json:"images"`
	Tier       *string      `json:"tier,omitempty"`
	EmoteType  *string      `json:"emote_type,omitempty"`
	EmoteSetID *string      `json:"emote_set_id,omitempty"`
	OwnerID    *string      `json:"owner_id,omitempty"`
	Format     []string     `json:"format"`
	Scale      []string     `json:"scale"`
	ThemeMode  []string     `json:"theme_mode"`
}

type EmotesImages struct {
	ImageURL1X string `json:"url_1x"`
	ImageURL2X string `json:"url_2x"`
	ImageURL4X string `json:"url_4x"`
}

func ptr(str string) *string {
	return &str
}
