// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package hype_train

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
var triggerSupported = []string{"hype-train-begin", "hype-train-progress", "hype-train-end"}
var triggerMapping = map[string]map[string]string{
	models.TransportEventSub: {
		"hype-train-progress": "channel.hype_train.progress",
		"hype-train-begin":    "channel.hype_train.begin",
		"hype-train-end":      "channel.hype_train.end",
	},
	models.TransportWebSocket: {
		"hype-train-progress": "channel.hype_train.progress",
		"hype-train-begin":    "channel.hype_train.begin",
		"hype-train-end":      "channel.hype_train.end",
	},
}

type Event struct{}

func (e Event) GenerateEvent(params events.MockEventParameters) (events.MockEventResponse, error) {
	var event []byte
	var err error
	lastUser := util.RandomUserID()
	lastTotal := util.RandomInt(10 * 100)
	lastType := util.RandomType()

	//Local variables which will be used for the trigger params below
	localLevel := util.RandomInt(4) + 1
	localTotal := util.RandomInt(10 * 100)
	localGoal := util.RandomInt(10*100*100) + localTotal
	localProgress := localTotal - util.RandomInt(100)

	tNow, _ := time.Parse(time.RFC3339Nano, params.Timestamp)

	switch params.Transport {
	case models.TransportEventSub, models.TransportWebSocket:
		body := models.HypeTrainEventSubResponse{
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
			Event: models.HypeTrainEventSubEvent{
				ID:                   params.ID,
				BroadcasterUserID:    params.ToUserID,
				BroadcasterUserLogin: params.ToUserName,
				BroadcasterUserName:  params.ToUserName,
				Total:                localTotal,
				Progress:             &localProgress,
				Goal:                 localGoal,
				TopContributions: []models.ContributionData{
					{
						TotalContribution:            util.RandomInt(10 * 100),
						TypeOfContribution:           util.RandomType(),
						UserWhoMadeContribution:      util.RandomUserID(),
						UserNameWhoMadeContribution:  "cli_user1",
						UserLoginWhoMadeContribution: "cli_user1",
					},
					{
						TotalContribution:            lastTotal,
						TypeOfContribution:           lastType,
						UserWhoMadeContribution:      lastUser,
						UserNameWhoMadeContribution:  "cli_user2",
						UserLoginWhoMadeContribution: "cli_user2",
					},
				},
				LastContribution: models.ContributionData{
					TotalContribution:            lastTotal,
					TypeOfContribution:           lastType,
					UserWhoMadeContribution:      lastUser,
					UserNameWhoMadeContribution:  "cli_user2",
					UserLoginWhoMadeContribution: "cli_user2",
				},
				Level:              localLevel,
				StartedAtTimestamp: params.Timestamp,
				ExpiresAtTimestamp: tNow.Add(5 * time.Minute).Format(time.RFC3339Nano),
			},
		}
		if params.Trigger == "hype-train-begin" {
			body.Event.Progress = &localTotal
		}
		if params.Trigger == "hype-train-end" {
			body.Event.CooldownEndsAtTimestamp = tNow.Add(1 * time.Hour).Format(time.RFC3339Nano)
			body.Event.EndedAtTimestamp = params.Timestamp
			body.Event.ExpiresAtTimestamp = ""
			body.Event.Goal = 0
			body.Event.Progress = nil
			body.Event.StartedAtTimestamp = tNow.Add(5 * -time.Minute).Format(time.RFC3339Nano)
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
