// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package charity

import (
	"encoding/json"
	"strings"

	"github.com/twitchdev/twitch-cli/internal/events"
	"github.com/twitchdev/twitch-cli/internal/models"
	"github.com/twitchdev/twitch-cli/internal/util"
)

var transportsSupported = map[string]bool{
	models.TransportEventSub: true,
}
var triggers = []string{"charity"}

var triggerMapping = map[string]map[string]string{
	models.TransportEventSub: {
		"charity": "channel.charity_campaign.donate",
	},
}

type Event struct{}

func (e Event) GenerateEvent(params events.MockEventParameters) (events.MockEventResponse, error) {
	var event []byte
	var err error

	switch params.Transport {
	case models.TransportEventSub:
		body := models.EventsubResponse{
			Subscription: models.EventsubSubscription{
				ID:      params.ID,
				Type:    "channel.charity_campaign.donate",
				Version: e.SubscriptionVersion(),
				Status:  params.SubscriptionStatus,
				Cost:    0,
				Condition: models.EventsubCondition{
					BroadcasterUserID: params.ToUserID,
				},
				Transport: models.EventsubTransport{
					Method:   "webhook",
					Callback: "null",
				},
				CreatedAt: params.Timestamp,
			},
			Event: models.CharityEventSubEvent{
				CampaignID:           util.RandomGUID(),
				BroadcasterUserID:    params.ToUserID,
				BroadcasterUserName:  params.ToUserName,
				BroadcasterUserLogin: params.ToUserName,
				UserID:               params.FromUserID,
				UserName:             params.FromUserName,
				UserLogin:            params.FromUserName,
				CharityName:          "Example Charity",
				CharityLogo:          "https://abc.cloudfront.net/ppgf/1000/100.png",
				Amount: models.CharityEventSubEventAmount{
					Value:         10000,
					DecimalPlaces: 2,
					Currency:      "USD",
				},
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
		ID:     params.ID,
		JSON:   event,
		ToUser: params.ToUserID,
	}, nil
}

func (e Event) ValidTransport(transport string) bool {
	return transportsSupported[transport]
}

func (e Event) ValidTrigger(trigger string) bool {
	for _, t := range triggers {
		if t == trigger {
			return true
		}
	}
	return false
}
func (e Event) GetTopic(transport string, trigger string) string {
	return triggerMapping[transport][trigger]
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
	return "beta"
}
