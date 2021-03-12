// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package ban

import (
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

var triggerSupported = []string{"ban"}

var triggerMapping = map[string]map[string]string{
	models.TransportWebSub: {
		"ban":   "moderation.user.ban",
		"unban": "moderation.user.ban",
	},
	models.TransportEventSub: {
		"ban": "channel.ban",
		"unban": "channel.unban",
	},
}

type Event struct{}

func (e Event) GenerateEvent(params events.MockEventParameters) (events.MockEventResponse, error) {
	var event []byte
	var err error

	if params.IsAnonymous == true {
		params.FromUserID = ""
		params.FromUserName = ""
	}

	if params.Cost <= 0 {
		params.Cost = 100
	}

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
			Event: models.BanEventSubEvent{
				UserID:               params.FromUserID,
				UserLogin:            params.FromUserName,
				UserName:             params.FromUserName,
				BroadcasterUserID:    params.ToUserID,
				BroadcasterUserLogin: params.ToUserName,
				BroadcasterUserName:  params.ToUserName,
				ModeratorUserId:      params.ModeratorUserID,
				ModeratorUserLogin:   params.ModeratorUserName,
				ModeratorUserName:    params.ModeratorUserName,
				Reason:               "This is a test event",
				EndsAt: 			  util.GetTimestamp().Format(time.RFC3339Nano),
				IsPermanent: 		  params.IsPermanent,
			},
		}

		event, err = json.Marshal(body)
		if err != nil {
			return events.MockEventResponse{}, err
		}

	case models.TransportWebSub:
		body := *&models.BanWebSubResponse{
			Data: []models.BanWebSubResponseData{
				{
					ID:             params.ID,
					EventType:      triggerMapping[params.Transport][params.Trigger],
					EventTimestamp: util.GetTimestamp().Format(time.RFC3339),
					Version:        "v1",
					EventData: models.BanWebSubEventData{
						BroadcasterID:   params.ToUserID,
						BroadcasterName: params.ToUserName,
						UserID:          params.FromUserID,
						UserName:        params.FromUserName,
					},
				},
			}}

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
