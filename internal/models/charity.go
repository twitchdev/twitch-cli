// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package models

type CharityEventSubEventAmount struct {
	Value         int    `json:"value"`
	DecimalPlaces int    `json:"decimal_places"`
	Currency      string `json:"currency"`
}

type CharityEventSubEvent struct {
	CampaignID           string                     `json:"campaign_id"`
	BroadcasterUserID    string                     `json:"broadcaster_id"`
	BroadcasterUserName  string                     `json:"broadcaster_name"`
	BroadcasterUserLogin string                     `json:"broadcaster_login"`
	UserID               string                     `json:"user_id"`
	UserName             string                     `json:"user_name"`
	UserLogin            string                     `json:"user_login"`
	CharityName          string                     `json:"charity_name"`
	CharityLogo          string                     `json:"charity_logo"`
	Amount               CharityEventSubEventAmount `json:"amount"`
}

type CharityEventSubResponse struct {
	Subscription EventsubSubscription `json:"subscription"`
	Event        CharityEventSubEvent `json:"event"`
}
