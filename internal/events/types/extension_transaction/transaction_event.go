// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package extension_transaction

import (
	"encoding/json"

	"github.com/spf13/viper"
	"github.com/twitchdev/twitch-cli/internal/events"
	"github.com/twitchdev/twitch-cli/internal/models"
	"github.com/twitchdev/twitch-cli/internal/util"
)

var transportsSupported = map[string]bool{
	models.TransportEventSub: true,
}

var triggerSupported = []string{"transaction"}

var triggerMapping = map[string]map[string]string{
	models.TransportEventSub: {
		"transaction": "extension.bits_transaction.create",
	},
}

type Event struct{}

func (e Event) GenerateEvent(params events.MockEventParameters) (events.MockEventResponse, error) {
	var event []byte
	var err error

	clientID := viper.GetString("clientId")

	// if not configured, generate a random one
	if clientID == "" {
		clientID = util.RandomClientID()
	}

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
	case models.TransportEventSub:
		body := &models.TransactionEventSubResponse{
			Subscription: models.EventsubSubscription{
				ID:      params.ID,
				Status:  "enabled",
				Type:    triggerMapping[params.Transport][params.Trigger],
				Version: e.SubscriptionVersion(),
				Condition: models.EventsubCondition{
					ExtensionClientID: clientID,
				},
				Transport: models.EventsubTransport{
					Method:   "webhook",
					Callback: "null",
				},
				Cost:      1,
				CreatedAt: params.Timestamp,
			},
			Event: models.TransactionEventSubEvent{
				ID:                   params.ID,
				ExtensionClientID:    clientID,
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
