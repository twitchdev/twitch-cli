// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package charity

import (
	"encoding/json"
	"testing"

	"github.com/twitchdev/twitch-cli/internal/events"
	"github.com/twitchdev/twitch-cli/internal/models"
	"github.com/twitchdev/twitch-cli/test_setup"
)

var fromUser = "1234"
var toUser = "4567"

func TestEventSubCharity(t *testing.T) {
	testEventSubCharity(t, "charity-donate")
	testEventSubCharity(t, "charity-start")
	testEventSubCharity(t, "charity-progress")
	testEventSubCharity(t, "charity-stop")
}

func testEventSubCharity(t *testing.T, trigger string) {
	a := test_setup.SetupTestEnv(t)
	params := events.MockEventParameters{
		FromUserID:         fromUser,
		ToUserID:           toUser,
		Transport:          models.TransportEventSub,
		Trigger:            trigger,
		SubscriptionStatus: "enabled",
	}

	r, err := Event{}.GenerateEvent(params)
	a.Nil(err, "Error generating body.")

	var body models.CharityEventSubResponse

	err = json.Unmarshal(r.JSON, &body)
	a.Nil(err, "Error unmarshalling JSON")

	a.Equal(toUser, body.Event.BroadcasterUserID, "Expected to user %v, got %v", toUser, body.Event.BroadcasterUserID)
	if trigger == "charity-donate" {
		a.Equal(fromUser, *body.Event.UserID, "Expected from user %v, got %v", r.ToUser, body.Event.UserID)

	}
}

func TestFakeTransport(t *testing.T) {
	a := test_setup.SetupTestEnv(t)

	params := events.MockEventParameters{
		FromUserID: fromUser,
		ToUserID:   toUser,
		Transport:  "fake_transport",
		Trigger:    "charity-donate",
	}

	r, err := Event{}.GenerateEvent(params)
	a.Nil(err)
	a.Empty(r)
}

func TestValidTrigger(t *testing.T) {
	a := test_setup.SetupTestEnv(t)

	r := Event{}.ValidTrigger("charity-donate")
	a.Equal(true, r)

	r = Event{}.ValidTrigger("charity-start")
	a.Equal(true, r)

	r = Event{}.ValidTrigger("charity-progress")
	a.Equal(true, r)

	r = Event{}.ValidTrigger("charity-stop")
	a.Equal(true, r)

	r = Event{}.ValidTrigger("notcharity")
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

	r := Event{}.GetTopic(models.TransportEventSub, "charity-donate")
	a.NotNil(r)

	r = Event{}.GetTopic(models.TransportEventSub, "charity-start")
	a.NotNil(r)

	r = Event{}.GetTopic(models.TransportEventSub, "charity-progress")
	a.NotNil(r)

	r = Event{}.GetTopic(models.TransportEventSub, "charity-stop")
	a.NotNil(r)
}
