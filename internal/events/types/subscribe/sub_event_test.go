// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package subscribe

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

	params := *&events.MockEventParameters{
		FromUserID: fromUser,
		ToUserID:   toUser,
		Transport:  models.TransportEventSub,
		Trigger:    "subscribe",
	}

	r, err := Event{}.GenerateEvent(params)
	a.Nil(err)

	var body models.SubEventSubResponse
	err = json.Unmarshal(r.JSON, &body)
	a.Nil(err)

	a.Equal(toUser, body.Event.BroadcasterUserID, "Expected to user %v, got %v", toUser, body.Event.BroadcasterUserID)
	a.Equal(fromUser, body.Event.UserID, "Expected from user %v, got %v", r.ToUser, body.Event.UserID)
}

func TestWebSub(t *testing.T) {
	a := test_setup.SetupTestEnv(t)

	params := *&events.MockEventParameters{
		FromUserID: fromUser,
		ToUserID:   toUser,
		Transport:  models.TransportWebSub,
		Trigger:    "unsubscribe",
	}

	r, err := Event{}.GenerateEvent(params)
	a.Nil(err)

	var body models.SubWebSubResponse
	err = json.Unmarshal(r.JSON, &body)
	a.Nil(err)

	a.Equal(toUser, body.Data[0].EventData.BroadcasterID, "Expected to user %v, got %v", toUser, body.Data[0].EventData.BroadcasterID)
	a.Equal(fromUser, body.Data[0].EventData.UserID, "Expected from user %v, got %v", fromUser, body.Data[0].EventData.UserID)

	a.Equal(false, body.Data[0].EventData.IsGift)

	params = *&events.MockEventParameters{
		FromUserID:  fromUser,
		ToUserID:    toUser,
		Transport:   models.TransportWebSub,
		IsGift:      true,
		IsAnonymous: true,
		Trigger:     "gift",
	}

	r, err = Event{}.GenerateEvent(params)
	a.Nil(err)

	err = json.Unmarshal(r.JSON, &body)
	a.Nil(err)

	a.Equal(toUser, body.Data[0].EventData.BroadcasterID, "Expected to user %v, got %v", toUser, body.Data[0].EventData.BroadcasterID)
	a.Equal(fromUser, body.Data[0].EventData.UserID, "Expected from user %v, got %v", fromUser, body.Data[0].EventData.UserID)
	a.Equal("274598607", body.Data[0].EventData.GifterID)

	a.Equal(true, body.Data[0].EventData.IsGift)
}
func TestFakeTransport(t *testing.T) {
	a := test_setup.SetupTestEnv(t)

	params := *&events.MockEventParameters{
		FromUserID: fromUser,
		ToUserID:   toUser,
		Transport:  "fake_transport",
		Trigger:    "unsubscribe",
	}

	r, err := Event{}.GenerateEvent(params)
	a.Nil(err)
	a.Empty(r)
}
func TestValidTrigger(t *testing.T) {
	a := test_setup.SetupTestEnv(t)

	r := Event{}.ValidTrigger("gift")
	a.Equal(true, r)

	r = Event{}.ValidTrigger("notgift")
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

	r := Event{}.GetTopic(models.TransportEventSub, "subscribe")
	a.NotNil(r)
}
