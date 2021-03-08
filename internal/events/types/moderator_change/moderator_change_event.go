// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package moderator_change

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

var triggerSupported = []string{"add-moderator", "remove-moderator"}

var triggerMapping = map[string]map[string]string{
	models.TransportWebSub: {
		"add-moderator":    "moderation.moderator.add",
		"remove-moderator": "moderation.moderator.remove",
	},
}

type Event struct{}

func (e Event) GenerateEvent(params events.MockEventParameters) (events.MockEventResponse, error) {
	var event []byte
	var err error

	switch params.Transport {
	case models.TransportEventSub:
		return events.MockEventResponse{}, errors.New("Moderator change events are currently in beta on EventSub")
	case models.TransportWebSub:
		body := *&models.ModeratorChangeWebSubResponse{
			Data: []models.ModeratorChangeWebSubEvent{
				{
					ID:             params.ID,
					EventType:      triggerMapping[params.Transport][params.Trigger],
					EventTimestamp: util.GetTimestamp().Format(time.RFC3339),
					Version:        "v1",
					EventData: models.ModeratorChangeEventData{
						BroadcasterID:   params.ToUserID,
						BroadcasterName: params.ToUserName,
						UserID:          params.FromUserID,
						UserName:        params.FromUserName,
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
