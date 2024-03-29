// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package models

type UnbanRequestCreateEventSubEvent struct {
	BroadcasterUserID    string `json:"broadcaster_user_id"`
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
	BroadcasterUserName  string `json:"broadcaster_user_name"`
	UserID               string `json:"user_id"`
	UserLogin            string `json:"user_login"`
	UserName             string `json:"user_name"`
	Text                 string `json:"text"`
	CreatedAt            string `json:"created_at"`
}

type UnbanRequestCreateEventSubResponse struct {
	Subscription EventsubSubscription     `json:"subscription"`
	Event        TransactionEventSubEvent `json:"event"`
}

type UnbanRequestResolveEventSubEvent struct {
	ID                   string  `json:"id"`
	BroadcasterUserID    string  `json:"broadcaster_user_id"`
	BroadcasterUserLogin string  `json:"broadcaster_user_login"`
	BroadcasterUserName  string  `json:"broadcaster_user_name"`
	ModeratorUserID      *string `json:"moderator_user_id"`
	ModeratorUserLogin   *string `json:"moderator_user_login"`
	ModeratorUserName    *string `json:"moderator_user_name"`
	UserID               string  `json:"user_id"`
	UserLogin            string  `json:"user_login"`
	UserName             string  `json:"user_name"`
	ResolutionText       string  `json:"resolution_text"`
	Status               string  `json:"status"`
}
