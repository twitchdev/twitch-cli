// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package models

type CheerEventSubEvent struct {
	UserID               string  `json:"user_id"`
	UserLogin            string  `json:"user_login"`
	UserName             string  `json:"user_name"`
	BroadcasterUserID    string  `json:"broadcaster_user_id"`
	BroadcasterUserLogin string  `json:"broadcaster_user_login"`
	BroadcasterUserName  string  `json:"broadcaster_user_name"`
	IsAnonymous          bool    `json:"is_anonymous"`
	Message              string  `json:"message"`
	Bits                 float64 `json:"bits"`
}

type CheerEventSubResponse struct {
	Subscription EventsubSubscription `json:"subscription"`
	Event        CheerEventSubEvent   `json:"event"`
}
