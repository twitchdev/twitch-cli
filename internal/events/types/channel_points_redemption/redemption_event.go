// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package channel_points_redemption

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/twitchdev/twitch-cli/internal/events"
	"github.com/twitchdev/twitch-cli/internal/models"
	"github.com/twitchdev/twitch-cli/internal/util"
)

var transportsSupported = map[string]bool{
	models.TransportWebSub:   false,
	models.TransportEventSub: true,
}

var triggerSupported = []string{"add-redemption", "update-redemption"}

var triggerMapping = map[string]map[string]string{
	models.TransportEventSub: {
		"add-redemption":    "channel.channel_points_custom_reward_redemption.add",
		"update-redemption": "channel.channel_points_custom_reward_redemption.update",
	},
}

type Event struct{}

func (e Event) GenerateEvent(params events.MockEventParameters) (events.MockEventResponse, error) {
	tNow := util.GetTimestamp().Format(time.RFC3339Nano)
	var event []byte
	var err error

	if params.Status == "" {
		params.Status = "unfulfilled"
	}

	if params.ItemID == "" {
		params.ItemID = util.RandomGUID()
	}

	if params.Cost <= 0 {
		params.Cost = 150
	}

	if params.ItemName == "" {
		params.ItemName = "Test Reward from CLI"
	}

	switch params.Transport {
	case models.TransportEventSub:
		body := *&models.RedemptionEventSubResponse{
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
				CreatedAt: tNow,
			},
			Event: models.RedemptionEventSubEvent{
				ID:                   params.ID,
				BroadcasterUserID:    params.ToUserID,
				BroadcasterUserLogin: params.ToUserName,
				BroadcasterUserName:  params.ToUserName,
				UserID:               params.FromUserID,
				UserLogin:            params.FromUserName,
				UserName:             params.FromUserName,
				UserInput:            "Test Input From CLI",
				Status:               params.Status,
				Reward: models.RedemptionReward{
					ID:     params.ItemID,
					Title:  params.ItemName,
					Cost:   params.Cost,
					Prompt: "Redeem Your Test Reward from CLI",
				},
				RedeemedAt: tNow,
			},
		}

		event, err = json.Marshal(body)
		if err != nil {
			return events.MockEventResponse{}, err
		}

	case models.TransportWebSub:
		return events.MockEventResponse{}, errors.New("Websub is unsupported for channel points events")
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
func (e Event) GetEventSubAlias(t string) string {
	// check for aliases
	for trigger, topic := range triggerMapping[models.TransportEventSub] {
		if topic == t {
			return trigger
		}
	}
	return ""
}
