// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package poll

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
		Trigger:    "poll-begin",
	}

	r, err := Event{}.GenerateEvent(params)
	a.Nil(err)

	var body models.PollEventSubResponse
	err = json.Unmarshal(r.JSON, &body)
	a.Nil(err)
	a.NotEmpty(body.Event.EndsAt)
	a.Empty(body.Event.EndedAt)

	params = *&events.MockEventParameters{
		FromUserID: fromUser,
		ToUserID:   toUser,
		Transport:  models.TransportEventSub,
		Trigger:    "poll-progress",
	}

	r, err = Event{}.GenerateEvent(params)
	a.Nil(err)

	body = models.PollEventSubResponse{}
	err = json.Unmarshal(r.JSON, &body)
	a.Nil(err)
	a.NotEmpty(body.Event.EndsAt)
	a.Empty(body.Event.EndedAt)

	params = *&events.MockEventParameters{
		FromUserID: fromUser,
		ToUserID:   toUser,
		Transport:  models.TransportEventSub,
		Trigger:    "poll-end",
	}

	r, err = Event{}.GenerateEvent(params)
	a.Nil(err)

	body = models.PollEventSubResponse{}
	err = json.Unmarshal(r.JSON, &body)
	a.Nil(err)
	a.Empty(body.Event.EndsAt)
	a.NotEmpty(body.Event.EndedAt)
}

func TestFakeTransport(t *testing.T) {
	a := test_setup.SetupTestEnv(t)

	params := *&events.MockEventParameters{
		FromUserID: fromUser,
		ToUserID:   toUser,
		Transport:  "fake_transport",
		Trigger:    "poll-begin",
	}

	r, err := Event{}.GenerateEvent(params)
	a.Nil(err)
	a.Empty(r)
}
func TestValidTrigger(t *testing.T) {
	a := test_setup.SetupTestEnv(t)

	r := Event{}.ValidTrigger("poll-begin")
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

	r := Event{}.GetTopic(models.TransportEventSub, "poll-begin")
	a.NotNil(r)
}
