// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package stream_change

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

	params := *&events.MockEventParameters{
		FromUserID: fromUser,
		ToUserID:   toUser,
		Transport:  models.TransportEventSub,
		Trigger:    "stream-change",
	}

	r, err := Event{}.GenerateEvent(params)
	a.Nil(err)

	var body models.ChannelUpdateEventSubResponse
	err = json.Unmarshal(r.JSON, &body)
	a.Nil(err, "Error unmarshalling JSON")

	// write actual tests here (making sure you set appropriate values and the like) for eventsub
	a.Equal(toUser, body.Event.BroadcasterUserID, "Expected Stream Channel %v, got %v", toUser, body.Event.BroadcasterUserID)

	// test for changing a title
	params = events.MockEventParameters{
		FromUserID: fromUser,
		ToUserID:   toUser,
		Transport:  models.TransportEventSub,
		Trigger:    "stream_change",
	}

	r, err = Event{}.GenerateEvent(params)
	a.Nil(err)

	err = json.Unmarshal(r.JSON, &body)
	a.Nil(err)

	a.Equal(toUser, body.Event.BroadcasterUserID, "Expected Stream Channel %v, got %v", toUser, body.Event.BroadcasterUserID)
	a.Equal("Example title from the CLI!", body.Event.StreamTitle, "Expected new stream title, got %v", body.Event.StreamTitle)
}

func TestWebSubStreamChange(t *testing.T) {
	a := util.SetupTestEnv(t)

	newStreamTitle := "Awesome new title from the CLI!"

	params := *&events.MockEventParameters{
		FromUserID:  fromUser,
		ToUserID:    toUser,
		Transport:   models.TransportWebSub,
		Trigger:     "stream-change",
		StreamTitle: newStreamTitle,
	}

	r, err := Event{}.GenerateEvent(params)
	a.Nil(err)

	var body models.StreamChangeWebSubResponse
	err = json.Unmarshal(r.JSON, &body)
	a.Nil(err)

	// write tests here for websub
	a.Equal(toUser, body.Data[0].BroadcasterUserID, "Expected Stream Channel %v, got %v", toUser, body.Data[0].BroadcasterUserID)
	a.Equal(newStreamTitle, body.Data[0].StreamTitle, "Expected new stream title, got %v", body.Data[0].StreamTitle)
}
func TestFakeTransport(t *testing.T) {
	a := util.SetupTestEnv(t)

	params := *&events.MockEventParameters{
		FromUserID: fromUser,
		ToUserID:   toUser,
		Transport:  "fake_transport",
		Trigger:    "stream-change",
	}

	r, err := Event{}.GenerateEvent(params)
	a.Nil(err)
	a.Empty(r)
}
func TestValidTrigger(t *testing.T) {
	a := util.SetupTestEnv(t)

	r := Event{}.ValidTrigger("stream-change")
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

	r := Event{}.GetTopic(models.TransportEventSub, "stream-change")
	a.NotNil(r)
}
