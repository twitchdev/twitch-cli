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

type RewardParams struct {
	Transport string
	Type      string
	ToUser    string
	Title     string
	Prompt    string
	Cost      int64
}

func GenerateRewardBody(p RewardParams) (TriggerResponse, error) {
	uuid := util.RandomGUID()
	tNow := util.GetTimestamp().Format(time.RFC3339)
	var event []byte
	var err error

	toUserName := "testBroadcaster"

	// handle default values for possible command line params
	if p.ToUser == "" {
		p.ToUser = util.RandomUserID()
	}

	if p.Title == "" {
		p.Title = "Test Reward from CLI"
	}

	if p.Prompt == "" {
		p.Prompt = "Redeem Your Test Reward from CLI"
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
			Event: models.RewardEventSubEvent{
				Id:                                uuid,
				BroadcasterUserId:                 p.ToUser,
				BroadcasterUserName:               toUserName,
				IsEnabled:                         true,
				IsPaused:                          false,
				IsInStock:                         true,
				Title:                             p.Title,
				Cost:                              p.Cost,
				Prompt:                            p.Prompt,
				IsUserInputRequired:               true,
				ShouldRedemptionsSkipRequestQueue: false,
				CooldownExpiresAt:                 tNow,
				RedemptionsRedeemedCurrentStream:  0,
				MaxPerStream:                      models.RewardMax{
					IsEnabled: true,
					Value:     100,
				},
				MaxPerUserPerStream:               models.RewardMax{
					IsEnabled: true,
					Value:     100,
				},
				GlobalCooldown:                    models.RewardGlobalCooldown{
					IsEnabled: true,
					Value:     300,
				},
				BackgroundColor:                   "#c0ffee",
				Image:                             models.RewardImage{
					Url1x: "https://static-cdn.jtvnw.net/image-1.png",
					Url2x: "https://static-cdn.jtvnw.net/image-2.png",
					Url4x: "https://static-cdn.jtvnw.net/image-4.png",
				},
				DefaultImage:                      models.RewardImage{
					Url1x: "https://static-cdn.jtvnw.net/default-1.png",
					Url2x: "https://static-cdn.jtvnw.net/default-2.png",
					Url4x: "https://static-cdn.jtvnw.net/default-4.png",
				},
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
		FromUser: p.ToUser,
		ToUser:   p.ToUser,
	}, nil
}
