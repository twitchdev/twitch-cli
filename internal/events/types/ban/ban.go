// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package ban

import (
	"encoding/json"
	"time"

	"github.com/twitchdev/twitch-cli/internal/events"
	"github.com/twitchdev/twitch-cli/internal/models"
	"github.com/twitchdev/twitch-cli/internal/util"
)

var transportsSupported = map[string]bool{
	models.TransportEventSub: true,
}

var triggerSupported = []string{"ban", "unban"}

var triggerMapping = map[string]map[string]string{
	models.TransportEventSub: {
		"ban":   "channel.ban",
		"unban": "channel.unban",
	},
}

type Event struct{}

func (e Event) GenerateEvent(params events.MockEventParameters) (events.MockEventResponse, error) {
	var event []byte
	var err error

	switch params.Transport {
	case models.TransportEventSub:
		e := models.BanEventSubEvent{
			UserID:               params.FromUserID,
			UserLogin:            params.FromUserName,
			UserName:             params.FromUserName,
			BroadcasterUserID:    params.ToUserID,
			BroadcasterUserLogin: params.ToUserName,
			BroadcasterUserName:  params.ToUserName,
			ModeratorUserId:      util.RandomUserID(),
			ModeratorUserLogin:   "CLIModerator",
			ModeratorUserName:    "CLIModerator",
		}

		if params.Trigger == "ban" {
			reason := "This is a test event"
			endsAt := util.GetTimestamp().Format(time.RFC3339Nano)
			e.Reason = &reason
			e.EndsAt = &endsAt
			e.IsPermanent = &params.IsPermanent
		}

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
			Event: e,
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
