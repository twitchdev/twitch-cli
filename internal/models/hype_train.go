// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package models

type ContributionData struct {
	TotalContribution            int64  `json:"total"`
	TypeOfContribution           string `json:"type"`
	UserWhoMadeContribution      string `json:"user_id,omitempty"`
	UserNameWhoMadeContribution  string `json:"user_name,omitempty"`
	UserLoginWhoMadeContribution string `json:"user_login,omitempty"`
}

type HypeTrainEventSubResponse struct {
	Subscription EventsubSubscription   `json:"subscription"`
	Event        HypeTrainEventSubEvent `json:"event"`
}

type HypeTrainEventSubEvent struct {
	ID                      string             `json:"id"`
	BroadcasterUserID       string             `json:"broadcaster_user_id"`
	BroadcasterUserLogin    string             `json:"broadcaster_user_login"`
	BroadcasterUserName     string             `json:"broadcaster_user_name"`
	Level                   int64              `json:"level,omitempty"`
	Total                   int64              `json:"total"`
	Progress                *int64             `json:"progress,omitempty"`
	Goal                    int64              `json:"goal,omitempty"`
	TopContributions        []ContributionData `json:"top_contributions"`
	LastContribution        ContributionData   `json:"last_contribution,omitempty"`
	StartedAtTimestamp      string             `json:"started_at,omitempty"`
	ExpiresAtTimestamp      string             `json:"expires_at,omitempty"`
	EndedAtTimestamp        string             `json:"ended_at,omitempty"`
	CooldownEndsAtTimestamp string             `json:"cooldown_ends_at,omitempty"`
}

type SharedTrainParticipant struct {
	BroadcasterUserID    string `json:"broadcaster_user_id"`
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
	BroadcasterUserName  string `json:"broadcaster_user_name"`
}

type HypeTrainEventSubEventV2 struct {
	ID                      string                   `json:"id"`
	BroadcasterUserID       string                   `json:"broadcaster_user_id"`
	BroadcasterUserLogin    string                   `json:"broadcaster_user_login"`
	BroadcasterUserName     string                   `json:"broadcaster_user_name"`
	Total                   int64                    `json:"total"`
	Progress                *int64                   `json:"progress,omitempty"`
	Goal                    int64                    `json:"goal,omitempty"`
	TopContributions        []ContributionData       `json:"top_contributions"`
	SharedTrainParticipants *[]SharedTrainParticipant `json:"shared_train_participants"`
	Level                   int64                    `json:"level,omitempty"`
	StartedAtTimestamp      string                   `json:"started_at,omitempty"`
	ExpiresAtTimestamp      string                   `json:"expires_at,omitempty"`
	EndedAtTimestamp        string                   `json:"ended_at,omitempty"`
	CooldownEndsAtTimestamp string                   `json:"cooldown_ends_at,omitempty"`
	IsSharedTrain           bool                     `json:"is_shared_train"`
	Type                    string                   `json:"type"`
	AllTimeHighLevel        int64                    `json:"all_time_high_level,omitempty"`
	AllTimeHighTotal        int64                    `json:"all_time_high_total,omitempty"`
}

type HypeTrainEventSubResponseV2 struct {
	Subscription EventsubSubscription     `json:"subscription"`
	Event        HypeTrainEventSubEventV2 `json:"event"`
}
