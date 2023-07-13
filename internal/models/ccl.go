// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package models

type ContentClassificationLabel struct {
	Description      string `json:"description"`
	ID               string `json:"id"`
	Name             string `json:"name"`
	RestrictedGaming bool   `json:"-"` // Restricts users from applying that CCL via the API. Currently only for MatureGame.
}

var CCL_MAP = map[string]ContentClassificationLabel{
	"DrugsIntoxication": {
		Description:      "Excessive tobacco glorification or promotion, any marijuana consumption/use, legal drug and alcohol induced intoxication, discussions of illegal drugs.",
		ID:               "DrugsIntoxication",
		Name:             "Drugs, Intoxication, or Excessive Tobacco Use",
		RestrictedGaming: false,
	},
	"Gambling": {
		Description:      "Participating in online or in-person gambling, poker or fantasy sports, that involve the exchange of real money.",
		ID:               "Gambling",
		Name:             "Gambling",
		RestrictedGaming: false,
	},
	"MatureGame": {
		Description:      "Games that are rated Mature or less suitable for a younger audience.",
		ID:               "MatureGame",
		Name:             "Mature-rated game",
		RestrictedGaming: true,
	},
	"ProfanityVulgarity": {
		Description:      "Prolonged, and repeated use of obscenities, profanities, and vulgarities, especially as a regular part of speech.",
		ID:               "ProfanityVulgarity",
		Name:             "Significant Profanity or Vulgarity",
		RestrictedGaming: false,
	},
	"SexualThemes": {
		Description:      "Content that focuses on sexualized physical attributes and activities, sexual topics, or experiences.",
		ID:               "SexualThemes",
		Name:             "Sexual Themes",
		RestrictedGaming: false,
	},
	"ViolentGraphic": {
		Description:      "Simulations and/or depictions of realistic violence, gore, extreme injury, or death.",
		ID:               "ViolentGraphic",
		Name:             "Violent and Graphic Depictions",
		RestrictedGaming: false,
	},
}
