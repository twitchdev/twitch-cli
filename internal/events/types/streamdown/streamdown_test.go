// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package streamdown

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
		Trigger:            "streamdown",
		SubscriptionStatus: "enabled",
	}

	r, err := Event{}.GenerateEvent(params)
	a.Nil(err)

	var body models.StreamDownEventSubResponse
	err = json.Unmarshal(r.JSON, &body)
	a.Nil(err)

	// write actual tests here (making sure you set appropriate values and the like) for eventsub
}

func TestFakeTransport(t *testing.T) {
	a := test_setup.SetupTestEnv(t)

	params := events.MockEventParameters{
		FromUserID:         fromUser,
		ToUserID:           toUser,
		Transport:          "fake_transport",
		Trigger:            "unsubscribe",
		SubscriptionStatus: "enabled",
	}

	r, err := Event{}.GenerateEvent(params)
	a.Nil(err)
	a.Empty(r)
}
func TestValidTrigger(t *testing.T) {
	a := test_setup.SetupTestEnv(t)

	r := Event{}.ValidTrigger("streamdown")
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

	r := Event{}.GetTopic(models.TransportWebhook, "streamdown")
	a.NotNil(r)
}
