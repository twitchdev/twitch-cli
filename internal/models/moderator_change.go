// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package models

type ModeratorChangeWebSubEvent struct {
	ID             string                   `json:"id"`
	EventType      string                   `json:"event_type"`
	EventTimestamp string                   `json:"event_timestamp"`
	Version        string                   `json:"version"`
	EventData      ModeratorChangeEventData `json:"event_data"`
}

type ModeratorChangeEventData struct {
	BroadcasterID   string `json:"broadcaster_id"`
	BroadcasterName string `json:"broadcaster_name"`
	UserID          string `json:"user_id"`
	UserName        string `json:"user_name"`
}

type ModeratorChangeWebSubResponse struct {
	Data []ModeratorChangeWebSubEvent `json:"data"`
}
