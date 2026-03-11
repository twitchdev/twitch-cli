// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package hype_train_v2

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

	var body models.HypeTrainEventSubResponseV2
	err = json.Unmarshal(r.JSON, &body)
	a.Nil(err)

	a.Equal("channel.hype_train.begin", body.Subscription.Type, "Expected event type %v, got %v", "channel.hype_train.begin", body.Subscription.Type)
	a.Equal(toUser, body.Event.BroadcasterUserID, "Expected to user %v, got %v", toUser, body.Event.BroadcasterUserID)
	a.Equal("2", body.Subscription.Version, "Expected version 2, got %v", body.Subscription.Version)
	a.Equal("regular", body.Event.Type, "Expected type regular, got %v", body.Event.Type)
	a.Equal(false, body.Event.IsSharedTrain, "Expected is_shared_train false")
	a.Nil(body.Event.SharedTrainParticipants, "Expected shared_train_participants nil")
	a.NotZero(body.Event.AllTimeHighLevel, "Expected all_time_high_level to be set for begin event")
	a.NotZero(body.Event.AllTimeHighTotal, "Expected all_time_high_total to be set for begin event")
	a.NotNil(body.Event.Progress, "Expected progress to be set for begin event")

	params = events.MockEventParameters{
		ToUserID:           toUser,
		Transport:          models.TransportWebhook,
		Trigger:            "hype-train-progress",
		SubscriptionStatus: "enabled",
	}

	r, err = Event{}.GenerateEvent(params)
	a.Nil(err)

	body = models.HypeTrainEventSubResponseV2{}
	err = json.Unmarshal(r.JSON, &body)
	a.Nil(err)

	a.Equal("channel.hype_train.progress", body.Subscription.Type, "Expected event type %v, got %v", "channel.hype_train.progress", body.Subscription.Type)
	a.Equal(toUser, body.Event.BroadcasterUserID, "Expected to user %v, got %v", toUser, body.Event.BroadcasterUserID)
	a.Equal("2", body.Subscription.Version)
	a.NotNil(body.Event.Progress, "Expected progress to be set for progress event")

	params = events.MockEventParameters{
		ToUserID:           toUser,
		Transport:          models.TransportWebhook,
		Trigger:            "hype-train-end",
		SubscriptionStatus: "enabled",
	}

	r, err = Event{}.GenerateEvent(params)
	a.Nil(err)

	body = models.HypeTrainEventSubResponseV2{}
	err = json.Unmarshal(r.JSON, &body)
	a.Nil(err)

	a.Equal("channel.hype_train.end", body.Subscription.Type, "Expected event type %v, got %v", "channel.hype_train.end", body.Subscription.Type)
	a.Equal(toUser, body.Event.BroadcasterUserID, "Expected to user %v, got %v", toUser, body.Event.BroadcasterUserID)
	a.Equal("2", body.Subscription.Version)
	a.Nil(body.Event.Progress, "Expected progress to be nil for end event")
	a.NotEmpty(body.Event.EndedAtTimestamp, "Expected ended_at to be set for end event")
	a.NotEmpty(body.Event.CooldownEndsAtTimestamp, "Expected cooldown_ends_at to be set for end event")
}

func TestWebSocketTransport(t *testing.T) {
	a := test_setup.SetupTestEnv(t)

	params := events.MockEventParameters{
		ToUserID:           toUser,
		Transport:          models.TransportWebSocket,
		Trigger:            "hype-train-begin",
		SubscriptionStatus: "enabled",
	}

	r, err := Event{}.GenerateEvent(params)
	a.Nil(err)

	var body models.HypeTrainEventSubResponseV2
	err = json.Unmarshal(r.JSON, &body)
	a.Nil(err)

	a.Equal("channel.hype_train.begin", body.Subscription.Type)
	a.Equal("2", body.Subscription.Version)
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

	r = Event{}.ValidTransport(models.TransportWebSocket)
	a.Equal(true, r)
}

func TestGetTopic(t *testing.T) {
	a := test_setup.SetupTestEnv(t)

	r := Event{}.GetTopic(models.TransportWebhook, "hype-train-progress")
	a.Equal("channel.hype_train.progress", r, "Expected %v, got %v", "channel.hype_train.progress", r)
}
