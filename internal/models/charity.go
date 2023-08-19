// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package models

type CharityEventSubEventAmount struct {
	Value         int    `json:"value"`
	DecimalPlaces int    `json:"decimal_places"`
	Currency      string `json:"currency"`
}

type CharityEventSubEvent struct {
	CampaignID           *string                     `json:"campaign_id,omitempty"` // Specific to channel.charity_campaign.donate
	ID                   string                      `json:"id,omitempty"`          // Used by everything else under channel.charity_campaign.*
	BroadcasterUserID    string                      `json:"broadcaster_user_id"`
	BroadcasterUserName  string                      `json:"broadcaster_user_name"`
	BroadcasterUserLogin string                      `json:"broadcaster_user_login"`
	UserID               *string                     `json:"user_id,omitempty"`
	UserName             *string                     `json:"user_name,omitempty"`
	UserLogin            *string                     `json:"user_login,omitempty"`
	CharityName          string                      `json:"charity_name"`
	CharityDescription   string                      `json:"charity_description,omitempty"`
	CharityLogo          string                      `json:"charity_logo"`
	CharityWebsite       string                      `json:"charity_website,omitempty"`
	Amount               *CharityEventSubEventAmount `json:"amount,omitempty"`
	CurrentAmount        *CharityEventSubEventAmount `json:"current_amount,omitempty"`
	TargetAmount         *CharityEventSubEventAmount `json:"target_amount,omitempty"`
	StartedAt            *string                     `json:"started_at,omitempty"`
	StoppedAt            *string                     `json:"stopped_at,omitempty"`
}

type CharityEventSubResponse struct {
	Subscription EventsubSubscription `json:"subscription"`
	Event        CharityEventSubEvent `json:"event"`
}
