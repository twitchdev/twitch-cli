// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package models

type ContributionData struct {
	TotalContribution            int64  `json:"total"`
	TypeOfContribution           string `json:"type"`
	WebSubUser                   string `json:"user,omitempty"`
	UserWhoMadeContribution      string `json:"user_id,omitempty"`
	UserNameWhoMadeContribution  string `json:"user_name,omitempty"`
	UserLoginWhoMadeContribution string `json:"user_login,omitempty"`
}

type HypeTrainWebSubEvent struct {
	ID             string                   `json:"id"`
	EventType      string                   `json:"event_type"`
	EventTimestamp string                   `json:"event_timestamp"`
	Version        string                   `json:"version"`
	EventData      HypeTrainWebsubEventData `json:"event_data"`
}

type HypeTrainWebsubEventData struct {
	BroadcasterID        string             `json:"broadcaster_id"`
	CooldownEndTimestamp string             `json:"cooldown_end_time"`
	ExpiresAtTimestamp   string             `json:"expires_at"`
	Goal                 int64              `json:"goal,omitempty"`
	Id                   string             `json:"id,omitempty"`
	LastContribution     ContributionData   `json:"last_contribution,omitempty"`
	Level                int64              `json:"level,omitempty"`
	StartedAtTimestamp   string             `json:"started_at,omitempty"`
	TopContributions     []ContributionData `json:"top_contributions"`
	Total                int64              `json:"total"`
}

type HypeTrainWebSubResponse struct {
	Data []HypeTrainWebSubEvent `json:"data"`
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
	Progress                int64              `json:"progress,omitempty"`
	Goal                    int64              `json:"goal,omitempty"`
	TopContributions        []ContributionData `json:"top_contributions"`
	LastContribution        ContributionData   `json:"last_contribution,omitempty"`
	StartedAtTimestamp      string             `json:"started_at,omitempty"`
	ExpiresAtTimestamp      string             `json:"expires_at,omitempty"`
	EndedAtTimestamp        string             `json:"ended_at,omitempty"`
	CooldownEndsAtTimestamp string             `json:"cooldown_ends_at,omitempty"`
}
