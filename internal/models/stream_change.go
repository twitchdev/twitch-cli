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

type ChannelUpdateEventSubResponse struct {
	Subscription EventsubSubscription       `json:"subscription"`
	Event        ChannelUpdateEventSubEvent `json:"event"`
}
