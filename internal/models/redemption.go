// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package models

type RedemptionEventSubEvent struct {
	Id                  string           `json:"id"`
	BroadcasterUserId   string           `json:"broadcaster_user_id"`
	BroadcasterUserName string           `json:"broadcaster_user_name"`
	UserId              string           `json:"user_id"`
	UserName            string           `json:"user_name"`
	UserInput           string           `json:"user_input"`
	Status              string           `json:"status"`
	Reward              RedemptionReward `json:"reward"`
	RedeemedAt          string           `json:"redeemed_at"`
}

type RedemptionReward struct {
	Id     string `json:"id"`
	Title  string `json:"title"`
	Cost   int64  `json:"cost"`
	Prompt string `json:"prompt"`
}
