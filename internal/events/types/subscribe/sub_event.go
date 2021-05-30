// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package subscribe

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

var triggerSupported = []string{"subscribe", "gift", "unsubscribe"}

var triggerMapping = map[string]map[string]string{
	models.TransportWebSub: {
		"subscribe":   "subscriptions.subscribe",
		"unsubscribe": "subscriptions.unsubscribe",
		"gift":        "subscriptions.subscribe",
	},
	models.TransportEventSub: {
		"subscribe":   "channel.subscribe",
		"unsubscribe": "channel.unsubscribe",
		"gift":        "channel.subscribe",
	},
}

type Event struct{}

func (e Event) GenerateEvent(params events.MockEventParameters) (events.MockEventResponse, error) {
	var (
		event        []byte
        err          error
        giftUserID   string
        giftUserName string
    )

	if params.Trigger == "gift" {
		params.IsGift = true
	}

	if params.IsGift == true {
		giftUserID = util.RandomUserID()
		giftUserName = "testGifter"
	}

	if params.IsAnonymous == true {
		giftUserID = "274598607"
		giftUserName = "ananonymousgifter"
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
				Cost:      0,
				CreatedAt: util.GetTimestamp().Format(time.RFC3339Nano),
			},
			Event: models.SubEventSubEvent{
				UserID:               params.FromUserID,
				UserLogin:            params.FromUserName,
				UserName:             params.FromUserName,
				BroadcasterUserID:    params.ToUserID,
				BroadcasterUserLogin: params.ToUserName,
				BroadcasterUserName:  params.ToUserName,
				Tier:                 "1000",
				IsGift:               params.IsGift,
			},
		}

		event, err = json.Marshal(body)
		if err != nil {
			return events.MockEventResponse{}, err
		}
	case models.TransportWebSub:
		body := *&models.SubWebSubResponse{
			Data: []models.SubWebSubResponseData{
				{
					ID:             params.ID,
					EventType:      triggerMapping[params.Transport][params.Trigger],
					EventTimestamp: util.GetTimestamp().Format(time.RFC3339),
					Version:        "1.0",
					EventData: models.SubWebSubEventData{
						BroadcasterID:   params.ToUserID,
						BroadcasterName: params.ToUserName,
						UserID:          params.FromUserID,
						UserName:        params.FromUserID,
						Tier:            "1000",
						PlanName:        "Tier 1 Test Sub",
						IsGift:          params.IsGift,
						GifterID:        giftUserID,
						GifterName:      giftUserName,
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
