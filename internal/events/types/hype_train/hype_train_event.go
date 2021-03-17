// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package hype_train

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
		"hype-train-begin":		"channel.hype_train.begin",
		"hype-train-end":		"channel.hype_train.end",
	},
}

type Event struct{}

func (e Event) GenerateEvent(params events.MockEventParameters) (events.MockEventResponse, error) {
	var event []byte
	var err error



	//Local variables which will be used for the trigger params below
	localTotal := util.RandomViewerCount()
	localGoal := util.RandomViewerCount()
	localProgress := (localTotal/localGoal)
	localRandomUser1 := util.RandomUserID()
	localRandomUser2 := util.RandomUserID()
	localRandomUser3 := util.RandomUserID()
	localLC := models.ContributionData{TotalContribution: util.RandomViewerCount(), TypeOfContribution: util.RandomType(), UserWhoMadeContribution: localRandomUser1, UserNameWhoMadeContribution: localRandomUser1, UserLoginWhoMadeContribution: localRandomUser1}
	localTC := []models.ContributionData{{TotalContribution: util.RandomViewerCount(), TypeOfContribution: util.RandomType(), UserWhoMadeContribution: localRandomUser2, UserNameWhoMadeContribution: localRandomUser2, UserLoginWhoMadeContribution: localRandomUser2},{TotalContribution: util.RandomViewerCount(), TypeOfContribution: util.RandomType(), UserWhoMadeContribution: localRandomUser3, UserNameWhoMadeContribution: localRandomUser3, UserLoginWhoMadeContribution: localRandomUser3}}


	if params.Trigger == "hype-train-begin" {
		switch params.Transport {
			case models.TransportEventSub:
				body := *&models.EventsubResponse{
					Subscription: models.EventsubSubscription{
						ID:      params.ID,
						Status:  "enabled",
						Type:    triggerMapping[params.Transport][params.Trigger],
						Version: "1.0",
						Condition: models.EventsubCondition{
							BroadcasterUserID: params.ToUserID,
						},
						Transport: models.EventsubTransport{
							Method:   "webhook",
							Callback: "null",
						},
						CreatedAt: util.GetTimestamp().Format(time.RFC3339Nano),
					},
					Event: models.HypeTrainEventSubEvent{
						BroadcasterUserID:    	params.ToUserID,
						BroadcasterUserLogin: 	params.ToUserName,
						BroadcasterUserName:  	params.ToUserName,
						Total:               	localTotal,
						Progress:            	localProgress,
						Goal:             		localGoal,
						TopContributions:    	localTC,
						LastContribution: 		localLC,
						StartedAtTimestamp:  	util.GetTimestamp().Format(time.RFC3339Nano),
						ExpiresAtTimestamp:  	util.GetTimestamp().Format(time.RFC3339Nano),
				},
			}

			event, err = json.Marshal(body)
			if err != nil {
				return events.MockEventResponse{}, err
			}

			default:
			return events.MockEventResponse{}, nil
		}
	}

	if params.Trigger == "hype-train-end" {
		switch params.Transport {
		case models.TransportEventSub:
			body := *&models.EventsubResponse{
				Subscription: models.EventsubSubscription{
					ID:      params.ID,
					Status:  "enabled",
					Type:    triggerMapping[params.Transport][params.Trigger],
					Version: "1.0",
					Condition: models.EventsubCondition{
						BroadcasterUserID: params.ToUserID,
					},
					Transport: models.EventsubTransport{
						Method:   "webhook",
						Callback: "null",
					},
					CreatedAt: util.GetTimestamp().Format(time.RFC3339Nano),
				},
				Event: models.HypeTrainEventSubEvent{
					BroadcasterUserID:    	 params.ToUserID,
					BroadcasterUserLogin: 	 params.ToUserName,
					BroadcasterUserName:  	 params.ToUserName,
					Level:               	 util.RandomViewerCount()%4,
					Total:             		 localTotal,
					TopContributions:    	 localTC,
					StartedAtTimestamp:  	 util.GetTimestamp().Format(time.RFC3339Nano),
					EndedAtTimestamp: 	 	 util.GetTimestamp().Format(time.RFC3339Nano),
					CooldownEndsAtTimestamp: util.GetTimestamp().Format(time.RFC3339Nano),
				},
			}

			event, err = json.Marshal(body)
			if err != nil {
				return events.MockEventResponse{}, err
			}
		default:
			return events.MockEventResponse{}, nil
		}
	}

	if params.Trigger == "hype-train-progress" {
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
					Event: models.HypeTrainEventSubEvent{
						BroadcasterUserID:    	params.ToUserID,
						BroadcasterUserLogin: 	params.ToUserName,
						BroadcasterUserName:  	params.ToUserName,
						Level:               	util.RandomViewerCount()%4,
						Total:               	localTotal,
						Progress:            	localProgress,
						Goal:             		localGoal,
						TopContributions:    	localTC,
						LastContribution: 		localLC,
						StartedAtTimestamp:  	util.GetTimestamp().Format(time.RFC3339Nano),
						ExpiresAtTimestamp:  	util.GetTimestamp().Format(time.RFC3339Nano),
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
							CooldownEndTimestamp: 		 util.GetTimestamp().Format(time.RFC3339),
							ExpiresAtTimestamp:          util.GetTimestamp().Format(time.RFC3339),
							Goal:       				 localGoal,
							Id:   						 util.RandomGUID(),
							LastContribution:  			 localLC,
							Level:          			 util.RandomViewerCount()%4,
							StartedAtTimestamp:        	 util.GetTimestamp().Format(time.RFC3339),
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