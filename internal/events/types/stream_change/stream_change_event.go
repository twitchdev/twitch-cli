// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package stream_change

import (
	"encoding/json"

	"github.com/twitchdev/twitch-cli/internal/events"
	"github.com/twitchdev/twitch-cli/internal/models"
	"github.com/twitchdev/twitch-cli/internal/util"
	"time"
)

var transportsSupported = map[string]bool{
	models.TransportWebSub:   true,
	models.TransportEventSub: true,
}

var triggerSupported = []string{"stream_change"}

var triggerMapping = map[string]map[string]string{
	models.TransportWebSub: {
		"stream_change": "streams",
	},
	models.TransportEventSub: {
		"stream_change": "channel.update",
	},
}

type Event struct{}

func (e Event) GenerateEvent(params events.MockEventParameters) (events.MockEventResponse, error) {
	var event []byte
	var err error

	if params.StreamTitle == "" {
		params.StreamTitle = "Default Title!"
	}

	switch params.Transport{
	case models.TransportEventSub:
		body := &models.EventsubResponse{
			// make the eventsub response (if supported)
			Subscription: models.EventsubSubscription{
				ID:      params.ID,
				Status:  "enabled",
				Type:    "channel.update",
				Version: "1",
				Condition: models.EventsubCondition{
					BroadcasterUserID: params.ToUserID,
				},
				Transport: models.EventsubTransport{
					Method:   "webhook",
					Callback: "null",
				},
				CreatedAt: util.GetTimestamp().Format(time.RFC3339Nano),
			},
			Event: models.ChannelUpdateEventSubEvent{
				BroadcasterUserID:               params.ToUserID,
				BroadcasterUserLogin:            params.ToUserID,
				BroadcasterUserName:             params.ToUserName,
				StreamTitle:    params.StreamTitle,
				StreamLanguage: "en",
				StreamCategoryID:  "509658",
				StreamCategoryName:  "Just Chatting",
				IsMature:  "true",
			},
		}
		event, err = json.Marshal(body)
		if err != nil {
			return events.MockEventResponse{}, err
		}
	case models.TransportWebSub:
		body := models.StreamChangeWebSubResponse{
			Data: []models.StreamChangeWebSubResponseData{
				{
					WebsubID:     params.FromUserID,
					BroadcasterUserID:   params.ToUserID,
					BroadcasterUserLogin:       params.ToUserID,
					BroadcasterUserName:     params.ToUserName,
					StreamCategoryID: "509658",
					StreamCategoryName:   "Just Chatting",
					StreamType:       "live",
					StreamTitle:     params.StreamTitle,
					StreamViewerCount: 9848,
					StreamStartedAt:   util.GetTimestamp().Format(time.RFC3339),
					StreamLanguage:       "en",
					StreamThumbnailURL:     "https://static-cdn.jtvnw.net/previews-ttv/live_user_lirik-{width}x{height}.jpg",
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
