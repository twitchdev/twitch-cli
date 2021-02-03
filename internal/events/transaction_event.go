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

type TransactionParams struct {
	IsGift    bool
	Transport string
	Type      string
	ToUser    string
	FromUser  string
	GiftUser  string
}

func GenerateTransactionBody(params TransactionParams) (TriggerResponse, error) {
	uuid := util.RandomGUID()
	var event []byte
	var err error

	if params.ToUser == "" {
		params.ToUser = util.RandomUserID()
	}

	if params.FromUser == "" {
		params.FromUser = util.RandomUserID()
	}

	switch params.Transport {
	case TransportEventSub:
		return TriggerResponse{}, errors.New("Extension transactions are unsupported on Eventsub")
	case TransportWebSub:
		body := *&models.TransactionWebSubResponse{
			Data: []models.TransactionWebsubEvent{
				{
					ID:              uuid,
					Timestamp:       util.GetTimestamp().Format(time.RFC3339Nano),
					BroadcasterID:   params.ToUser,
					BroadcasterName: "testBroadcaster",
					UserID:          params.FromUser,
					UserName:        "testUser",
					ProductType:     "BITS_IN_EXTENSION",
					Product: models.TransactionProduct{
						Sku:           "testItemSku",
						DisplayName:   "Test Trigger Item from CLI",
						Broadcast:     false,
						InDevelopment: true,
						Domain:        "",
						Expiration:    "",
						Cost: models.TransactionCost{
							Amount: 100,
							Type:   "bits",
						},
					},
				},
			},
		}
		event, err = json.Marshal(body)
		if err != nil {
			return TriggerResponse{}, err
		}

		return TriggerResponse{
			ID:       uuid,
			JSON:     event,
			FromUser: params.FromUser,
			ToUser:   params.ToUser,
		}, nil
	default:
		return TriggerResponse{}, nil
	}
}
