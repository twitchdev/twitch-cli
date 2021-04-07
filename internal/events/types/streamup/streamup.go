// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package streamup

import (
	"encoding/json"
	"time"

	"github.com/twitchdev/twitch-cli/internal/events"
	"github.com/twitchdev/twitch-cli/internal/models"
	"github.com/twitchdev/twitch-cli/internal/util"
)

var transportsSupported = map[string]bool{
	models.TransportWebSub:   true,
	models.TransportEventSub: true,
}

var triggerSupported = []string{"streamup"}

var triggerMapping = map[string]map[string]string{
	models.TransportWebSub: {
		"streamup": "channel.update",
	},
	models.TransportEventSub: {
		"streamup": "stream.online",
	},
}

type Event struct{}

func (e Event) GenerateEvent(params events.MockEventParameters) (events.MockEventResponse, error) {
	var event []byte
	var err error

	if params.StreamTitle == "" {
		params.StreamTitle = "Example title from the CLI!"
	}

	switch params.Transport {
	case models.TransportEventSub:
		body := &models.EventsubResponse{
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
			Event: models.StreamUpEventSubEvent{
				ID:                   util.RandomUserID(),
				BroadcasterUserID:    params.ToUserID,
				BroadcasterUserLogin: params.ToUserName,
				BroadcasterUserName:  params.ToUserName,
				Type:                 "live",
				StartedAt:            util.GetTimestamp().Format(time.RFC3339Nano),
			},
		}
		event, err = json.Marshal(body)
		if err != nil {
			return events.MockEventResponse{}, err
		}
	case models.TransportWebSub:
		body := models.StreamUpWebSubResponse{
			Data: []models.StreamUpWebSubResponseData{
				{
					ID:           params.ID,
					UserID:       params.ToUserID,
					UserLogin:    params.ToUserName,
					UserName:     params.ToUserName,
					GameID:       "509658",
					Type:         "live",
					Title:        params.StreamTitle,
					ViewerCount:  1337,
					StartedAt:    util.GetTimestamp().Format(time.RFC3339),
					Language:     "en",
					ThumbnailURL: "https://static-cdn.jtvnw.net/ttv-static/404_preview-440x248.jpg",
					TagIDs:       make([]string, 0),
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
