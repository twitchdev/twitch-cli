// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package models

type SubEventSubResponse struct {
	Subscription EventsubSubscription `json:"subscription"`
	Event        SubEventSubEvent     `json:"event"`
}

type SubEventSubEvent struct {
	UserID               string `json:"user_id"`
	UserLogin            string `json:"user_login"`
	UserName             string `json:"user_name"`
	BroadcasterUserID    string `json:"broadcaster_user_id"`
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
	BroadcasterUserName  string `json:"broadcaster_user_name"`
	Tier                 string `json:"tier"`
	IsGift               bool   `json:"is_gift"`
}

type GiftEventSubResponse struct {
	Subscription EventsubSubscription `json:"subscription"`
	Event        GiftEventSubEvent    `json:"event"`
}

type GiftEventSubEvent struct {
	UserID               string `json:"user_id"`
	UserLogin            string `json:"user_login"`
	UserName             string `json:"user_name"`
	BroadcasterUserID    string `json:"broadcaster_user_id"`
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
	BroadcasterUserName  string `json:"broadcaster_user_name"`
	Tier                 string `json:"tier"`
	Total                int    `json:"total"`
	IsAnonymous          bool   `json:"is_anonymous"`
	CumulativeTotal      *int   `json:"cumulative_total"`
}

type SubscribeMessageEventSubResponse struct {
	Subscription EventsubSubscription          `json:"subscription"`
	Event        SubscribeMessageEventSubEvent `json:"event"`
}

type SubscribeMessageEventSubEvent struct {
	UserID               string                          `json:"user_id"`
	UserLogin            string                          `json:"user_login"`
	UserName             string                          `json:"user_name"`
	BroadcasterUserID    string                          `json:"broadcaster_user_id"`
	BroadcasterUserLogin string                          `json:"broadcaster_user_login"`
	BroadcasterUserName  string                          `json:"broadcaster_user_name"`
	Tier                 string                          `json:"tier"`
	Message              SubscribeMessageEventSubMessage `json:"message"`
	CumulativeMonths     int                             `json:"cumulative_months"`
	StreakMonths         *int                            `json:"streak_months"`
	DurationMonths       int                             `json:"duration_months"`
}

type SubscribeMessageEventSubMessage struct {
	Text   string                                 `json:"text"`
	Emotes []SubscribeMessageEventSubMessageEmote `json:"emotes"`
}

type SubscribeMessageEventSubMessageEmote struct {
	Begin int    `json:"begin"`
	End   int    `json:"end"`
	ID    string `json:"id"`
}
