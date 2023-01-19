// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package shield_mode

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
		Trigger:            "shield-mode-begin",
		SubscriptionStatus: "enabled",
		Cost:               0,
	}
	endParams := *&events.MockEventParameters{
		FromUserID:         fromUser,
		ToUserID:           toUser,
		Transport:          models.TransportEventSub,
		Trigger:            "shield-mode-end",
		SubscriptionStatus: "enabled",
		Cost:               0,
	}

	r1, err := Event{}.GenerateEvent(beginParams)
	a.Nil(err)

	r2, err := Event{}.GenerateEvent(endParams)
	a.Nil(err)

	var body models.ShieldModeEventSubResponse
	err = json.Unmarshal(r1.JSON, &body)
	a.Nil(err)

	err = json.Unmarshal(r2.JSON, &body)
	a.Nil(err)
}

func TestFakeTransport(t *testing.T) {
	a := test_setup.SetupTestEnv(t)

	beginParams := *&events.MockEventParameters{
		FromUserID:         fromUser,
		ToUserID:           toUser,
		Transport:          "fake_transport",
		Trigger:            "shield-mode-begin",
		SubscriptionStatus: "enabled",
	}
	endParams := *&events.MockEventParameters{
		FromUserID:         fromUser,
		ToUserID:           toUser,
		Transport:          "fake_transport",
		Trigger:            "shield-mode-end",
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

	r := Event{}.ValidTrigger("shield-mode-begin")
	a.Equal(true, r)

	r = Event{}.ValidTrigger("shield-mode-end")
	a.Equal(true, r)

	r = Event{}.ValidTrigger("notshield")
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

	r := Event{}.GetTopic(models.TransportEventSub, "shield-mode-begin")
	a.NotNil(r)

	r = Event{}.GetTopic(models.TransportEventSub, "shield-mode-end")
	a.NotNil(r)
}
