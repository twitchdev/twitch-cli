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
var triggers = []string{"charity-donate", "charity-start", "charity-progress", "charity-stop"}

var triggerMapping = map[string]map[string]string{
	models.TransportEventSub: {
		"charity-donate":   "channel.charity_campaign.donate",
		"charity-start":    "channel.charity_campaign.start",
		"charity-progress": "channel.charity_campaign.progress",
		"charity-stop":     "channel.charity_campaign.stop",
	},
}

type Event struct{}

func (e Event) GenerateEvent(params events.MockEventParameters) (events.MockEventResponse, error) {
	var event []byte
	var err error
	var campaign_id *string // only used by channel.charity_campaign.donate
	var id *string          // used by the rest of channel.charity_campaign.*
	var user_id *string
	var user_login_name *string
	var charity_description *string
	var charity_website *string
	var amount *models.CharityEventSubEventAmount
	var current_amount *models.CharityEventSubEventAmount
	var target_amount *models.CharityEventSubEventAmount
	var started_at *string
	var stopped_at *string

	randomID := util.RandomGUID()
	charityName := "Example Charity"
	charityLogo := "https://abc.cloudfront.net/ppgf/1000/100.png"
	charityDescription := "Example Description"
	charityWebsite := "https://www.example.com"

	if params.Trigger == "charity-donate" {
		campaign_id = &randomID
		user_id = &params.FromUserID
		user_login_name = &params.FromUserName
		amount = &models.CharityEventSubEventAmount{
			Value:         10000,
			DecimalPlaces: 2,
			Currency:      "USD",
		}
	}
	if params.Trigger == "charity-start" {
		id = &randomID
		charity_description = &charityDescription
		charity_website = &charityWebsite
		current_amount = &models.CharityEventSubEventAmount{
			Value:         0,
			DecimalPlaces: 2,
			Currency:      "USD",
		}
		target_amount = &models.CharityEventSubEventAmount{
			Value:         1500000,
			DecimalPlaces: 2,
			Currency:      "USD",
		}
		started_at = &params.Timestamp
	}
	if params.Trigger == "charity-progress" {
		id = &randomID
		current_amount = &models.CharityEventSubEventAmount{
			Value:         260000,
			DecimalPlaces: 2,
			Currency:      "USD",
		}
		target_amount = &models.CharityEventSubEventAmount{
			Value:         1500000,
			DecimalPlaces: 2,
			Currency:      "USD",
		}
	}
	if params.Trigger == "charity-stop" {
		id = &randomID
		charity_description = &charityDescription
		charity_website = &charityWebsite
		current_amount = &models.CharityEventSubEventAmount{
			Value:         1450000,
			DecimalPlaces: 2,
			Currency:      "USD",
		}
		target_amount = &models.CharityEventSubEventAmount{
			Value:         1500000,
			DecimalPlaces: 2,
			Currency:      "USD",
		}
		stopped_at = &params.Timestamp
	}

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
				CampaignID:           campaign_id,
				ID:                   id,
				BroadcasterUserID:    params.ToUserID,
				BroadcasterUserName:  params.ToUserName,
				BroadcasterUserLogin: params.ToUserName,
				UserID:               user_id,
				UserName:             user_login_name,
				UserLogin:            user_login_name,
				CharityName:          charityName,
				CharityDescription:   charity_description,
				CharityLogo:          charityLogo,
				CharityWebsite:       charity_website,
				Amount:               amount,
				CurrentAmount:        current_amount,
				TargetAmount:         target_amount,
				StartedAt:            started_at,
				StoppedAt:            stopped_at,
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
