// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package event_name

import (
	"encoding/json"

	"github.com/twitchdev/twitch-cli/internal/events"
	"github.com/twitchdev/twitch-cli/internal/models"
)

var transportsSupported = map[string]bool{
	models.TransportWebSub:   true,
	models.TransportEventSub: true,
}

var triggerSupported = []string{"trigger_keyword"}

var triggerMapping = map[string]map[string]string{
	models.TransportWebSub: {
		"trigger_keyword": "topic_name_ws",
	},
	models.TransportEventSub: {
		"trigger_keyword": "topic_name_es",
	},
}

type Event struct{}

func (e Event) GenerateEvent(params events.MockEventParameters) (events.MockEventRespose, err) {
	var event []byte
	var err error

	switch params.Transport {
	case models.TransportEventSub:
		body := &models.EventsubResponse{
			// make the eventsub response (if supported)
		}
		event, err = json.Marshal(body)
		if err != nil {
			return events.MockEventRespose{}, err
		}
	case models.TransportWebSub:
		body := models.FollowWebSubResponse{} // replace with actual model in internal/models
		event, err = json.Marshal(body)
		if err != nil {
			return events.MockEventRespose{}, err
		}
	default:
		return events.MockEventRespose{}, nil
	}

	return events.MockEventRespose{
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
