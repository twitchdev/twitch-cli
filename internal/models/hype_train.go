// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package models

type ContributionData struct{
	TotalContribution   		 	 int64 	`json:"total"`
	TypeOfContribution        		 string `json:"type"`
	UserWhoMadeContribution	    	 string `json:"user_id"`
	UserNameWhoMadeContribution	     string `json:"user_name"`
	UserLoginWhoMadeContribution     string `json:"user_login"`
}

type HypeTrainWebSubEvent struct {
	ID             string                   `json:"id"`
	EventType      string                   `json:"event_type"`
	EventTimestamp string                   `json:"event_timestamp"`
	Version        string                   `json:"version"`
	EventData      HypeTrainWebsubEventData `json:"event_data"`
}

type HypeTrainWebsubEventData struct {
	BroadcasterID   			string				 	`json:"broadcaster_id,omitempty"`
	CooldownEndTimestamp        string				 	`json:"cooldown_end_time,omitempty"`
	ExpiresAtTimestamp	        string 				 	`json:"expires_at,omitempty"`
	Goal       				    int64 				 	`json:"goal,omitempty"`
	Id       				    string 				 	`json:"id,omitempty"`
	LastContribution 			ContributionData 	 	`json:"last_contribution,omitempty"`
	Level      				    int64 				 	`json:"level,omitempty"`
	StartedAtTimestamp		    string 				 	`json:"started_at,omitempty"`
	TopContributions 			[]ContributionData  	`json:"top_contributions,omitempty"`
	Total      				    int64 				 	`json:"total,omitempty"`
}

type HypeTrainWebSubResponse struct {
	Data []HypeTrainWebSubEvent `json:"data"`
}

type HypeTrainEventSubResponse struct {
	Subscription EventsubSubscription         	`json:"subscription"`
	Event        HypeTrainEventSubEvent 		`json:"event"`
}

type HypeTrainEventSubEvent struct {
	BroadcasterUserID    		string 					`json:"broadcaster_user_id,omitempty"`
	BroadcasterUserLogin 		string 					`json:"broadcaster_user_login,omitempty"`
	BroadcasterUserName  		string 					`json:"broadcaster_user_name,omitempty"`
	Level   			 		int64 					`json:"level,omitempty"`
	Total   			 		int64 					`json:"total,omitempty"`
	Progress		     		int64 					`json:"progress,omitempty"`
	Goal       			 		int64					`json:"goal,omitempty"`
	TopContributions 	 		[]ContributionData  	`json:"top_contributions,omitempty"`
	LastContribution 			ContributionData 		`json:"last_contribution,omitempty"`
	StartedAtTimestamp		    string 					`json:"started_at,omitempty"`
	ExpiresAtTimestamp	        string 					`json:"expires_at,omitempty"`
	EndedAtTimestamp	        string 					`json:"ended_at,omitempty"`
	CooldownEndsAtTimestamp	    string 					`json:"cooldown_ends_at,omitempty"`
}