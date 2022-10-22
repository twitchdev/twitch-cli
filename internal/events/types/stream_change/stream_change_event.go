// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package stream_change

import (
	"encoding/json"

	"github.com/twitchdev/twitch-cli/internal/events"
	"github.com/twitchdev/twitch-cli/internal/models"
)

var transportsSupported = map[string]bool{
	models.TransportEventSub: true,
}

var triggerSupported = []string{"stream-change"}

var triggerMapping = map[string]map[string]string{
	models.TransportEventSub: {
		"stream-change": "channel.update",
	},
}

type Event struct{}

func (e Event) GenerateEvent(params events.MockEventParameters) (events.MockEventResponse, error) {
	var event []byte
	var err error

	if params.Description == "" {
		params.Description = "Example title from the CLI!"
	}
	if params.ItemID == "" && params.GameID == "" {
		params.GameID = "509658"
	} else if params.ItemID != "" && params.GameID == "" {
		params.GameID = params.ItemID
	}
	if params.ItemName == "" {
		params.ItemName = "Just Chatting"
	}

	switch params.Transport {
	case models.TransportEventSub:
		body := &models.EventsubResponse{
			// make the eventsub response (if supported)
			Subscription: models.EventsubSubscription{
				ID:      params.ID,
				Status:  "enabled",
				Type:    "channel.update",
				Version: e.SubscriptionVersion(),
				Condition: models.EventsubCondition{
					BroadcasterUserID: params.ToUserID,
				},
				Transport: models.EventsubTransport{
					Method:   "webhook",
					Callback: "null",
				},
				Cost:      0,
				CreatedAt: params.Timestamp,
			},
			Event: models.ChannelUpdateEventSubEvent{
				BroadcasterUserID:    params.ToUserID,
				BroadcasterUserLogin: params.ToUserName,
				BroadcasterUserName:  params.ToUserName,
				StreamTitle:          params.Description,
				StreamLanguage:       "en",
				StreamCategoryID:     params.GameID,
				StreamCategoryName:   params.ItemName,
				IsMature:             false,
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
func (e Event) GetEventSubAlias(t string) string {
	// check for aliases
	for trigger, topic := range triggerMapping[models.TransportEventSub] {
		if topic == t {
			return trigger
		}
	}
	return ""
}

func (e Event) SubscriptionVersion() string {
	return "1"
}
