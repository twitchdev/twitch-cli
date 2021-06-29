// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package streams

import "github.com/twitchdev/twitch-cli/internal/database"

var db database.CLIDatabase

type TagResponse struct {
	TagID                    string         `json:"tag_id"`
	IsAuto                   bool           `json:"is_auto"`
	LocalizationNames        []Localization `json:"localization_names"`
	LocalizationDescriptions []Localization `json:"localization_descriptions"`
}

type Localization struct {
	EnglishUS string `json:"en-us"`
}

func convertTags(tags []database.Tag) []TagResponse {
	response := []TagResponse{}

	for _, tag := range tags {
		t := TagResponse{
			TagID:                    tag.ID,
			IsAuto:                   tag.IsAuto,
			LocalizationNames:        []Localization{{EnglishUS: tag.Name}},
			LocalizationDescriptions: []Localization{{EnglishUS: tag.Name}},
		}

		response = append(response, t)
	}
	return response
}
