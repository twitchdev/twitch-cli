// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package ban

import (
	"encoding/json"
	"testing"

	"github.com/twitchdev/twitch-cli/internal/events"
	"github.com/twitchdev/twitch-cli/internal/models"
	"github.com/twitchdev/twitch-cli/internal/util"
)

var fromUser = "1234"
var toUser = "4567"

func TestEventSubBan(t *testing.T) {
	a := util.SetupTestEnv(t)
	params := events.MockEventParameters{
		FromUserID: fromUser,
		ToUserID:   toUser,
		Transport:  models.TransportEventSub,
		Trigger:    "ban",
	}

	r, err := Event{}.GenerateEvent(params)
	a.Nil(err, "Error generating body.")

	var body models.BanEventSubResponse

	err = json.Unmarshal(r.JSON, &body)
	a.Nil(err, "Error unmarshalling JSON")

	a.Equal(toUser, body.Event.BroadcasterUserID, "Expected to user %v, got %v", toUser, body.Event.BroadcasterUserID)
	a.Equal(fromUser, body.Event.UserID, "Expected from user %v, got %v", r.ToUser, body.Event.UserID)

	// test for unban
	params = events.MockEventParameters{
		FromUserID:  fromUser,
		ToUserID:    toUser,
		Transport:   models.TransportEventSub,
		Trigger:     "unban",
	}

	r, err = Event{}.GenerateEvent(params)
	a.Nil(err)

	err = json.Unmarshal(r.JSON, &body)
	a.Nil(err)

	a.Equal(toUser, body.Event.BroadcasterUserID, "Expected to user %v, got %v", toUser, body.Event.BroadcasterUserID)
	a.Equal(fromUser, body.Event.UserID, "Expected  from user %v, got %v", fromUser, body.Event.UserID)
}

func TestWebSubBan(t *testing.T) {
	a := util.SetupTestEnv(t)

	params := *&events.MockEventParameters{
		FromUserID: fromUser,
		ToUserID:   toUser,
		Transport:  models.TransportWebSub,
		Trigger:    "ban",
	}

	r, err := Event{}.GenerateEvent(params)
	a.Nil(err)

	var body models.BanWebSubResponse
	err = json.Unmarshal(r.JSON, &body)
	a.Nil(err)

	a.Equal(toUser, body.Data[0].EventData.BroadcasterID, "Expected to user %v, got %v", toUser, body.Data[0].EventData.BroadcasterID)
	a.Equal(fromUser, body.Data[0].EventData.UserID, "Expected from user %v, got %v", fromUser, body.Data[0].EventData.UserID)


	params = *&events.MockEventParameters{
		FromUserID:  fromUser,
		ToUserID:    toUser,
		Transport:   models.TransportWebSub,
		Trigger:     "unban",
	}

	r, err = Event{}.GenerateEvent(params)
	a.Nil(err)

	err = json.Unmarshal(r.JSON, &body)
	a.Nil(err)

	a.Equal(toUser, body.Data[0].EventData.BroadcasterID, "Expected to user %v, got %v", toUser, body.Data[0].EventData.BroadcasterID)
	a.Equal(fromUser, body.Data[0].EventData.UserID, "Expected from user %v, got %v", fromUser, body.Data[0].EventData.UserID)

}

func TestFakeTransport(t *testing.T) {
	a := util.SetupTestEnv(t)

	params := events.MockEventParameters{
		FromUserID: fromUser,
		ToUserID:   toUser,
		Transport:  "fake_transport",
		Trigger:    "unban",
	}

	r, err := Event{}.GenerateEvent(params)
	a.Nil(err)
	a.Empty(r)
}

func TestValidTrigger(t *testing.T) {
	a := util.SetupTestEnv(t)

	r := Event{}.ValidTrigger("ban")
	a.Equal(true, r)

	r = Event{}.ValidTrigger("unban")
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

	r := Event{}.GetTopic(models.TransportEventSub, "ban")
	a.NotNil(r)
}
