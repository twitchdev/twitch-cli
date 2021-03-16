// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package hype_train_progress

import (
	"encoding/json"
	"testing"

	"github.com/twitchdev/twitch-cli/internal/events"
	"github.com/twitchdev/twitch-cli/internal/models"
	"github.com/twitchdev/twitch-cli/internal/util"
)

var toUser = "4567"

func TestEventSub(t *testing.T) {
	a := util.SetupTestEnv(t)

	params := *&events.MockEventParameters{
		FromUserID: fromUser,
		ToUserID:   toUser,
		Transport:  models.TransportEventSub,
		Trigger:    "hype-train-progress",
	}

	r, err := Event{}.GenerateEvent(params)
	a.Nil(err)

	var body models.HypeTrainEventProgressSubResponse
	err = json.Unmarshal(r.JSON, &body)
	a.Nil(err)

	a.Equal("channel.hype_train.progress", body.Subscription.Type, "Expected event type %v, got %v", "channel.hype_train.progress", body.Subscription.Type)
	a.Equal(toUser, body.Event.BroadcasterUserID, "Expected to user %v, got %v", toUser, body.Event.BroadcasterUserID)
	
}

func TestWebSub(t *testing.T) {
	a := util.SetupTestEnv(t)

	params := *&events.MockEventParameters{
		ToUserID:   toUser,
		Transport:  models.TransportWebSub,
		Trigger:    "hype-train-progress",
	}

	r, err := Event{}.GenerateEvent(params)
	a.Nil(err)

	var body models.HypeTrainWebSubResponse
	err = json.Unmarshal(r.JSON, &body)
	a.Nil(err)

	a.Equal("hypetrain.progression", body.Data[0].EventType, "Expected event type %v, got %v", "hypetrain.progression", body.Data[0].EventType)
	a.Equal(toUser, body.Data[0].EventData.BroadcasterID, "Expected to user %v, got %v", toUser, body.Data[0].EventData.BroadcasterID)
	

	params = *&events.MockEventParameters{
		ToUserID:   toUser,
		Transport:  models.TransportWebSub,
	}
}
func TestFakeTransport(t *testing.T) {
	a := util.SetupTestEnv(t)

	params := *&events.MockEventParameters{
		ToUserID:   toUser,
		Transport:  "fake_transport",
		Trigger:    "hype-train-progress",
	}

	r, err := Event{}.GenerateEvent(params)
	a.Nil(err)
	a.Empty(r)
}
func TestValidTrigger(t *testing.T) {
	a := util.SetupTestEnv(t)

	r := Event{}.ValidTrigger("hype-train-progress")
	a.Equal(true, r)

}

func TestValidTransport(t *testing.T) {
	a := util.SetupTestEnv(t)

	r := Event{}.ValidTransport(models.TransportWebSub)
	a.Equal(true, r)

	r = Event{}.ValidTransport(models.TransportEventSub)
	a.Equal(true, r)
}

func TestGetTopic(t *testing.T) {
	a := util.SetupTestEnv(t)

	r := Event{}.GetTopic(models.TransportWebSub, "hype-train-progress")
	a.Equal("hypetrain.progression", r, "Expected %v, got %v", "hypetrain.progression", r)
}