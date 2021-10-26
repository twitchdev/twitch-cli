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
	Progress                int64              `json:"progress,omitempty"`
	Goal                    int64              `json:"goal,omitempty"`
	TopContributions        []ContributionData `json:"top_contributions"`
	LastContribution        ContributionData   `json:"last_contribution,omitempty"`
	StartedAtTimestamp      string             `json:"started_at,omitempty"`
	ExpiresAtTimestamp      string             `json:"expires_at,omitempty"`
	EndedAtTimestamp        string             `json:"ended_at,omitempty"`
	CooldownEndsAtTimestamp string             `json:"cooldown_ends_at,omitempty"`
}
