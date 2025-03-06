// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package models

type MessagePartialEmote struct {
	ID         string `json:"id"`
	EmoteSetID string `json:"emote_set_id"`
}

type MessageCheermote struct {
	Prefix string `json:"prefix"`
	Bits   int    `json:"bits"`
	Tier   int    `json:"tier"`
}
