// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package events

import (
	"encoding/json"
	"time"

	"github.com/twitchdev/twitch-cli/internal/models"
	"github.com/twitchdev/twitch-cli/internal/util"
)

type FollowParams struct {
	Transport string
	FromUser  string
	ToUser    string
	Type      string
}

func GenerateFollowBody(p FollowParams) (TriggerResponse, error) {
	uuid := util.RandomGUID()
	var event []byte
	var err error

	fromUserName := "testFromuser"

	toUserName := "testBroadcaster"

	if p.ToUser == "" {
		p.ToUser = util.RandomUserID()
	}

	if p.FromUser == "" {
		p.FromUser = util.RandomUserID()
	}

	switch p.Transport {
	case TransportEventSub:
		body := models.EventsubResponse{
			Subscription: models.EventsubSubscription{
				ID:      uuid,
				Status:  "enabled",
				Type:    p.Type,
				Version: "1",
				Condition: models.EventsubCondition{
					BroadcasterUserID: p.ToUser,
				},
				Transport: models.EventsubTransport{
					Method:   "webhook",
					Callback: "null",
				},
				CreatedAt: util.GetTimestamp().Format(time.RFC3339Nano),
			},
			Event: models.FollowEventSubEvent{
				UserID:               p.FromUser,
				UserLogin:            fromUserName,
				UserName:             fromUserName,
				BroadcasterUserID:    p.ToUser,
				BroadcasterUserLogin: toUserName,
				BroadcasterUserName:  toUserName,
			},
		}

		event, err = json.Marshal(body)
		if err != nil {
			return TriggerResponse{}, err
		}

	case TransportWebSub:
		body := models.FollowWebSubResponse{
			Data: []models.FollowWebSubResponseData{
				{
					FromID:     p.FromUser,
					FromName:   fromUserName,
					ToID:       p.ToUser,
					ToName:     toUserName,
					FollowedAt: util.GetTimestamp().Format(time.RFC3339),
				},
			},
		}
		event, err = json.Marshal(body)
		if err != nil {
			return TriggerResponse{}, err
		}

	default:
		return TriggerResponse{}, nil
	}

	return TriggerResponse{
		ID:       uuid,
		JSON:     event,
		FromUser: p.FromUser,
		ToUser:   p.ToUser,
	}, nil
}
