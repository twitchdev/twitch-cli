// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package shoutout

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/twitchdev/twitch-cli/internal/events"
	"github.com/twitchdev/twitch-cli/internal/models"
	"github.com/twitchdev/twitch-cli/internal/util"
)

var transportsSupported = map[string]bool{
	models.TransportEventSub:  true,
	models.TransportWebSocket: true,
}
var triggers = []string{"shoutout-create", "shoutout-received"}

var triggerMapping = map[string]map[string]string{
	models.TransportEventSub: {
		"shoutout-create":   "channel.shoutout.create",
		"shoutout-received": "channel.shoutout.receive",
	},
	models.TransportWebSocket: {
		"shoutout-create":   "channel.shoutout.create",
		"shoutout-received": "channel.shoutout.receive",
	},
}

type Event struct{}

func (e Event) GenerateEvent(params events.MockEventParameters) (events.MockEventResponse, error) {
	var event []byte
	var err error

	switch params.Transport {
	case models.TransportEventSub, models.TransportWebSocket:
		viewerCount := util.RandomInt(2000)
		startedAt := util.GetTimestamp()

		moderatorUserID := "3502151007"

		if params.Trigger == "shoutout-create" {
			body := models.ShoutoutCreateEventSubResponse{
				Subscription: models.EventsubSubscription{
					ID:      params.ID,
					Status:  params.SubscriptionStatus,
					Type:    triggerMapping[params.Transport][params.Trigger],
					Version: e.SubscriptionVersion(),
					Condition: models.EventsubCondition{
						BroadcasterUserID: params.FromUserID,
						ModeratorUserID:   moderatorUserID,
					},
					Transport: models.EventsubTransport{
						Method:   "webhook",
						Callback: "null",
					},
					Cost:      0,
					CreatedAt: params.Timestamp,
				},
				Event: models.ShoutoutCreateEventSubEvent{
					BroadcasterUserID:      params.FromUserID,
					BroadcasterUserName:    params.FromUserName,
					BroadcasterUserLogin:   params.FromUserName,
					ToBroadcasterUserID:    params.ToUserID,
					ToBroadcasterUserName:  params.ToUserName,
					ToBroadcasterUserLogin: params.ToUserName,
					ModeratorUserID:        moderatorUserID,
					ModeratorUserName:      "TrustedUser123",
					ModeratorUserLogin:     "trusteduser123",
					ViewerCount:            int(viewerCount),
					StartedAt:              startedAt.Format(time.RFC3339Nano),
					CooldownEndsAt:         startedAt.Add(2 * time.Minute).Format(time.RFC3339Nano),
					TargetCooldownEndsAt:   startedAt.Add(1 * time.Hour).Format(time.RFC3339Nano),
				},
			}

			event, err = json.Marshal(body)
			if err != nil {
				return events.MockEventResponse{}, err
			}
		} else if params.Trigger == "shoutout-received" {
			body := models.ShoutoutReceivedEventSubResponse{
				Subscription: models.EventsubSubscription{
					ID:      params.ID,
					Status:  params.SubscriptionStatus,
					Type:    triggerMapping[params.Transport][params.Trigger],
					Version: e.SubscriptionVersion(),
					Condition: models.EventsubCondition{
						BroadcasterUserID: params.ToUserID,
						ModeratorUserID:   moderatorUserID,
					},
					Transport: models.EventsubTransport{
						Method:   "webhook",
						Callback: "null",
					},
					Cost:      0,
					CreatedAt: params.Timestamp,
				},
				Event: models.ShoutoutReceivedEventSubEvent{
					BroadcasterUserID:        params.ToUserID,
					BroadcasterUserName:      params.ToUserName,
					BroadcasterUserLogin:     params.ToUserName,
					FromBroadcasterUserID:    params.FromUserID,
					FromBroadcasterUserName:  params.FromUserName,
					FromBroadcasterUserLogin: params.FromUserName,
					ViewerCount:              int(viewerCount),
					StartedAt:                startedAt.Format(time.RFC3339Nano),
				},
			}

			event, err = json.Marshal(body)
			if err != nil {
				return events.MockEventResponse{}, err
			}
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
		ToUser:   params.ToUserID,
		FromUser: params.FromUserID,
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
func (e Event) GetAllTopicsByTransport(transport string) []string {
	allTopics := []string{}
	for _, topic := range triggerMapping[transport] {
		allTopics = append(allTopics, topic)
	}
	return allTopics
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
