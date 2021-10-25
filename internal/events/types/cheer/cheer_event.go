// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package cheer

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

var triggerSupported = []string{"cheer"}

var triggerMapping = map[string]map[string]string{
	models.TransportEventSub: {
		"cheer": "channel.cheer",
	},
}

type Event struct{}

func (e Event) GenerateEvent(params events.MockEventParameters) (events.MockEventResponse, error) {
	var event []byte
	var err error

	if params.IsAnonymous == true {
		params.FromUserID = ""
		params.FromUserName = ""
	}

	if params.Cost <= 0 {
		params.Cost = 100
	}

	switch params.Transport {
	case models.TransportEventSub:
		body := *&models.EventsubResponse{
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
				CreatedAt: util.GetTimestamp().Format(time.RFC3339Nano),
			},
			Event: models.CheerEventSubEvent{
				UserID:               params.FromUserID,
				UserLogin:            params.FromUserName,
				UserName:             params.FromUserName,
				BroadcasterUserID:    params.ToUserID,
				BroadcasterUserLogin: params.ToUserName,
				BroadcasterUserName:  params.ToUserName,
				IsAnonymous:          params.IsAnonymous,
				Message:              "This is a test event.",
				Bits:                 params.Cost,
			},
		}

		event, err = json.Marshal(body)
		if err != nil {
			return events.MockEventResponse{}, err
		}

	case models.TransportWebSub:
		return events.MockEventResponse{}, errors.New("Websub is unsupported for cheer events")
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
func (e Event) GetEventbusAlias(t string) string {
	// check for aliases
	for trigger, topic := range triggerMapping[models.TransportEventSub] {
		if topic == t {
			return trigger
		}
	}
	return ""
}
