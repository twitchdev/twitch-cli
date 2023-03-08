// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package websockets_cmd

import (
	"encoding/json"

	"github.com/twitchdev/twitch-cli/internal/events"
	"github.com/twitchdev/twitch-cli/internal/models"
	"github.com/twitchdev/twitch-cli/internal/models/models_mock_websocket"
)

var transportsSupported = map[string]bool{
	models.TransportEventSub:  false,
	models.TransportWebSocket: true,
}

var triggerSupported = []string{"mock.websocket.reconnect"}

var triggerMapping = map[string]map[string]string{
	models.TransportWebSocket: {
		"mock.websocket.reconnect": "mock.websocket.reconnect",
	},
}

type ReconnectEvent struct{}

func (e ReconnectEvent) GenerateEvent(params events.MockEventParameters) (events.MockEventResponse, error) {
	var event []byte
	var err error

	switch params.Transport {
	case models.TransportWebSocket: // WebSocket only
		body := *&models.EventsubResponse{
			Subscription: models.EventsubSubscription{
				Type:      triggerMapping[params.Transport][params.Trigger],
				Version:   e.SubscriptionVersion(),
				Transport: models.EventsubTransport{},
				CreatedAt: params.Timestamp,
			},
			Event: models_mock_websocket.ReconnectWebSocketEvent{},
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

func (e ReconnectEvent) ValidTransport(t string) bool {
	return transportsSupported[t]
}

func (e ReconnectEvent) ValidTrigger(t string) bool {
	for _, ts := range triggerSupported {
		if ts == t {
			return true
		}
	}
	return false
}
func (e ReconnectEvent) GetTopic(transport string, trigger string) string {
	return triggerMapping[transport][trigger]
}
func (e ReconnectEvent) GetAllTopicsByTransport(transport string) []string {
	allTopics := []string{}
	for _, topic := range triggerMapping[transport] {
		allTopics = append(allTopics, topic)
	}
	return allTopics
}
func (e ReconnectEvent) GetEventSubAlias(t string) string {
	// check for aliases
	for trigger, topic := range triggerMapping[models.TransportEventSub] {
		if topic == t {
			return trigger
		}
	}
	return ""
}

func (e ReconnectEvent) SubscriptionVersion() string {
	return "1"
}
