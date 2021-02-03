// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package events

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/twitchdev/twitch-cli/internal/models"
	"github.com/twitchdev/twitch-cli/internal/util"
)

type CheerParams struct {
	IsAnonymous bool
	Transport   string
	Type        string
	ToUser      string
	FromUser    string
	Message     string
	Bits        int64
}

func GenerateCheerBody(p CheerParams) (TriggerResponse, error) {
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

	if p.IsAnonymous == true {
		p.FromUser = ""
		fromUserName = ""
	}

	if p.Bits <= 0 {
		p.Bits = 100
	}

	switch p.Transport {
	case TransportEventSub:
		body := *&models.EventsubResponse{
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
			Event: models.CheerEventSubEvent{
				UserID:               p.FromUser,
				UserLogin:            toUserName,
				UserName:             toUserName,
				BroadcasterUserID:    p.ToUser,
				BroadcasterUserLogin: fromUserName,
				BroadcasterUserName:  fromUserName,
				IsAnonymous:          p.IsAnonymous,
				Message:              "This is a test event.",
				Bits:                 p.Bits,
			},
		}

		event, err = json.Marshal(body)
		if err != nil {
			return TriggerResponse{}, err
		}

	case TransportWebSub:
		return TriggerResponse{}, errors.New("Websub is unsupported for cheer events")
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
