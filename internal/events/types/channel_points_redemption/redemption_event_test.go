// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package channel_points_redemption

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
		Transport:          models.TransportWebhook,
		Trigger:            "add-redemption",
		SubscriptionStatus: "enabled",
		ToUserID:           toUser,
		FromUserID:         fromUser,
		EventStatus:        "tested",
		ItemID:             "12345678-1234-abcd-5678-000000000000",
		Cost:               1337,
		ItemName:           "Testing",
	}

	r, err := Event{}.GenerateEvent(params)
	a.Nil(err)

	var body models.RedemptionEventSubResponse
	err = json.Unmarshal(r.JSON, &body)
	a.Nil(err)

	a.Equal(toUser, body.Event.BroadcasterUserID, "Expected to user %v, got %v", toUser, body.Event.BroadcasterUserID)
	a.Equal(fromUser, body.Event.UserID, "Expected from user %v, got %v", r.ToUser, body.Event.UserID)
	a.Equal(params.EventStatus, body.Event.Status)
	a.Equal(params.Cost, body.Event.Reward.Cost)
	a.Equal(params.ItemID, body.Event.Reward.ID)
	a.Equal(params.ItemName, body.Event.Reward.Title)

	params = events.MockEventParameters{
		Transport:  models.TransportWebhook,
		ToUserID:   toUser,
		FromUserID: fromUser,
	}

	r, err = Event{}.GenerateEvent(params)
	a.Nil(err)

	err = json.Unmarshal(r.JSON, &body)
	a.Nil(err)

	a.NotNil(body.Event.BroadcasterUserID)
	a.NotNil(body.Event.UserID)
	a.NotNil(body.Event.Reward.ID)
}

func TestFakeTransport(t *testing.T) {
	a := test_setup.SetupTestEnv(t)

	params := events.MockEventParameters{
		FromUserID:         fromUser,
		ToUserID:           toUser,
		Transport:          "fake_transport",
		Trigger:            "add-redemption",
		SubscriptionStatus: "enabled",
	}

	r, err := Event{}.GenerateEvent(params)
	a.Nil(err)
	a.Empty(r)
}
func TestValidTrigger(t *testing.T) {
	a := test_setup.SetupTestEnv(t)

	r := Event{}.ValidTrigger("add-redemption")
	a.Equal(true, r)

	r = Event{}.ValidTrigger("not_trigger_keyword")
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

	r := Event{}.GetTopic(models.TransportWebhook, "add-redemption")
	a.NotNil(r)
}
