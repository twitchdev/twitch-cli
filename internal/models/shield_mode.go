// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package models

type ShieldModeEventSubResponse struct {
	Subscription EventsubSubscription    `json:"subscription"`
	Event        ShieldModeEventSubEvent `json:"event"`
}

type ShieldModeEventSubEvent struct {
	BroadcasterUserID    string `json:"broadcaster_user_id"`
	BroadcasterUserName  string `json:"broadcaster_user_name"`
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
	ModeratorUserID      string `json:"moderator_user_id"`
	ModeratorUserName    string `json:"moderator_user_name"`
	ModeratorUserLogin   string `json:"moderator_user_login"`
	StartedAt            string `json:"started_at,omitempty"`
	EndedAt              string `json:"ended_at,omitempty"`
}
