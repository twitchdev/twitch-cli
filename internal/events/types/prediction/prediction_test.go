// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package prediction

import (
	"encoding/json"
	"testing"

	"github.com/twitchdev/twitch-cli/internal/events"
	"github.com/twitchdev/twitch-cli/internal/models"
	"github.com/twitchdev/twitch-cli/test_setup"
)

var fromUser = "1234"
var toUser = "4567"

func TestEventSub(t *testing.T) {
	a := test_setup.SetupTestEnv(t)

	params := events.MockEventParameters{
		FromUserID:         fromUser,
		ToUserID:           toUser,
		Transport:          models.TransportWebhook,
		Trigger:            "prediction-begin",
		SubscriptionStatus: "enabled",
	}

	r, err := Event{}.GenerateEvent(params)
	a.Nil(err)

	var body models.PredictionEventSubResponse
	err = json.Unmarshal(r.JSON, &body)
	a.Nil(err)

	params = events.MockEventParameters{
		FromUserID:         fromUser,
		ToUserID:           toUser,
		Transport:          models.TransportWebhook,
		Trigger:            "prediction-progress",
		SubscriptionStatus: "enabled",
	}

	r, err = Event{}.GenerateEvent(params)
	a.Nil(err)

	body = models.PredictionEventSubResponse{}
	err = json.Unmarshal(r.JSON, &body)
	a.Nil(err)

	params = events.MockEventParameters{
		FromUserID: fromUser,
		ToUserID:   toUser,
		Transport:  models.TransportWebhook,
		Trigger:    "prediction-lock",
	}

	r, err = Event{}.GenerateEvent(params)
	a.Nil(err)

	body = models.PredictionEventSubResponse{}
	err = json.Unmarshal(r.JSON, &body)
	a.Nil(err)

	params = events.MockEventParameters{
		FromUserID: fromUser,
		ToUserID:   toUser,
		Transport:  models.TransportWebhook,
		Trigger:    "prediction-end",
	}

	r, err = Event{}.GenerateEvent(params)
	a.Nil(err)

	body = models.PredictionEventSubResponse{}
	err = json.Unmarshal(r.JSON, &body)
	a.Nil(err)
}

func TestFakeTransport(t *testing.T) {
	a := test_setup.SetupTestEnv(t)

	params := events.MockEventParameters{
		FromUserID: fromUser,
		ToUserID:   toUser,
		Transport:  "fake_transport",
		Trigger:    "unsubscribe",
	}

	r, err := Event{}.GenerateEvent(params)
	a.Nil(err)
	a.Empty(r)
}
func TestValidTrigger(t *testing.T) {
	a := test_setup.SetupTestEnv(t)

	r := Event{}.ValidTrigger("prediction-begin")
	a.Equal(true, r)

	r = Event{}.ValidTrigger("notgift")
	a.Equal(false, r)
}

func TestValidTransport(t *testing.T) {
	a := test_setup.SetupTestEnv(t)

	r := Event{}.ValidTransport(models.TransportWebhook)
	a.Equal(true, r)

	r = Event{}.ValidTransport("noteventsub")
	a.Equal(false, r)
}
func TestGetTopic(t *testing.T) {
	a := test_setup.SetupTestEnv(t)

	r := Event{}.GetTopic(models.TransportWebhook, "prediction-begin")
	a.NotNil(r)
}
