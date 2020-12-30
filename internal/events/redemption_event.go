// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package trigger

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/twitchdev/twitch-cli/internal/models"
	"github.com/twitchdev/twitch-cli/internal/util"
)

type RedemptionParams struct {
	Transport string
	Type      string
	ToUser    string
	FromUser  string
	Title     string
	Prompt    string
	Status    string
	RewardId  string
	Cost      int64
}

func GenerateRedemptionBody(p RedemptionParams) (TriggerResponse, error) {
	uuid := util.RandomGUID()
	tNow := util.GetTimestamp().Format(time.RFC3339)
	var event []byte
	var err error

	fromUserName := "testFromuser"

	toUserName := "testBroadcaster"

	// handle default values for possible command line params
	if p.ToUser == "" {
		p.ToUser = util.RandomUserID()
	}

	if p.FromUser == "" {
		p.FromUser = util.RandomUserID()
	}

	if p.Title == "" {
		p.Title = "Test Reward from CLI"
	}

	if p.Prompt == "" {
		p.Prompt = "Redeem Your Test Reward from CLI"
	}

	if p.Status == "" {
		p.Status = "unfulfilled"
	}

	if p.RewardId == "" {
		p.RewardId = util.RandomGUID()
	}

	if p.Cost <= 0 {
		p.Cost = 150
	}

	switch p.Transport {
	case "eventsub":
		body := *&models.EventsubResponse{
			Subscription: models.EventsubSubscription{
				ID:      uuid,
				Type:    p.Type,
				Version: "test",
				Condition: models.EventsubCondition{
					BroadcasterUserID: p.ToUser,
				},
				Transport: models.EventsubTransport{
					Method:   "webhook",
					Callback: "null",
				},
				CreatedAt: tNow,
			},
			Event: models.RedemptionEventSubEvent{
				Id:                  uuid,
				BroadcasterUserId:   p.ToUser,
				BroadcasterUserName: toUserName,
				UserId:              p.FromUser,
				UserName:            fromUserName,
				UserInput:           "Test Input From CLI",
				Status:              p.Status,
				Reward: models.RedemptionReward{
					Id:     p.RewardId,
					Title:  p.Title,
					Cost:   p.Cost,
					Prompt: p.Prompt,
				},
				RedeemedAt: tNow,
			},
		}

		event, err = json.Marshal(body)
		if err != nil {
			return TriggerResponse{}, err
		}

	case "websub":
		return TriggerResponse{}, errors.New("Websub is unsupported for channel points events")
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
