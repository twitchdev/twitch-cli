// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package chat

import "github.com/twitchdev/twitch-cli/internal/database"

var db database.CLIDatabase

type BadgesResponse struct {
	SetID    string          `json:"set_id"`
	Versions []BadgesVersion `json:"versions"`
}

type BadgesVersion struct {
	ID         string `json:"id"`
	ImageURL1X string `json:"image_url_1x"`
	ImageURL2X string `json:"image_url_2x"`
	ImageURL4X string `json:"image_url_4x"`
}
