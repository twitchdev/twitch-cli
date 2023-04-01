// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package subscription_message

import (
	"encoding/json"
	"strings"

	"github.com/twitchdev/twitch-cli/internal/events"
	"github.com/twitchdev/twitch-cli/internal/models"
	"github.com/twitchdev/twitch-cli/internal/util"
)

var transportsSupported = map[string]bool{
	models.TransportWebhook:   true,
	models.TransportWebSocket: true,
}

var triggerSupported = []string{"subscribe-message"}

var triggerMapping = map[string]map[string]string{
	models.TransportWebhook: {
		"subscribe-message": "channel.subscription.message",
	},
	models.TransportWebSocket: {
		"subscribe-message": "channel.subscription.message",
	},
}

type Event struct{}

func (e Event) GenerateEvent(params events.MockEventParameters) (events.MockEventResponse, error) {
	var event []byte
	var err error

	if params.Cost == 0 {
		params.Cost = util.RandomInt(120) + 1
	}

	switch params.Transport {
	case models.TransportWebhook, models.TransportWebSocket:
		body := &models.SubscribeMessageEventSubResponse{
			Subscription: models.EventsubSubscription{
				ID:      params.ID,
				Status:  params.SubscriptionStatus,
				Type:    triggerMapping[params.Transport][params.Trigger],
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
			Event: models.SubscribeMessageEventSubEvent{
				UserID:               params.FromUserID,
				UserLogin:            params.FromUserName,
				UserName:             params.FromUserName,
				BroadcasterUserID:    params.ToUserID,
				BroadcasterUserLogin: params.ToUserName,
				BroadcasterUserName:  params.ToUserName,
				Tier:                 params.Tier,
				Message: models.SubscribeMessageEventSubMessage{
					Text: "Hello from the Twitch CLI! twitchdevLeek",
					Emotes: []models.SubscribeMessageEventSubMessageEmote{
						{
							Begin: 26,
							End:   39,
							ID:    "304456816",
						},
					},
				},
				CumulativeMonths: int(params.Cost) + int(util.RandomInt(10)),
				DurationMonths:   1,
			},
		}

		if !params.IsAnonymous {
			streak := int(params.Cost)
			body.Event.StreakMonths = &streak
		}
		event, err = json.Marshal(body)
		if err != nil {
			return events.MockEventResponse{}, err
		}

		// Delete event info if Subscription.Status is not set to "enabled"
		if !strings.EqualFold(params.SubscriptionStatus, "enabled") {
			var i interface{}
			if err := json.Unmarshal([]byte(event), &i); err != nil {
				return events.MockEventResponse{}, err
			}
			if m, ok := i.(map[string]interface{}); ok {
				delete(m, "event") // Matches JSON key defined in body variable above
			}

			event, err = json.Marshal(i)
			if err != nil {
				return events.MockEventResponse{}, err
			}
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
func (e Event) GetAllTopicsByTransport(transport string) []string {
	allTopics := []string{}
	for _, topic := range triggerMapping[transport] {
		allTopics = append(allTopics, topic)
	}
	return allTopics
}
func (e Event) GetEventSubAlias(t string) string {
	// check for aliases
	for trigger, topic := range triggerMapping[models.TransportWebhook] {
		if topic == t {
			return trigger
		}
	}
	return ""
}

func (e Event) SubscriptionVersion() string {
	return "1"
}
