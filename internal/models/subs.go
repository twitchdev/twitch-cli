// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package models

type SubWebSubEventData struct {
	BroadcasterID   string `json:"broadcaster_id"`
	BroadcasterName string `json:"broadcaster_name"`
	IsGift          bool   `json:"is_gift"`
	Tier            string `json:"tier"`
	PlanName        string `json:"plan_name"`
	UserID          string `json:"user_id"`
	UserName        string `json:"user_name"`
	GifterID        string `json:"gifter_id"`
	GifterName      string `json:"gifter_name"`
}

type SubWebSubResponse struct {
	Data []SubWebSubResponseData `json:"data"`
}

type SubWebSubResponseData struct {
	ID             string             `json:"id"`
	EventType      string             `json:"event_type"`
	EventTimestamp string             `json:"event_timestamp"`
	Version        string             `json:"version"`
	EventData      SubWebSubEventData `json:"event_data"`
}

type SubEventSubResponse struct {
	Subscription EventsubSubscription `json:"subscription"`
	Event        SubEventSubEvent     `json:"event"`
}

type SubEventSubEvent struct {
	UserID               string `json:"user_id"`
	UserLogin            string `json:"user_login"`
	UserName             string `json:"user_name"`
	BroadcasterUserID    string `json:"broadcaster_user_id"`
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
	BroadcasterUserName  string `json:"broadcaster_user_name"`
	Tier                 string `json:"tier"`
	IsGift               bool   `json:"is_gift"`
}
