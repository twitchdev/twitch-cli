// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package unban_requests

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/twitchdev/twitch-cli/internal/events"
	"github.com/twitchdev/twitch-cli/internal/models"
	"github.com/twitchdev/twitch-cli/internal/util"
)

var transportsSupported = map[string]bool{
	models.TransportWebhook:   true,
	models.TransportWebSocket: true,
}
var triggers = []string{"unban-request-create", "unban-request-resolve"}

var triggerMapping = map[string]map[string]string{
	models.TransportWebhook: {
		"unban-request-create":  "channel.unban_request.create",
		"unban-request-resolve": "channel.unban_request.resolve",
	},
	models.TransportWebSocket: {
		"unban-request-create":  "channel.unban_request.create",
		"unban-request-resolve": "channel.unban_request.resolve",
	},
}

type Event struct{}

func (e Event) GenerateEvent(params events.MockEventParameters) (events.MockEventResponse, error) {
	var event []byte
	var err error

	var unbanRequestEvent interface{}

	if params.Trigger == "unban-request-create" {
		unbanRequestEvent = models.UnbanRequestCreateEventSubEvent{
			BroadcasterUserID:    params.ToUserID,
			BroadcasterUserName:  params.ToUserName,
			BroadcasterUserLogin: strings.ToLower(params.ToUserName),
			UserID:               params.FromUserID,
			UserName:             params.FromUserName,
			UserLogin:            strings.ToLower(params.FromUserName),
			Text:                 "Please unban me!",
			CreatedAt:            util.GetTimestamp().Add(-30 * time.Minute).Format(time.RFC3339Nano),
		}
	}

	if params.Trigger == "unban-request-resolve" {
		mod_user := util.RandomUserID()
		mod_user_lower := strings.ToLower(mod_user)
		mod_user_id := util.RandomUserID()

		unbanRequestEvent = models.UnbanRequestResolveEventSubEvent{
			ID:                   util.RandomGUID(),
			BroadcasterUserID:    params.ToUserID,
			BroadcasterUserName:  params.ToUserName,
			BroadcasterUserLogin: strings.ToLower(params.ToUserName),
			ModeratorUserID:      &mod_user_id,
			ModeratorUserName:    &mod_user,
			ModeratorUserLogin:   &mod_user_lower,
			UserID:               params.FromUserID,
			UserName:             params.FromUserName,
			UserLogin:            strings.ToLower(params.FromUserName),
			ResolutionText:       "We forgive you",
			Status:               "approved",
		}
	}

	switch params.Transport {
	case models.TransportWebhook, models.TransportWebSocket:
		body := models.EventsubResponse{
			Subscription: models.EventsubSubscription{
				ID:      params.SubscriptionID,
				Type:    triggerMapping[params.Transport][params.Trigger],
				Version: e.SubscriptionVersion(),
				Status:  params.SubscriptionStatus,
				Cost:    0,
				Condition: models.EventsubCondition{
					BroadcasterUserID: params.ToUserID,
					ModeratorUserID:   params.FromUserID,
				},
				Transport: models.EventsubTransport{
					Method:   "webhook",
					Callback: "null",
				},
				CreatedAt: params.Timestamp,
			},
			Event: unbanRequestEvent,
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
		ID:     params.EventMessageID,
		JSON:   event,
		ToUser: params.ToUserID,
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
