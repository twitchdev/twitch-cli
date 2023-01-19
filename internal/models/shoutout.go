// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package models

// channel.shoutout.create

type ShoutoutCreateEventSubResponse struct {
	Subscription EventsubSubscription        `json:"subscription"`
	Event        ShoutoutCreateEventSubEvent `json:"event"`
}

type ShoutoutCreateEventSubEvent struct {
	BroadcasterUserID      string `json:"broadcaster_user_id"`
	BroadcasterUserName    string `json:"broadcaster_user_name"`
	BroadcasterUserLogin   string `json:"broadcaster_user_login"`
	ToBroadcasterUserID    string `json:"to_broadcaster_user_id"`
	ToBroadcasterUserName  string `json:"to_broadcaster_user_name"`
	ToBroadcasterUserLogin string `json:"to_broadcaster_user_login"`
	ModeratorUserID        string `json:"moderator_user_id"`
	ModeratorUserName      string `json:"moderator_user_name"`
	ModeratorUserLogin     string `json:"moderator_user_login"`
	ViewerCount            int    `json:"viewer_count"`
	StartedAt              string `json:"started_at"`
	CooldownEndsAt         string `json:"cooldown_ends_at"`
	TargetCooldownEndsAt   string `json:"target_cooldown_ends_at"`
}

// channel.shoutout.receive

type ShoutoutReceivedEventSubResponse struct {
	Subscription EventsubSubscription          `json:"subscription"`
	Event        ShoutoutReceivedEventSubEvent `json:"event"`
}

type ShoutoutReceivedEventSubEvent struct {
	BroadcasterUserID        string `json:"broadcaster_user_id"`
	BroadcasterUserName      string `json:"broadcaster_user_name"`
	BroadcasterUserLogin     string `json:"broadcaster_user_login"`
	FromBroadcasterUserID    string `json:"from_broadcaster_user_id"`
	FromBroadcasterUserName  string `json:"from_broadcaster_user_name"`
	FromBroadcasterUserLogin string `json:"from_broadcaster_user_login"`
	ViewerCount              int    `json:"viewer_count"`
	StartedAt                string `json:"started_at"`
}
