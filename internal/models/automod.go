// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package models

type AutomodMessageHoldEvent struct {
	BroadcasterUserID    string                           `json:"broadcaster_user_id"`
	BroadcasterUserLogin string                           `json:"broadcaster_user_login"`
	BroadcasterUserName  string                           `json:"broadcaster_user_name"`
	UserID               string                           `json:"user_id"`
	UserLogin            string                           `json:"user_login"`
	UserName             string                           `json:"user_name"`
	MessageID            string                           `json:"message_id"`
	Message              AutomodMessage                   `json:"message"`
	HeldAt               string                           `json:"held_at"`
	Reason               string                           `json:"reason"`
	Automod              *AutomodMessageAutomodReason     `json:"automod"`
	BlockedTerm          *AutomodMessageBlockedTermReason `json:"blocked_term"`
}

type AutomodMessageHoldEventSubResponse struct {
	Subscription EventsubSubscription    `json:"subscription"`
	Event        AutomodMessageHoldEvent `json:"event"`
}

type AutomodMessageBoundary struct {
	StartPos int `json:"start_pos"`
	EndPos   int `json:"end_pos"`
}

type AutomodMessageAutomodReason struct {
	Category   string                   `json:"category"`
	Level      int                      `json:"level"`
	Boundaries []AutomodMessageBoundary `json:"boundaries"`
}

type AutomodMessageFoundTerm struct {
	TermID                    string                 `json:"term_id"`
	Boundary                  AutomodMessageBoundary `json:"boundary"`
	OwnerBroadcasterUserID    string                 `json:"owner_broadcaster_user_id"`
	OwnerBroadcasterUserLogin string                 `json:"owner_broadcaster_user_login"`
	OwnerBroadcasterUserName  string                 `json:"owner_broadcaster_user_name"`
}

type AutomodMessageBlockedTermReason struct {
	TermsFound []AutomodMessageFoundTerm `json:"terms_found"`
}

type AutomodMessage struct {
	Text      string                   `json:"text"`
	Fragments []AutomodMessageFragment `json:"fragments"`
}

type AutomodMessageFragment struct {
	Type      string               `json:"type"`
	Text      string               `json:"text"`
	Emote     *MessagePartialEmote `json:"emote"`
	Cheermote *MessageCheermote    `json:"cheermote"`
}
