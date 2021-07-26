// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package drop

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
		Trigger:    "drop",
	}

	r, err := Event{}.GenerateEvent(params)
	a.Nil(err)

	var body models.DropsEntitlementEventSubResponse // replace with actual value
	err = json.Unmarshal(r.JSON, &body)
	a.Nil(err)

	a.Len(body.Events, 1)
}

func TestFakeTransport(t *testing.T) {
	a := test_setup.SetupTestEnv(t)

	params := *&events.MockEventParameters{
		FromUserID: fromUser,
		ToUserID:   toUser,
		Transport:  "fake_transport",
		Trigger:    "drop",
	}

	r, err := Event{}.GenerateEvent(params)
	a.Nil(err)
	a.Empty(r)
}
func TestValidTrigger(t *testing.T) {
	a := test_setup.SetupTestEnv(t)

	r := Event{}.ValidTrigger("drop")
	a.Equal(true, r)

	r = Event{}.ValidTrigger("notdrop")
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

	r := Event{}.GetTopic(models.TransportEventSub, "drop")
	a.NotNil(r)
}
