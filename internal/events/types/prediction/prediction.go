// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package prediction

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

var triggerSupported = []string{"prediction-begin", "prediction-progress", "prediction-end", "prediction-lock"}

var triggerMapping = map[string]map[string]string{
	models.TransportEventSub: {
		"prediction-begin":    "channel.prediction.begin",
		"prediction-progress": "channel.prediction.progress",
		"prediction-lock":     "channel.prediction.lock",
		"prediction-end":      "channel.prediction.end",
	},
}

type Event struct{}

func (e Event) GenerateEvent(params events.MockEventParameters) (events.MockEventResponse, error) {
	var event []byte
	var err error
	if params.Description == "" {
		params.Description = "Will the developer finish this program?"
	}

	switch params.Transport {
	case models.TransportEventSub:
		var outcomes []models.PredictionEventSubEventOutcomes
		for i := 0; i < 2; i++ {
			color := "blue"
			title := "yes"

			if i == 1 {
				color = "pink"
				title = "no"
			}

			o := models.PredictionEventSubEventOutcomes{
				ID:    util.RandomGUID(),
				Title: title,
				Color: color,
			}

			if params.Trigger != "prediction-begin" {
				tp := []models.PredictionEventSubEventTopPredictors{}
				sum := 0
				for j := 0; j < int(util.RandomInt(10))+1; j++ {
					t := models.PredictionEventSubEventTopPredictors{
						UserID:            util.RandomUserID(),
						UserLogin:         "testLogin",
						UserName:          "testLogin",
						ChannelPointsUsed: int(util.RandomInt(10*1000)) + 100,
					}
					sum += t.ChannelPointsUsed
					if params.Trigger == "prediction-lock" || params.Trigger == "prediction-end" {
						if i == 0 {
							t.ChannelPointsWon = intPointer(t.ChannelPointsUsed * 2)
						} else {
							t.ChannelPointsWon = intPointer(0)
						}
					}
					tp = append(tp, t)
					o.TopPredictors = &tp
				}
				length := len(*o.TopPredictors)
				o.Users = &length
				o.ChannelPoints = &sum
			}

			outcomes = append(outcomes, o)
		}

		body := &models.PredictionEventSubResponse{
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
			Event: models.PredictionEventSubEvent{
				ID:                   util.RandomGUID(),
				BroadcasterUserID:    params.ToUserID,
				BroadcasterUserLogin: params.ToUserName,
				BroadcasterUserName:  params.ToUserName,
				Title:                params.Description,
				Outcomes:             outcomes,
				StartedAt:            util.GetTimestamp().Format(time.RFC3339Nano),
			},
		}

		if params.Trigger == "prediction-begin" || params.Trigger == "prediction-progress" {
			body.Event.LocksAt = util.GetTimestamp().Add(time.Minute * 10).Format(time.RFC3339Nano)
		} else if params.Trigger == "prediction-lock" {
			body.Event.LockedAt = util.GetTimestamp().Add(time.Minute * 10).Format(time.RFC3339Nano)
		} else if params.Trigger == "prediction-end" {
			body.Event.WinningOutcomeID = outcomes[0].ID
			body.Event.EndedAt = util.GetTimestamp().Add(time.Minute * 10).Format(time.RFC3339Nano)
			body.Event.Status = "resolved"
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

func intPointer(i int) *int {
	return &i
}
func (e Event) GetEventbusAlias(t string) string {
	// check for aliases
	for trigger, topic := range triggerMapping[models.TransportEventSub] {
		if topic == t {
			return trigger
		}
	}
	return ""
}
