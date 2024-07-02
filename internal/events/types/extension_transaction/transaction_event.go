// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package extension_transaction

import (
	"encoding/json"
	"strings"

	"github.com/twitchdev/twitch-cli/internal/events"
	"github.com/twitchdev/twitch-cli/internal/models"
	"github.com/twitchdev/twitch-cli/internal/util"
)

var transportsSupported = map[string]bool{
	models.TransportWebhook: true,
}

var triggerSupported = []string{"transaction"}

var triggerMapping = map[string]map[string]string{
	models.TransportWebhook: {
		"transaction": "extension.bits_transaction.create",
	},
}

type Event struct{}

func (e Event) GenerateEvent(params events.MockEventParameters) (events.MockEventResponse, error) {
	var event []byte
	var err error

	if params.Cost == 0 {
		params.Cost = 100
	}

	if params.ItemID == "" {
		params.ItemID = "testItemSku"
	}

	if params.ItemName == "" {
		params.ItemName = "Test Trigger Item from CLI"
	}

	switch params.Transport {
	case models.TransportWebhook:
		body := &models.TransactionEventSubResponse{
			Subscription: models.EventsubSubscription{
				ID:      params.SubscriptionID,
				Status:  params.SubscriptionStatus,
				Type:    triggerMapping[params.Transport][params.Trigger],
				Version: e.SubscriptionVersion(),
				Condition: models.EventsubCondition{
					ExtensionClientID: params.ClientID,
				},
				Transport: models.EventsubTransport{
					Method:   "webhook",
					Callback: "null",
				},
				Cost:      1,
				CreatedAt: params.Timestamp,
			},
			Event: models.TransactionEventSubEvent{
				ID:                   util.RandomGUID(),
				ExtensionClientID:    params.ClientID,
				BroadcasterUserID:    params.ToUserID,
				BroadcasterUserLogin: "testBroadcaster",
				BroadcasterUserName:  "testBroadcaster",
				UserName:             "testUser",
				UserLogin:            "testUser",
				UserID:               params.FromUserID,
				Product: models.TransactionEventSubProduct{
					Name:          params.ItemName,
					Sku:           params.ItemID,
					Bits:          params.Cost,
					InDevelopment: true,
				},
			},
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
		ID:       params.EventMessageID,
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
