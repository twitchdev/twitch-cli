// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package models


type HypeTrainWebSubEvent struct {
	ID             string                   `json:"id"`
	EventType      string                   `json:"event_type"`
	EventTimestamp string                   `json:"event_timestamp"`
	Version        string                   `json:"version"`
	EventData      HypeTrainWebsubEventData `json:"event_data"`
}

type HypeTrainWebsubEventData struct {
	BroadcasterID   			string				 	`json:"broadcaster_id"`
	CooldownEndTimestamp        string				 	`json:"cooldown_end_time"`
	ExpiresAtTimestamp	        string 				 	`json:"expires_at"`
	Goal       				    int64 				 	`json:"goal"`
	Id       				    string 				 	`json:"id"`
	LastContribution 			ContributionData 	 	`json:"last_contribution"`
	Level      				    int64 				 	`json:"level"`
	StartedAtTimestamp		    string 				 	`json:"started_at"`
	TopContributions 			[]ContributionData  	`json:"top_contributions"`
	Total      				    int64 				 	`json:"total"`
}

type ContributionData struct{
	TotalContribution   		 	 int64 	`json:"total"`
	TypeOfContribution        		 string `json:"type"`
	UserWhoMadeContribution	    	 string `json:"user_id"`
	UserNameWhoMadeContribution	     string `json:"user_name"`
	UserLoginWhoMadeContribution     string `json:"user_login"`
}

type HypeTrainWebSubResponse struct {
	Data []HypeTrainWebSubEvent `json:"data"`
}

type HypeTrainEventBeginSubResponse struct {
	Subscription EventsubSubscription         	`json:"subscription"`
	Event        HypeTrainEventBeginSubEvent 	`json:"event"`
}

type HypeTrainEventProgressSubResponse struct {
	Subscription EventsubSubscription         	`json:"subscription"`
	Event        HypeTrainEventProgressSubEvent `json:"event"`
}

type HypeTrainEventEndSubResponse struct {
	Subscription EventsubSubscription         	`json:"subscription"`
	Event        HypeTrainEventEndSubEvent 		`json:"event"`
}

type HypeTrainEventBeginSubEvent struct {
	BroadcasterUserID    		string 					`json:"broadcaster_user_id"`
	BroadcasterUserLogin 		string 					`json:"broadcaster_user_login"`
	BroadcasterUserName  		string 					`json:"broadcaster_user_name"`
	Total   			 		int64 					`json:"total"`
	Progress		     		int64 					`json:"progress"`
	Goal       			 		int64					`json:"goal"`
	TopContributions 	 		[]ContributionData  	`json:"top_contributions"`
	LastContribution 			ContributionData 		`json:"last_contribution"`
	StartedAtTimestamp		    string 					`json:"started_at"`
	ExpiresAtTimestamp	        string 					`json:"expires_at"`
}

type HypeTrainEventProgressSubEvent struct {
	BroadcasterUserID    		string 					`json:"broadcaster_user_id"`
	BroadcasterUserLogin 		string 					`json:"broadcaster_user_login"`
	BroadcasterUserName  		string 					`json:"broadcaster_user_name"`
	Level   			 		int64 					`json:"level"`
	Total   			 		int64 					`json:"total"`
	Progress		     		int64 					`json:"progress"`
	Goal       			 		int64					`json:"goal"`
	TopContributions 	 		[]ContributionData  	`json:"top_contributions"`
	LastContribution 			ContributionData 		`json:"last_contribution"`
	StartedAtTimestamp		    string 					`json:"started_at"`
	ExpiresAtTimestamp	        string 					`json:"expires_at"`
}

type HypeTrainEventEndSubEvent struct {
	BroadcasterUserID    		string 					`json:"broadcaster_user_id"`
	BroadcasterUserLogin 		string 					`json:"broadcaster_user_login"`
	BroadcasterUserName  		string 					`json:"broadcaster_user_name"`
	Level   			 		int64 					`json:"level"`
	Total   			 		int64 					`json:"total"`
	TopContributions 	 		[]ContributionData  	`json:"top_contributions"`
	StartedAtTimestamp		    string 					`json:"started_at"`
	EndedAtTimestamp	        string 					`json:"ended_at"`
	CooldownEndsAtTimestamp	    string 					`json:"cooldown_ends_at"`
}