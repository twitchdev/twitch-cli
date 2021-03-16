// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package hype_train_end

import(
	"encoding/json"
	"time"

	"github.com/twitchdev/twitch-cli/internal/events"
	"github.com/twitchdev/twitch-cli/internal/models"
	"github.com/twitchdev/twitch-cli/internal/util"
)

var transportsSupported = map[string]bool{
	models.TransportWebSub:   false,
	models.TransportEventSub: true,
}

var triggerSupported = []string{"hype-train-end"}

var triggerMapping = map[string]map[string]string{
	models.TransportEventSub: {
		"hype-train-end":        	"channel.hype_train.end",
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
	localRandomUser2 := util.RandomUserID()
	localRandomUser3 := util.RandomUserID()
	localTC := []models.ContributionData{{TotalContribution: util.RandomViewerCount(), TypeOfContribution: util.RandomType(), UserWhoMadeContribution: localRandomUser2, UserNameWhoMadeContribution: localRandomUser2, UserLoginWhoMadeContribution: localRandomUser2},{TotalContribution: util.RandomViewerCount(), TypeOfContribution: util.RandomType(), UserWhoMadeContribution: localRandomUser3, UserNameWhoMadeContribution: localRandomUser3, UserLoginWhoMadeContribution: localRandomUser3}}
	EndedAtTime := util.GetTimestamp()
	CooldownAtTime := util.GetTimestamp()

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
			Event: models.HypeTrainEventEndSubEvent{
				BroadcasterUserID:    	 params.ToUserID,
				BroadcasterUserLogin: 	 params.ToUserName,
				BroadcasterUserName:  	 params.ToUserName,
				Level:               	 localLevel,
				Total:             		 localTotal,
				TopContributions:    	 localTC,
				StartedAtTimestamp:  	 StartedAtTime.String(),
				EndedAtTimestamp: 	 	 EndedAtTime.String(),
				CooldownEndsAtTimestamp: CooldownAtTime.String(),
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