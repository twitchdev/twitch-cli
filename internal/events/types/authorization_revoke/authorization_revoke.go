// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package authorization_revoke

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/spf13/viper"
	"github.com/twitchdev/twitch-cli/internal/events"
	"github.com/twitchdev/twitch-cli/internal/models"
	"github.com/twitchdev/twitch-cli/internal/util"
)

var transportsSupported = map[string]bool{
	models.TransportWebSub:   false,
	models.TransportEventSub: true,
}

var triggerSupported = []string{"revoke"}

var triggerMapping = map[string]map[string]string{
	models.TransportEventSub: {
		"trigger_keyword": "revoke",
	},
}

type Event struct{}

func (e Event) GenerateEvent(params events.MockEventParameters) (events.MockEventResponse, error) {
	var event []byte
	var err error
	clientID := viper.GetString("ClientID")

	// if not configured, generate a random one
	if clientID == "" {
		clientID = util.RandomClientID()
	}
	switch params.Transport {
	case models.TransportEventSub:
		body := &models.AuthorizationRevokeEventSubResponse{
			Subscription: models.EventsubSubscription{
				ID:      params.ID,
				Status:  "enabled",
				Type:    triggerMapping[params.Transport][params.Trigger],
				Version: "1",
				Condition: models.EventsubCondition{
					ClientID: clientID,
				},
				Transport: models.EventsubTransport{
					Method:   "webhook",
					Callback: "null",
				},
				CreatedAt: util.GetTimestamp().Format(time.RFC3339Nano),
			},
			Event: models.AuthorizationRevokeEvent{
				ClientID:  clientID,
				UserID:    params.FromUserID,
				UserLogin: params.FromUserName,
				UserName:  params.FromUserName,
			},
		}
		event, err = json.Marshal(body)
		if err != nil {
			return events.MockEventResponse{}, err
		}
	case models.TransportWebSub:
		return events.MockEventResponse{}, errors.New("Websub is unsupported for authorization revoke events")
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
