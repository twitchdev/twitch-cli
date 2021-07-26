// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package subscription_message

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
	ten := 10

	params := *&events.MockEventParameters{
		FromUserID: fromUser,
		ToUserID:   toUser,
		Transport:  models.TransportEventSub,
		Trigger:    "subscribe-message",
		Cost:       int64(ten),
	}

	r, err := Event{}.GenerateEvent(params)
	a.Nil(err)

	var body models.SubscribeMessageEventSubResponse // replace with actual value
	err = json.Unmarshal(r.JSON, &body)
	a.Nil(err)
	a.Equal(&ten, body.Event.StreakMonths)
	a.GreaterOrEqual(body.Event.CumulativeMonths, 10)

	params = *&events.MockEventParameters{
		FromUserID:  fromUser,
		ToUserID:    toUser,
		Transport:   models.TransportEventSub,
		Trigger:     "subscribe-message",
		Cost:        int64(ten),
		IsAnonymous: true,
	}

	r, err = Event{}.GenerateEvent(params)
	a.Nil(err)

	err = json.Unmarshal(r.JSON, &body)
	a.Nil(err)
	a.Nil(body.Event.StreakMonths)
	a.GreaterOrEqual(body.Event.CumulativeMonths, 10)
}

func TestFakeTransport(t *testing.T) {
	a := test_setup.SetupTestEnv(t)

	params := *&events.MockEventParameters{
		FromUserID: fromUser,
		ToUserID:   toUser,
		Transport:  "fake_transport",
		Trigger:    "subscribe-message",
	}

	r, err := Event{}.GenerateEvent(params)
	a.Nil(err)
	a.Empty(r)
}
func TestValidTrigger(t *testing.T) {
	a := test_setup.SetupTestEnv(t)

	r := Event{}.ValidTrigger("subscribe-message")
	a.Equal(true, r)

	r = Event{}.ValidTrigger("notmessage")
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

	r := Event{}.GetTopic(models.TransportEventSub, "subscribe-message")
	a.NotNil(r)
}
