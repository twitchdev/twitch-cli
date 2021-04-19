// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package models

type FollowEventSubEvent struct {
	UserID               string `json:"user_id"`
	UserLogin            string `json:"user_login"`
	UserName             string `json:"user_name"`
	BroadcasterUserID    string `json:"broadcaster_user_id"`
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
	BroadcasterUserName  string `json:"broadcaster_user_name"`,
        FollowedAt           string `json:"followed_at"`
}

type FollowWebSubResponse struct {
	Data []FollowWebSubResponseData `json:"data"`
}

type FollowWebSubResponseData struct {
	FromID     string `json:"from_id"`
	FromName   string `json:"from_name"`
	ToID       string `json:"to_id"`
	ToName     string `json:"to_name"`
	FollowedAt string `json:"followed_at"`
}

type FollowEventSubResponse struct {
	Subscription EventsubSubscription `json:"subscription"`
	Event        FollowEventSubEvent  `json:"event"`
}
