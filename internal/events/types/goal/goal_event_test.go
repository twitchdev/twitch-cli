// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package goal

import (
	"encoding/json"
	"testing"

	"github.com/twitchdev/twitch-cli/internal/events"
	"github.com/twitchdev/twitch-cli/internal/models"
	"github.com/twitchdev/twitch-cli/test_setup"
)

var user = "1234"

func TestEventSub(t *testing.T) {
	a := test_setup.SetupTestEnv(t)

	params := *&events.MockEventParameters{
		ToUserID:           user,
		Description:        "Twitch Subscriber Goal",
		Transport:          models.TransportEventSub,
		Trigger:            "goal-begin",
		SubscriptionStatus: "enabled",
		EventStatus:        "subscriber",
	}

	r, err := Event{}.GenerateEvent(params)
	a.Nil(err)

	var body models.GoalEventSubResponse
	err = json.Unmarshal(r.JSON, &body)
	a.Nil(err)

	a.Equal("channel.goal.begin", body.Subscription.Type, "Expected event type %v, got %v", "channel.goal.begin", body.Subscription.Type)
	a.Equal(user, body.Event.BroadcasterUserID, "Expected from user %v, got %v", r.ToUser, body.Event.BroadcasterUserID)
	a.Equal("Twitch Subscriber Goal", body.Event.Description, "Expected from goal type %v, got %v", "Twitch Subscriber Goal", body.Event.Type)

	params = *&events.MockEventParameters{
		ToUserID:           user,
		Description:        "Twitch Follower Goal",
		Transport:          models.TransportEventSub,
		Trigger:            "goal-progress",
		SubscriptionStatus: "enabled",
	}

	r, err = Event{}.GenerateEvent(params)
	a.Nil(err)

	err = json.Unmarshal(r.JSON, &body)
	a.Nil(err)

	a.Equal("channel.goal.progress", body.Subscription.Type, "Expected event type %v, got %v", "channel.goal.progress", body.Subscription.Type)
	a.Equal(user, body.Event.BroadcasterUserID, "Expected from user %v, got %v", r.ToUser, body.Event.BroadcasterUserID)

	params = *&events.MockEventParameters{
		ToUserID:           user,
		Description:        "Twitch Follower Goal",
		Transport:          models.TransportEventSub,
		Trigger:            "goal-end",
		SubscriptionStatus: "enabled",
	}

	r, err = Event{}.GenerateEvent(params)
	a.Nil(err)

	err = json.Unmarshal(r.JSON, &body)
	a.Nil(err)

	a.Equal("channel.goal.end", body.Subscription.Type, "Expected event type %v, got %v", "channel.goal.end", body.Subscription.Type)
	a.Equal(user, body.Event.BroadcasterUserID, "Expected from user %v, got %v", r.ToUser, body.Event.BroadcasterUserID)
	a.NotNil(body.Event.EndedAt, "Expected endedDate to be nil got %v", body.Event.EndedAt)
	a.NotNil(body.Event.IsAchieved, "Expected endedDate to be nil got %v", body.Event.IsAchieved)
}

func TestFakeTransport(t *testing.T) {
	a := test_setup.SetupTestEnv(t)

	params := *&events.MockEventParameters{
		FromUserID: user,
		Transport:  "fake_transport",
		Trigger:    "unsubscribe",
	}

	r, err := Event{}.GenerateEvent(params)
	a.Nil(err)
	a.Empty(r)
}

func TestValidTrigger(t *testing.T) {
	a := test_setup.SetupTestEnv(t)

	r := Event{}.ValidTrigger("goal-begin")
	a.Equal(true, r)

	r = Event{}.ValidTrigger("goal-started")
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

	r := Event{}.GetTopic(models.TransportEventSub, "goal-progress")
	a.NotNil(r)
}
