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
	var endDate *string
	var goalType string

	createdAt := util.GetTimestamp()
	switch params.Trigger {
	case "goal-end":
		date := createdAt.Add(time.Hour * 24).Format(time.RFC3339)
		endDate = &date

		achieved := util.RandomInt(1) == 1
		isAchieved = &achieved
	}

	goalType = params.ItemID
	if goalType == "" {
		goalType = "follower"
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
		Event: models.GoalEventSubEvent{
			ID:                   params.ID,
			BroadcasterUserID:    params.ToUserID,
			BroadcasterUserLogin: params.ToUserName,
			BroadcasterUserName:  params.ToUserName,
			Type:                 goalType,
			Description:          params.Description,
			CurrentAmount:        util.RandomInt(10 * 10),
			TargetAmount:         util.RandomInt(10 * 100),
			StartedAt:            util.GetTimestamp().Format(time.RFC3339Nano),
			EndedAt:              endDate,
			IsAchieved:           isAchieved,
		},
	}

	event, err = json.Marshal(body)
	if err != nil {
		return events.MockEventResponse{}, err
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
