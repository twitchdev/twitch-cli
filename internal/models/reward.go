// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package models

type RewardEventSubEvent struct {
	Id                                string               `json:"id"`
	BroadcasterUserId                 string               `json:"broadcaster_user_id"`
	BroadcasterUserName               string               `json:"broadcaster_user_name"`
	IsEnabled                         bool                 `json:"is_enabled"`
	IsPaused                          bool                 `json:"is_paused"`
	IsInStock                         bool                 `json:"is_in_stock"`
	Title                             string               `json:"title"`
	Cost                              int64                `json:"cost"`
	Prompt                            string               `json:"prompt"`
	IsUserInputRequired               bool                 `json:"is_user_input_required"`
	ShouldRedemptionsSkipRequestQueue bool                 `json:"should_redemptions_skip_request_queue"`
	CooldownExpiresAt                 string               `json:"cooldown_expires_at"`
	RedemptionsRedeemedCurrentStream  int64                `json:"redemptions_redeemed_current_stream"`
	MaxPerStream                      RewardMax            `json:"max_per_stream"`
	MaxPerUserPerStream               RewardMax            `json:"max_per_user_per_stream"`
	GlobalCooldown                    RewardGlobalCooldown `json:"global_cooldown"`
	BackgroundColor                   string               `json:"background_color"`
	Image                             RewardImage          `json:"image"`
	DefaultImage                      RewardImage          `json:"default_image"`
}

type RewardMax struct {
	IsEnabled bool  `json:"is_enabled"`
	Value     int64 `json:"value"`
}

type RewardGlobalCooldown struct {
	IsEnabled bool  `json:"is_enabled"`
	Value     int64 `json:"value"`
}

type RewardImage struct {
	Url1x string `json:"url_1x"`
	Url2x string `json:"url_2x"`
	Url4x string `json:"url_4x"`
}
