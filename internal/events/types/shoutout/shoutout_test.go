// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package shoutout

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

	beginParams := *&events.MockEventParameters{
		FromUserID:         fromUser,
		ToUserID:           toUser,
		Transport:          models.TransportEventSub,
		Trigger:            "shoutout-create",
		SubscriptionStatus: "enabled",
		Cost:               0,
	}
	endParams := *&events.MockEventParameters{
		FromUserID:         fromUser,
		ToUserID:           toUser,
		Transport:          models.TransportEventSub,
		Trigger:            "shoutout-received",
		SubscriptionStatus: "enabled",
		Cost:               0,
	}

	r1, err := Event{}.GenerateEvent(beginParams)
	a.Nil(err)

	r2, err := Event{}.GenerateEvent(endParams)
	a.Nil(err)

	var body1 models.ShoutoutCreateEventSubResponse
	err = json.Unmarshal(r1.JSON, &body1)
	a.Nil(err)

	var body2 models.ShoutoutReceivedEventSubResponse
	err = json.Unmarshal(r2.JSON, &body2)
	a.Nil(err)
}

func TestFakeTransport(t *testing.T) {
	a := test_setup.SetupTestEnv(t)

	beginParams := *&events.MockEventParameters{
		FromUserID:         fromUser,
		ToUserID:           toUser,
		Transport:          "fake_transport",
		Trigger:            "shoutout-create",
		SubscriptionStatus: "enabled",
	}
	endParams := *&events.MockEventParameters{
		FromUserID:         fromUser,
		ToUserID:           toUser,
		Transport:          "fake_transport",
		Trigger:            "shoutout-received",
		SubscriptionStatus: "enabled",
	}

	r1, err1 := Event{}.GenerateEvent(beginParams)
	r2, err2 := Event{}.GenerateEvent(endParams)
	a.Nil(err1)
	a.Nil(err2)
	a.Empty(r1)
	a.Empty(r2)
}
func TestValidTrigger(t *testing.T) {
	a := test_setup.SetupTestEnv(t)

	r := Event{}.ValidTrigger("shoutout-create")
	a.Equal(true, r)

	r = Event{}.ValidTrigger("shoutout-received")
	a.Equal(true, r)

	r = Event{}.ValidTrigger("notshoutout")
	a.Equal(false, r)
}

func TestValidTransport(t *testing.T) {
	a := test_setup.SetupTestEnv(t)

	r := Event{}.ValidTransport(models.TransportEventSub)
	a.Equal(true, r)

	r = Event{}.ValidTransport("noteventsub")
	a.Equal(false, r)
}
func TestGetTopic(t *testing.T) {
	a := test_setup.SetupTestEnv(t)

	r := Event{}.GetTopic(models.TransportEventSub, "shoutout-create")
	a.NotNil(r)

	r = Event{}.GetTopic(models.TransportEventSub, "shoutout-receieve")
	a.NotNil(r)
}
