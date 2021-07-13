// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package stream_change

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

var triggerSupported = []string{"stream-change"}

var triggerMapping = map[string]map[string]string{
	models.TransportWebSub: {
		"stream-change": "streams",
	},
	models.TransportEventSub: {
		"stream-change": "channel.update",
	},
}

type Event struct{}

func (e Event) GenerateEvent(params events.MockEventParameters) (events.MockEventResponse, error) {
	var event []byte
	var err error

	if params.Description == "" {
		params.Description = "Example title from the CLI!"
	}
	if params.ItemID == "" {
		params.ItemID = "509658"
	}
	if params.ItemName == "" {
		params.ItemName = "Just Chatting"
	}

	switch params.Transport {
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
				Cost:      0,
				CreatedAt: util.GetTimestamp().Format(time.RFC3339Nano),
			},
			Event: models.ChannelUpdateEventSubEvent{
				BroadcasterUserID:    params.ToUserID,
				BroadcasterUserLogin: params.ToUserName,
				BroadcasterUserName:  params.ToUserName,
				StreamTitle:          params.Description,
				StreamLanguage:       "en",
				StreamCategoryID:     params.ItemID,
				StreamCategoryName:   params.ItemName,
				IsMature:             false,
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
					WebsubID:             params.ID,
					BroadcasterUserID:    params.ToUserID,
					BroadcasterUserLogin: params.ToUserName,
					BroadcasterUserName:  params.ToUserName,
					StreamCategoryID:     params.ItemID,
					StreamCategoryName:   "Just Chatting",
					StreamType:           "live",
					StreamTitle:          params.Description,
					StreamViewerCount:    9848,
					StreamStartedAt:      util.GetTimestamp().Format(time.RFC3339),
					StreamLanguage:       "en",
					StreamThumbnailURL:   "https://static-cdn.jtvnw.net/previews-ttv/live_twitch_user-{width}x{height}.jpg",
					TagIDs:               make([]string, 0),
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
