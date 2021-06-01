// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package channel_points_reward

import (
	"encoding/json"
	"testing"

	"github.com/twitchdev/twitch-cli/internal/events"
	"github.com/twitchdev/twitch-cli/internal/models"
	"github.com/twitchdev/twitch-cli/internal/util"
)

var fromUser = "1234"
var toUser = "4567"

func TestEventSub(t *testing.T) {
	a := util.SetupTestEnv(t)

	params := events.MockEventParameters{
		Transport:  models.TransportEventSub,
		Trigger:    "add-redemption",
		ToUserID:   toUser,
		FromUserID: fromUser,
		Status:     "tested",
		ItemID:     "12345678-1234-abcd-5678-000000000000",
		Cost:       1337,
		ItemName:   "Testing",
	}

	r, err := Event{}.GenerateEvent(params)
	a.Nil(err)

	var body models.RewardEventSubResponse
	err = json.Unmarshal(r.JSON, &body)
	a.Nil(err)

	a.Equal(toUser, body.Event.BroadcasterUserID, "Expected to user %v, got %v", toUser, body.Event.BroadcasterUserID)
	a.Equal(params.Cost, body.Event.Cost, "Expected cost %v, got %v", params.Cost, body.Event.Cost)
	a.Equal(params.ItemName, body.Event.Title)
}

func TestWebSub(t *testing.T) {
	a := util.SetupTestEnv(t)

	params := events.MockEventParameters{
		Transport: models.TransportWebSub,
	}

	_, err := Event{}.GenerateEvent(params)
	a.NotNil(err)

}

func TestFakeTransport(t *testing.T) {
	a := util.SetupTestEnv(t)

	params := events.MockEventParameters{
		FromUserID: fromUser,
		ToUserID:   toUser,
		Transport:  "fake_transport",
		Trigger:    "add-reward",
	}

	r, err := Event{}.GenerateEvent(params)
	a.Nil(err)
	a.Empty(r)
}

func TestValidTrigger(t *testing.T) {
	a := util.SetupTestEnv(t)

	r := Event{}.ValidTrigger("add-reward")
	a.Equal(true, r)

	r = Event{}.ValidTrigger("not_trigger_keyword")
	a.Equal(false, r)
}

func TestValidTransport(t *testing.T) {
	a := util.SetupTestEnv(t)

	r := Event{}.ValidTransport(models.TransportEventSub)
	a.Equal(true, r)

	r = Event{}.ValidTransport("noteventsub")
	a.Equal(false, r)
}
func TestGetTopic(t *testing.T) {
	a := util.SetupTestEnv(t)

	r := Event{}.GetTopic(models.TransportEventSub, "add-reward")
	a.NotNil(r)
}
