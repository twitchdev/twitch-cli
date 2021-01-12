// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package models

type RedemptionEventSubEvent struct {
	ID                   string           `json:"id"`
	BroadcasterUserID    string           `json:"broadcaster_user_id"`
	BroadcasterUserLogin string           `json:"broadcaster_user_login"`
	BroadcasterUserName  string           `json:"broadcaster_user_name"`
	UserID               string           `json:"user_id"`
	UserLogin            string           `json:"user_login"`
	UserName             string           `json:"user_name"`
	UserInput            string           `json:"user_input"`
	Status               string           `json:"status"`
	Reward               RedemptionReward `json:"reward"`
	RedeemedAt           string           `json:"redeemed_at"`
}

type RedemptionReward struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Cost   int64  `json:"cost"`
	Prompt string `json:"prompt"`
}

type RedemptionEventSubResponse struct {
	Subscription EventsubSubscription    `json:"subscription"`
	Event        RedemptionEventSubEvent `json:"event"`
}
