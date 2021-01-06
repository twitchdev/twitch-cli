// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package trigger

import (
	"encoding/json"
	"time"

	"github.com/twitchdev/twitch-cli/internal/models"
	"github.com/twitchdev/twitch-cli/internal/util"
)

type SubscribeParams struct {
	IsGift          bool
	IsAnonymousGift bool
	Transport       string
	Type            string
	ToUser          string
	FromUser        string
	GiftUser        string
}

func GenerateSubBody(params SubscribeParams) (TriggerResponse, error) {
	uuid := util.RandomGUID()
	var event []byte
	var err error

	fromUserName := "testFromuser"

	toUserName := "testBroadcaster"

	gifterUserName := ""

	if params.ToUser == "" {
		params.ToUser = util.RandomUserID()
	}

	if params.FromUser == "" {
		params.FromUser = util.RandomUserID()
	}

	if params.IsGift == true && params.GiftUser == "" {
		params.GiftUser = util.RandomUserID()
		gifterUserName = "testGifter"
	}

	if params.IsAnonymousGift == true {
		params.GiftUser = "274598607"
		gifterUserName = "ananonymousgifter"
	}

	switch params.Transport {
	case "eventsub":
		body := *&models.EventsubResponse{
			Subscription: models.EventsubSubscription{
				ID:      uuid,
				Type:    params.Type,
				Version: "test",
				Condition: models.EventsubCondition{
					BroadcasterUserID: params.ToUser,
				},
				Transport: models.EventsubTransport{
					Method:   "webhook",
					Callback: "null",
				},
				CreatedAt: util.GetTimestamp().Format(time.RFC3339),
			},
			Event: models.SubEventSubEvent{
				UserID:               params.FromUser,
				UserLogin:            fromUserName,
				UserName:             fromUserName,
				BroadcasterUserID:    params.ToUser,
				BroadcasterUserLogin: toUserName,
				BroadcasterUserName:  toUserName,
			},
		}

		event, err = json.Marshal(body)
		if err != nil {
			return TriggerResponse{}, err
		}
	case "websub":
		body := *&models.SubWebSubResponse{
			Data: []models.SubWebSubResponseData{
				{
					ID:             "test",
					EventType:      params.Type,
					EventTimestamp: time.Now().Format(time.RFC3339),
					Version:        "1.0",
					EventData: models.SubWebSubEventData{
						BroadcasterID:   params.ToUser,
						BroadcasterName: toUserName,
						UserID:          params.FromUser,
						UserName:        fromUserName,
						Tier:            "1000",
						PlanName:        "Tier 1 Test Sub",
						IsGift:          params.IsGift,
						GifterID:        params.GiftUser,
						GifterName:      gifterUserName,
					},
				},
			}}

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
		FromUser: params.FromUser,
		ToUser:   params.ToUser,
	}, nil
}
