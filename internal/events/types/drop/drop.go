// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package drop

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/twitchdev/twitch-cli/internal/events"
	"github.com/twitchdev/twitch-cli/internal/models"
	"github.com/twitchdev/twitch-cli/internal/util"
)

var transportsSupported = map[string]bool{
	models.TransportEventSub: true,
}

var triggerSupported = []string{"drop"}

var triggerMapping = map[string]map[string]string{
	models.TransportEventSub: {
		"drop": "drop.entitlement.grant",
	},
}

type Event struct{}

func (e Event) GenerateEvent(params events.MockEventParameters) (events.MockEventResponse, error) {
	var event []byte
	var err error

	if params.ItemID == "" {
		params.ItemID = util.RandomGUID()
	}

	if params.Description == "" {
		params.Description = fmt.Sprintf("%v", util.RandomInt(1000))
	}
	switch params.Transport {
	case models.TransportEventSub:
		body := &models.DropsEntitlementEventSubResponse{
			Subscription: models.EventsubSubscription{
				ID:      params.ID,
				Status:  params.SubscriptionStatus,
				Type:    triggerMapping[params.Transport][params.Trigger],
				Version: e.SubscriptionVersion(),
				Condition: models.EventsubCondition{
					OrganizationID: params.FromUserID,
				},
				Transport: models.EventsubTransport{
					Method:   "webhook",
					Callback: "null",
				},
				Cost:      0,
				CreatedAt: params.Timestamp,
			},
			Events: []models.DropsEntitlementEventSubEvent{
				{
					ID: util.RandomGUID(),
					Data: models.DropsEntitlementEventSubEventData{
						OrganizationID: params.FromUserID,
						CategoryID:     params.GameID,
						CategoryName:   "",
						CampaignID:     util.RandomGUID(),
						EntitlementID:  util.RandomGUID(),
						BenefitID:      params.ItemID,
						UserID:         params.ToUserID,
						UserName:       params.ToUserName,
						UserLogin:      params.ToUserName,
						CreatedAt:      params.Timestamp,
					},
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
				delete(m, "events") // Matches JSON key defined in body variable above
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
