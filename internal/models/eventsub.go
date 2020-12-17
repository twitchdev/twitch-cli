// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package models

type EventsubSubscription struct {
	ID        string            `json:"id"`
	Type      string            `json:"type"`
	Version   string            `json:"version"`
	Condition EventsubCondition `json:"condition"`
	Transport EventsubTransport `json:"transport"`
	CreatedAt string            `json:"created_at"`
}

type EventsubTransport struct {
	Method   string `json:"method"`
	Callback string `json:"callback"`
}

type EventsubCondition struct {
	BroadcasterUserID string `json:"broadcaster_user_id"`
}

type EventsubResponse struct {
	Subscription EventsubSubscription `json:"subscription"`
	Event        interface{}          `json:"event"`
}
