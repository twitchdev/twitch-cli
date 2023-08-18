// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package hype_train

import (
	"encoding/json"
	"testing"

	"github.com/twitchdev/twitch-cli/internal/events"
	"github.com/twitchdev/twitch-cli/internal/models"
	"github.com/twitchdev/twitch-cli/test_setup"
)

var toUser = "4567"

func TestEventSub(t *testing.T) {
	a := test_setup.SetupTestEnv(t)

	params := events.MockEventParameters{
		ToUserID:           toUser,
		Transport:          models.TransportWebhook,
		Trigger:            "hype-train-begin",
		SubscriptionStatus: "enabled",
	}

	r, err := Event{}.GenerateEvent(params)
	a.Nil(err)

	var body models.HypeTrainEventSubResponse
	err = json.Unmarshal(r.JSON, &body)
	a.Nil(err)

	a.Equal("channel.hype_train.begin", body.Subscription.Type, "Expected event type %v, got %v", "channel.hype_train.begin", body.Subscription.Type)
	a.Equal(toUser, body.Event.BroadcasterUserID, "Expected to user %v, got %v", toUser, body.Event.BroadcasterUserID)

	params = events.MockEventParameters{
		ToUserID:           toUser,
		Transport:          models.TransportWebhook,
		Trigger:            "hype-train-progress",
		SubscriptionStatus: "enabled",
	}

	r, err = Event{}.GenerateEvent(params)
	a.Nil(err)

	//var body models.HypeTrainEventProgressSubResponse
	err = json.Unmarshal(r.JSON, &body)
	a.Nil(err)

	a.Equal("channel.hype_train.progress", body.Subscription.Type, "Expected event type %v, got %v", "channel.hype_train.progress", body.Subscription.Type)
	a.Equal(toUser, body.Event.BroadcasterUserID, "Expected to user %v, got %v", toUser, body.Event.BroadcasterUserID)

	params = events.MockEventParameters{
		ToUserID:  toUser,
		Transport: models.TransportWebhook,
		Trigger:   "hype-train-end",
	}

	r, err = Event{}.GenerateEvent(params)
	a.Nil(err)

	//var body models.HypeTrainEventProgressSubResponse
	err = json.Unmarshal(r.JSON, &body)
	a.Nil(err)

	a.Equal("channel.hype_train.end", body.Subscription.Type, "Expected event type %v, got %v", "channel.hype_train.end", body.Subscription.Type)
	a.Equal(toUser, body.Event.BroadcasterUserID, "Expected to user %v, got %v", toUser, body.Event.BroadcasterUserID)

}

func TestFakeTransport(t *testing.T) {
	a := test_setup.SetupTestEnv(t)

	params := events.MockEventParameters{
		ToUserID:  toUser,
		Transport: "fake_transport",
		Trigger:   "hype-train-progress",
	}

	r, err := Event{}.GenerateEvent(params)
	a.Nil(err)
	a.Empty(r)
}
func TestValidTrigger(t *testing.T) {
	a := test_setup.SetupTestEnv(t)

	r := Event{}.ValidTrigger("hype-train-begin")
	a.Equal(true, r)

	r = Event{}.ValidTrigger("hype-train-progress")
	a.Equal(true, r)

	r = Event{}.ValidTrigger("hype-train-end")
	a.Equal(true, r)

}

func TestValidTransport(t *testing.T) {
	a := test_setup.SetupTestEnv(t)

	r := Event{}.ValidTransport(models.TransportWebhook)
	a.Equal(true, r)
}

func TestGetTopic(t *testing.T) {
	a := test_setup.SetupTestEnv(t)

	r := Event{}.GetTopic(models.TransportWebhook, "hype-train-progress")
	a.Equal("channel.hype_train.progress", r, "Expected %v, got %v", "channel.hype_train.progress", r)
}
