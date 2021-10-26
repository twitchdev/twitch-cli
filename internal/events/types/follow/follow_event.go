// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package follow

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
var triggers = []string{"follow"}

var triggerMapping = map[string]map[string]string{
	models.TransportEventSub: {
		"follow": "channel.follow",
	},
	models.TransportWebSub: {
		"follow": "follow",
	},
}

type Event struct{}

func (e Event) GenerateEvent(p events.MockEventParameters) (events.MockEventResponse, error) {
	var event []byte
	var err error

	switch p.Transport {
	case models.TransportEventSub:
		body := models.EventsubResponse{
			Subscription: models.EventsubSubscription{
				ID:      p.ID,
				Status:  "enabled",
				Type:    "channel.follow",
				Version: "1",
				Condition: models.EventsubCondition{
					BroadcasterUserID: p.ToUserID,
				},
				Transport: models.EventsubTransport{
					Method:   "webhook",
					Callback: "null",
				},
				Cost:      0,
				CreatedAt: util.GetTimestamp().Format(time.RFC3339Nano),
			},
			Event: models.FollowEventSubEvent{
				UserID:               p.FromUserID,
				UserLogin:            p.FromUserName,
				UserName:             p.FromUserName,
				BroadcasterUserID:    p.ToUserID,
				BroadcasterUserLogin: p.ToUserID,
				BroadcasterUserName:  p.ToUserName,
				FollowedAt:           util.GetTimestamp().Format(time.RFC3339Nano),
			},
		}

		event, err = json.Marshal(body)
		if err != nil {
			return events.MockEventResponse{}, err
		}

	case models.TransportWebSub:
		body := models.FollowWebSubResponse{
			Data: []models.FollowWebSubResponseData{
				{
					FromID:     p.FromUserID,
					FromName:   p.FromUserName,
					ToID:       p.ToUserID,
					ToName:     p.ToUserName,
					FollowedAt: util.GetTimestamp().Format(time.RFC3339),
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
		ID:       p.ID,
		JSON:     event,
		FromUser: p.FromUserID,
		ToUser:   p.ToUserID,
	}, nil
}

func (e Event) ValidTransport(transport string) bool {
	return transportsSupported[transport]
}

func (e Event) ValidTrigger(trigger string) bool {
	for _, t := range triggers {
		if t == trigger {
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
