// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package models

type ChannelUpdateEventSubEvent struct {
	BroadcasterUserID    string `json:"broadcaster_user_id"`
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
	BroadcasterUserName  string `json:"broadcaster_user_name"`
	StreamTitle          string `json:"title"`
	StreamLanguage       string `json:"language"`
	StreamCategoryID     string `json:"category_id"`
	StreamCategoryName   string `json:"category_name"`
	IsMature             bool   `json:"is_mature"`
}

type StreamChangeWebSubResponse struct {
	Data []StreamChangeWebSubResponseData `json:"data"`
}

type StreamChangeWebSubResponseData struct {
	WebsubID             string   `json:"id"`
	BroadcasterUserID    string   `json:"user_id"`
	BroadcasterUserLogin string   `json:"user_login"`
	BroadcasterUserName  string   `json:"user_name"`
	StreamCategoryID     string   `json:"game_id"`
	StreamCategoryName   string   `json:"game_name"`
	StreamType           string   `json:"type"`
	StreamTitle          string   `json:"title"`
	StreamViewerCount    int      `json:"viewer_count"`
	StreamStartedAt      string   `json:"started_at"`
	StreamLanguage       string   `json:"language"`
	StreamThumbnailURL   string   `json:"thumbnail_url"`
	TagIDs               []string `json:"tag_ids"`
}

type ChannelUpdateEventSubResponse struct {
	Subscription EventsubSubscription       `json:"subscription"`
	Event        ChannelUpdateEventSubEvent `json:"event"`
}
