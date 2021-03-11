// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package models

type StreamUpEventSubResponse struct {
	Subscription EventsubSubscription  `json:"subscription"`
	Event        StreamUpEventSubEvent `json:"event"`
}

type StreamUpEventSubEvent struct {
	ID                   string `json:"id"`
	BroadcasterUserID    string `json:"broadcaster_user_id"`
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
	BroadcasterUserName  string `json:"broadcaster_user_name"`
	Type                 string `json:"type"`
	StartedAt            string `json:"started_at"`
}

type StreamUpWebSubResponse struct {
	Data []StreamUpWebSubResponseData `json:"data"`
}

type StreamUpWebSubResponseData struct {
	ID    			string  	`json:"id"`
	UserID  		string  	`json:"user_id"`
	UserLogin  		string  	`json:"user_login"`
	UserName  		string  	`json:"user_name"`
	GameID  		string  	`json:"game_id"`
	Type    		string 		`json:"type"`
	Title    		string 		`json:"title"`
	ViewerCount     int64       `json:"viewer_count"`
	StartedAt       string 		`json:"started_at"`
	Language 		string 		`json:language`
	ThumbnailURL 	string 		`json:thumbnail_url`
	TagIDs			[]string  	`json:"tag_ids"`
}