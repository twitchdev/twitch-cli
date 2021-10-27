// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package goal

import (
	"encoding/json"
	"time"

	"github.com/twitchdev/twitch-cli/internal/events"
	"github.com/twitchdev/twitch-cli/internal/models"
	"github.com/twitchdev/twitch-cli/internal/util"
)

var transportsSupported = map[string]bool{
	models.TransportWebSub:   false,
	models.TransportEventSub: true,
}

var triggerSupported = []string{"goal-begin", "goal-progress", "goal-end"}

var triggerMapping = map[string]map[string]string{
	models.TransportWebSub: {},
	models.TransportEventSub: {
		"goal-progress": "channel.goal.progress",
		"goal-begin":    "channel.goal.begin",
		"goal-end":      "channel.goal.end",
	},
}

type Event struct{}

func (e Event) GenerateEvent(params events.MockEventParameters) (events.MockEventResponse, error) {
	var event []byte
	var err error
	var isAchieved *bool
	var goalEndDate *string
	var goalType string
	var currentAmount int64
	var targetAmount int64

	goalStartedAt := util.GetTimestamp()
	currentAmount = util.RandomInt(10 * 10)
	targetAmount = util.RandomInt(10 * 100)

	if params.Trigger == "goal-end" {
		endDate := goalStartedAt.Add(time.Hour * 24).Format(time.RFC3339)
		goalEndDate = &endDate

		achieved := util.RandomInt(1) == 1
		if achieved {
			currentAmount = 100
			targetAmount = 100
		}

		isAchieved = &achieved
	}

	goalType = params.ItemName
	if goalType == "" {
		goalType = "follower"
	}

	switch params.Transport {
	case models.TransportEventSub:

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
			Event: models.GoalEventSubEvent{
				ID:                   params.ID,
				BroadcasterUserID:    params.ToUserID,
				BroadcasterUserLogin: params.ToUserName,
				BroadcasterUserName:  params.ToUserName,
				Type:                 goalType,
				Description:          params.Description,
				CurrentAmount:        currentAmount,
				TargetAmount:         targetAmount,
				StartedAt:            goalStartedAt.Format(time.RFC3339Nano),
				EndedAt:              goalEndDate,
				IsAchieved:           isAchieved,
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
