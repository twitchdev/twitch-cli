// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package goal

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/twitchdev/twitch-cli/internal/events"
	"github.com/twitchdev/twitch-cli/internal/models"
	"github.com/twitchdev/twitch-cli/internal/util"
)

var transportsSupported = map[string]bool{
	models.TransportEventSub:  true,
	models.TransportWebSocket: true,
}

var triggerSupported = []string{"goal-begin", "goal-progress", "goal-end"}

var triggerMapping = map[string]map[string]string{
	models.TransportEventSub: {
		"goal-progress": "channel.goal.progress",
		"goal-begin":    "channel.goal.begin",
		"goal-end":      "channel.goal.end",
	},
	models.TransportWebSocket: {
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

	goalStartedAt := params.Timestamp
	currentAmount = util.RandomInt(10 * 10)
	targetAmount = util.RandomInt(10 * 100)

	if params.Trigger == "goal-end" {
		tNow, _ := time.Parse(time.RFC3339Nano, params.Timestamp)
		endDate := tNow.Add(time.Hour * 24).Format(time.RFC3339)
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
	case models.TransportEventSub, models.TransportWebSocket:

		body := *&models.EventsubResponse{
			Subscription: models.EventsubSubscription{
				ID:      params.ID,
				Status:  params.SubscriptionStatus,
				Type:    triggerMapping[params.Transport][params.Trigger],
				Version: e.SubscriptionVersion(),
				Condition: models.EventsubCondition{
					BroadcasterUserID: params.ToUserID,
				},
				Transport: models.EventsubTransport{
					Method:   "webhook",
					Callback: "null",
				},
				Cost:      0,
				CreatedAt: params.Timestamp,
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
				StartedAt:            goalStartedAt,
				EndedAt:              goalEndDate,
				IsAchieved:           isAchieved,
			},
		}

		event, err = json.Marshal(body)
		if err != nil {
			return events.MockEventResponse{}, err
		}

		// Delete event info if Subscription.Status is not set to "enabled"
		if !strings.EqualFold(params.SubscriptionStatus, "enabled") {
			var i interface{}
			if err := json.Unmarshal([]byte(event), &i); err != nil {
				return events.MockEventResponse{}, err
			}
			if m, ok := i.(map[string]interface{}); ok {
				delete(m, "event") // Matches JSON key defined in body variable above
			}

			event, err = json.Marshal(i)
			if err != nil {
				return events.MockEventResponse{}, err
			}
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

func (e Event) ValidTrigger(trigger string) bool {
	for _, t := range triggerSupported {
		if t == trigger {
			return true
		}
	}
	return false
}

func (e Event) GetTopic(transport string, trigger string) string {
	return triggerMapping[transport][trigger]
}

func (e Event) GetAllTopicsByTransport(transport string) []string {
	allTopics := []string{}
	for _, topic := range triggerMapping[transport] {
		allTopics = append(allTopics, topic)
	}
	return allTopics
}

func (e Event) GetEventSubAlias(t string) string {
	// check for aliases
	for trigger, topic := range triggerMapping[models.TransportEventSub] {
		if topic == t {
			return trigger
		}
	}
	return ""
}

func (e Event) SubscriptionVersion() string {
	return "1"
}
