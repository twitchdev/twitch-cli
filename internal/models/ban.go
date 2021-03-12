// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package models

type BanEventSubEvent struct {
	UserID      			 string `json:"user_id"`
	UserLogin     			 string `json:"user_login"`
	UserName 				 string `json:"user_name"`
	BroadcasterUserID    	 string `json:"broadcaster_user_id"`
	BroadcasterUserLogin 	 string `json:"broadcaster_user_login"`
	BroadcasterUserName  	 string `json:"broadcaster_user_name"`
	ModeratorUserId      	 string `json:"moderator_user_id"`
	ModeratorUserLogin   	 string `json:"moderator_user_login"`
	ModeratorUserName    	 string `json:"moderator_user_name"`
	Reason               	 string `json:"reason"`
	EndsAt 			  		 string `json:"ends_at"`
	IsPermanent 		  	 bool `json:"is_permanent"`
}

type BanEventSubResponse struct {
	Subscription EventsubSubscription `json:"ban"`
	Event        BanEventSubEvent           `json:"event"`
}


type BanWebSubResponse struct {
	Data []BanWebSubResponseData `json:"data"`
}

type BanWebSubResponseData struct {
	ID             string             `json:"id"`
	EventType      string             `json:"event_type"`
	EventTimestamp string             `json:"event_timestamp"`
	Version        string             `json:"version"`
	EventData      BanWebSubEventData `json:"event_data"`
}

type BanWebSubEventData struct {
	BroadcasterID   string `json:"broadcaster_id"`
	BroadcasterName string `json:"broadcaster_name"`
	UserID      	string `json:"user_id"`
	UserName 		string `json:"user_name"`
}