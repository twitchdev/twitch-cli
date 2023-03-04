// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package models

type EventsubSubscription struct {
	ID        string            `json:"id"`
	Status    string            `json:"status"`
	Type      string            `json:"type"`
	Version   string            `json:"version"`
	Condition EventsubCondition `json:"condition"`
	Transport EventsubTransport `json:"transport"`
	CreatedAt string            `json:"created_at"`
	Cost      int64             `json:"cost"`
}

type EventsubTransport struct {
	Method    string `json:"method"`
	Callback  string `json:"callback,omitempty"`
	SessionID string `json:"session_id,omitempty"`
}

type EventsubCondition struct {
	BroadcasterUserID     string `json:"broadcaster_user_id,omitempty"`
	ToBroadcasterUserID   string `json:"to_broadcaster_user_id,omitempty"`
	UserID                string `json:"user_id,omitempty"`
	FromBroadcasterUserID string `json:"from_broadcaster_user_id,omitempty"`
	ModeratorUserID       string `json:"moderator_user_id,omitempty"`
	ClientID              string `json:"client_id,omitempty"`
	ExtensionClientID     string `json:"extension_client_id,omitempty"`
	OrganizationID        string `json:"organization_id,omitempty"`
	CategoryID            string `json:"category_id,omitempty"`
	CampaignID            string `json:"campaign_id,omitempty"`
}

type EventsubResponse struct {
	Subscription EventsubSubscription `json:"subscription"`
	Event        interface{}          `json:"event"`
}

type EventsubSubscriptionVerification struct {
	Challenge    string               `json:"challenge"`
	Subscription EventsubSubscription `json:"subscription"`
}
