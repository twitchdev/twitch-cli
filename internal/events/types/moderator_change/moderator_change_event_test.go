// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package moderator_change

import (
	"encoding/json"
	"testing"

	"github.com/twitchdev/twitch-cli/internal/events"
	"github.com/twitchdev/twitch-cli/internal/models"
	"github.com/twitchdev/twitch-cli/internal/util"
)

var fromUser = "1234"
var toUser = "4567"

func TestWebSub(t *testing.T) {
	a := util.SetupTestEnv(t)

	params := *&events.MockEventParameters{
		FromUserID: fromUser,
		ToUserID:   toUser,
		Transport:  models.TransportWebSub,
		Trigger:    "add-moderator",
	}

	r, err := Event{}.GenerateEvent(params)
	a.Nil(err)

	var body models.ModeratorChangeWebSubResponse
	err = json.Unmarshal(r.JSON, &body)
	a.Nil(err)

	a.Equal("moderation.moderator.add", body.Data[0].EventType, "Expected event type %v, got %v", "moderation.moderator.add", body.Data[0].EventType)
	a.Equal(toUser, body.Data[0].EventData.BroadcasterID, "Expected to user %v, got %v", toUser, body.Data[0].EventData.BroadcasterID)
	a.Equal(fromUser, body.Data[0].EventData.UserID, "Expected from user %v, got %v", fromUser, body.Data[0].EventData.UserID)

	params = *&events.MockEventParameters{
		FromUserID: fromUser,
		ToUserID:   toUser,
		Transport:  models.TransportWebSub,
	}
}
func TestFakeTransport(t *testing.T) {
	a := util.SetupTestEnv(t)

	params := *&events.MockEventParameters{
		FromUserID: fromUser,
		ToUserID:   toUser,
		Transport:  "fake_transport",
		Trigger:    "add-moderator",
	}

	r, err := Event{}.GenerateEvent(params)
	a.Nil(err)
	a.Empty(r)
}
func TestValidTrigger(t *testing.T) {
	a := util.SetupTestEnv(t)

	r := Event{}.ValidTrigger("add-moderator")
	a.Equal(true, r)

	r = Event{}.ValidTrigger("remove-moderator")
	a.Equal(true, r)

	r = Event{}.ValidTrigger("update-moderator")
	a.Equal(false, r)
}

func TestValidTransport(t *testing.T) {
	a := util.SetupTestEnv(t)

	r := Event{}.ValidTransport(models.TransportWebSub)
	a.Equal(true, r)

	r = Event{}.ValidTransport(models.TransportEventSub)
	a.Equal(false, r)
}

func TestGetTopic(t *testing.T) {
	a := util.SetupTestEnv(t)

	r := Event{}.GetTopic(models.TransportWebSub, "add-moderator")
	a.Equal("moderation.moderator.add", r, "Expected %v, got %v", "moderation.moderator.add", r)

	r = Event{}.GetTopic(models.TransportWebSub, "remove-moderator")
	a.Equal("moderation.moderator.remove", r, "Expected %v, got %v", "moderation.moderator.remove", r)
}
