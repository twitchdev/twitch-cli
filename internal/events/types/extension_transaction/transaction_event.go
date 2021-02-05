// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package extension_transaction

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/twitchdev/twitch-cli/internal/events"
	"github.com/twitchdev/twitch-cli/internal/models"
	"github.com/twitchdev/twitch-cli/internal/util"
)

var transportsSupported = map[string]bool{
	models.TransportWebSub:   true,
	models.TransportEventSub: false,
}

var triggerSupported = []string{"transaction"}

var triggerMapping = map[string]map[string]string{
	models.TransportWebSub: {
		"transaction": "transaction",
	},
}

type Event struct{}

func (e Event) GenerateEvent(params events.MockEventParameters) (events.MockEventResponse, error) {
	var event []byte
	var err error

	if params.Cost == 0 {
		params.Cost = 100
	}
	switch params.Transport {
	case models.TransportEventSub:
		return events.MockEventResponse{}, errors.New("Extension transactions are unsupported on Eventsub")
	case models.TransportWebSub:
		body := *&models.TransactionWebSubResponse{
			Data: []models.TransactionWebsubEvent{
				{
					ID:              params.ID,
					Timestamp:       util.GetTimestamp().Format(time.RFC3339),
					BroadcasterID:   params.ToUserID,
					BroadcasterName: "testBroadcaster",
					UserID:          params.FromUserID,
					UserName:        "testUser",
					ProductType:     "BITS_IN_EXTENSION",
					Product: models.TransactionProduct{
						Sku:           "testItemSku",
						DisplayName:   "Test Trigger Item from CLI",
						Broadcast:     false,
						InDevelopment: true,
						Domain:        "",
						Expiration:    "",
						Cost: models.TransactionCost{
							Amount: params.Cost,
							Type:   "bits",
						},
					},
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
