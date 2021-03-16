// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package hype_train_begin

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
		ToUserID:   toUser,
		Transport:  models.TransportEventSub,
		Trigger:    "hype-train-begin",
	}

	r, err := Event{}.GenerateEvent(params)
	a.Nil(err)

	var body models.HypeTrainEventBeginSubResponse
	err = json.Unmarshal(r.JSON, &body)
	a.Nil(err)

	a.Equal("channel.hype_train.begin", body.Subscription.Type, "Expected event type %v, got %v", "channel.hype_train.begin", body.Subscription.Type)
	a.Equal(toUser, body.Event.BroadcasterUserID, "Expected to user %v, got %v", toUser, body.Event.BroadcasterUserID)

}

func TestWebSub(t *testing.T) {

	// Hype_Train_begin does not have a websub part
	// Remove?

	a := util.SetupTestEnv(t)

	params := *&events.MockEventParameters{
		ToUserID:   toUser,
		Transport:  models.TransportWebSub,
		Trigger:    "hype-train-begin",
	}

	_, err := Event{}.GenerateEvent(params)
	a.NotNil(err)

	
}
func TestFakeTransport(t *testing.T) {
	a := util.SetupTestEnv(t)

	params := *&events.MockEventParameters{
		ToUserID:   toUser,
		Transport:  "fake_transport",
		Trigger:    "hype-train-begin",
	}

	r, err := Event{}.GenerateEvent(params)
	a.Nil(err)
	a.Empty(r)
}
func TestValidTrigger(t *testing.T) {
	a := util.SetupTestEnv(t)

	r := Event{}.ValidTrigger("hype-train-begin")
	a.Equal(true, r)
}

func TestValidTransport(t *testing.T) {
	a := util.SetupTestEnv(t)

	r := Event{}.ValidTransport(models.TransportWebSub)
	a.Equal(false, r)

	r = Event{}.ValidTransport(models.TransportEventSub)
	a.Equal(true, r)
}

func TestGetTopic(t *testing.T) {
	a := util.SetupTestEnv(t)
}
