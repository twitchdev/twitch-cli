// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package channel_points_reward

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

var triggerSupported = []string{"add-reward", "update-reward", "remove-reward"}

var triggerMapping = map[string]map[string]string{
	models.TransportEventSub: {
		"add-reward":    "channel.channel_points_custom_reward.add",
		"update-reward": "channel.channel_points_custom_reward.update",
		"remove-reward": "channel.channel_points_custom_reward.remove",
	},
}

type Event struct{}

func (e Event) GenerateEvent(params events.MockEventParameters) (events.MockEventResponse, error) {
	tNow := util.GetTimestamp().Format(time.RFC3339Nano)
	var event []byte
	var err error

	if params.Cost <= 0 {
		params.Cost = 150
	}

	if params.ItemName == "" {
		params.ItemName = "Test Reward from CLI"
	}

	switch params.Transport {
	case models.TransportEventSub:
		body := models.EventsubResponse{
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
			Event: models.RewardEventSubEvent{
				ID:                                params.ID,
				BroadcasterUserID:                 params.ToUserID,
				BroadcasterUserLogin:              params.ToUserName,
				BroadcasterUserName:               params.ToUserName,
				IsEnabled:                         true,
				IsPaused:                          false,
				IsInStock:                         true,
				Title:                             params.ItemName,
				Cost:                              params.Cost,
				Prompt:                            "Redeem Your Test Reward from CLI",
				IsUserInputRequired:               true,
				ShouldRedemptionsSkipRequestQueue: false,
				CooldownExpiresAt:                 tNow,
				RedemptionsRedeemedCurrentStream:  0,
				MaxPerStream: models.RewardMax{
					IsEnabled: true,
					Value:     100,
				},
				MaxPerUserPerStream: models.RewardMax{
					IsEnabled: true,
					Value:     100,
				},
				GlobalCooldown: models.RewardGlobalCooldown{
					IsEnabled: true,
					Seconds:   300,
				},
				BackgroundColor: "#c0ffee",
				Image: models.RewardImage{
					URL1x: "https://static-cdn.jtvnw.net/image-1.png",
					URL2x: "https://static-cdn.jtvnw.net/image-2.png",
					URL4x: "https://static-cdn.jtvnw.net/image-4.png",
				},
				DefaultImage: models.RewardImage{
					URL1x: "https://static-cdn.jtvnw.net/default-1.png",
					URL2x: "https://static-cdn.jtvnw.net/default-2.png",
					URL4x: "https://static-cdn.jtvnw.net/default-4.png",
				},
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
		FromUser: params.ToUserID,
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
func (e Event) GetEventbusAlias(t string) string {
	// check for aliases
	for trigger, topic := range triggerMapping[models.TransportEventSub] {
		if topic == t {
			return trigger
		}
	}
	return ""
}
