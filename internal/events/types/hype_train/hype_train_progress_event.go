// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package hype_train_progress

import(
	"encoding/json"
	"time"

	"github.com/twitchdev/twitch-cli/internal/events"
	"github.com/twitchdev/twitch-cli/internal/models"
	"github.com/twitchdev/twitch-cli/internal/util"
)

var transportsSupported = map[string]bool{
	models.TransportWebSub:   true,
	models.TransportEventSub: true,
}

var triggerSupported = []string{"hype-train-begin", "hype-train-progress", "hype-train-end"}

var triggerMapping = map[string]map[string]string{
	models.TransportWebSub: {
		"hype-train-progress":  "hypetrain.progression",
	},
	models.TransportEventSub: {
		"hype-train-progress": 	"channel.hype_train.progress",
	},
}

type Event struct{}

func (e Event) GenerateEvent(params events.MockEventParameters) (events.MockEventResponse, error) {
	var event []byte
	var err error

	//Some values
	StartedAtTime := util.GetTimestamp()
	localTotal := util.RandomViewerCount()
	localLevel := util.RandomViewerCount()%4
	localGoal := util.RandomViewerCount()
	localProgress := (localTotal/localGoal)
	localRandomUser1 := util.RandomUserID()
	localRandomUser2 := util.RandomUserID()
	localRandomUser3 := util.RandomUserID()
	localLC := models.ContributionData{TotalContribution: util.RandomViewerCount(), TypeOfContribution: util.RandomType(), UserWhoMadeContribution: localRandomUser1, UserNameWhoMadeContribution: localRandomUser1, UserLoginWhoMadeContribution: localRandomUser1}
	localTC := []models.ContributionData{{TotalContribution: util.RandomViewerCount(), TypeOfContribution: util.RandomType(), UserWhoMadeContribution: localRandomUser2, UserNameWhoMadeContribution: localRandomUser2, UserLoginWhoMadeContribution: localRandomUser2},{TotalContribution: util.RandomViewerCount(), TypeOfContribution: util.RandomType(), UserWhoMadeContribution: localRandomUser3, UserNameWhoMadeContribution: localRandomUser3, UserLoginWhoMadeContribution: localRandomUser3}}
	ExpiresAtTime := util.GetTimestamp()
	CooldownAtTime := util.GetTimestamp()

	switch params.Transport {
	case models.TransportEventSub:
		body := *&models.EventsubResponse{
			Subscription: models.EventsubSubscription{
				ID:      params.ID,
				Status:  "enabled",
				Type:    triggerMapping[params.Transport][params.Trigger],
				Version: "1",
				Condition: models.EventsubCondition{
					BroadcasterUserID: params.ToUserID,
				},
				Transport: models.EventsubTransport{
					Method:   "webhook",
					Callback: "null",
				},
				CreatedAt: util.GetTimestamp().Format(time.RFC3339Nano),
			},
			Event: models.HypeTrainEventProgressSubEvent{
				BroadcasterUserID:    	params.ToUserID,
				BroadcasterUserLogin: 	params.ToUserName,
				BroadcasterUserName:  	params.ToUserName,
				Level:               	localLevel,
				Total:               	localTotal,
				Progress:            	localProgress,
				Goal:             		localGoal,
				TopContributions:    	localTC,
				LastContribution: 		localLC,
				StartedAtTimestamp:  	StartedAtTime.String(),
				ExpiresAtTimestamp:  	ExpiresAtTime.String(),
			},
		}

		event, err = json.Marshal(body)
		if err != nil {
			return events.MockEventResponse{}, err
		}
	case models.TransportWebSub:
		body := *&models.HypeTrainWebSubResponse{
			Data: []models.HypeTrainWebSubEvent{
				{
					ID:             params.ID,
					EventType:      triggerMapping[params.Transport][params.Trigger],
					EventTimestamp: util.GetTimestamp().Format(time.RFC3339),
					Version:        "1.0",
					EventData: models.HypeTrainWebsubEventData{
						BroadcasterID:  			 params.ToUserID,
						CooldownEndTimestamp: 		 CooldownAtTime.String(),
						ExpiresAtTimestamp:          ExpiresAtTime.String(),
						Goal:       				 localGoal,
						Id:   						 util.RandomGUID(),
						LastContribution:  			 localLC,
						Level:          			 localLevel,
						StartedAtTimestamp:        	 StartedAtTime.String(),
						TopContributions: 			 localTC,
						Total:          			 localTotal,
					},
				},
			},
		}
		event, err = json.Marshal(body)
		if err != nil {
			return events.MockEventResponse{}, err
		}

	default:
		return events.MockEventResponse{}, nil
	}

	return events.MockEventResponse{
		ID:       params.ID,
		JSON:     event,
		FromUser: params.FromUserID,
		ToUser:   params.ToUserID,
	}, nil
}

func (e Event) ValidTransport(t string) bool {
	return transportsSupported[t]
}

func (e Event) ValidTrigger(t string) bool {
	for _, ts := range triggerSupported {
		if ts == t {
			return true
		}
	}
	return false
}
func (e Event) GetTopic(transport string, trigger string) string {
	return triggerMapping[transport][trigger]
}