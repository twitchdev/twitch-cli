// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package authorization

import (
	"encoding/json"
	"testing"

	"github.com/spf13/viper"
	"github.com/twitchdev/twitch-cli/internal/events"
	"github.com/twitchdev/twitch-cli/internal/models"
	"github.com/twitchdev/twitch-cli/test_setup"
)

var fromUser = "1234"
var toUser = "4567"

func TestEventSub(t *testing.T) {
	a := test_setup.SetupTestEnv(t)

	viper.Set("ClientID", "1234")
	params := *&events.MockEventParameters{
		FromUserID:         fromUser,
		ToUserID:           toUser,
		Transport:          models.TransportEventSub,
		Trigger:            "subscribe",
		SubscriptionStatus: "enabled",
	}

	r, err := Event{}.GenerateEvent(params)
	a.Nil(err)

	var body models.AuthorizationRevokeEventSubResponse
	err = json.Unmarshal(r.JSON, &body)
	a.Nil(err)

	a.NotEmpty(body.Event.ClientID)
	a.Equal(body.Event.ClientID, body.Subscription.Condition.ClientID)
	a.Equal("1234", body.Event.ClientID)
}
func TestFakeTransport(t *testing.T) {
	a := test_setup.SetupTestEnv(t)

	params := *&events.MockEventParameters{
		FromUserID:         fromUser,
		ToUserID:           toUser,
		Transport:          "fake_transport",
		Trigger:            "unsubscribe",
		SubscriptionStatus: "enabled",
	}

	r, err := Event{}.GenerateEvent(params)
	a.Nil(err)
	a.Empty(r)
}
func TestValidTrigger(t *testing.T) {
	a := test_setup.SetupTestEnv(t)

	r := Event{}.ValidTrigger("revoke")
	a.Equal(true, r)

	r = Event{}.ValidTrigger("fake_revoke")
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

	r := Event{}.GetTopic(models.TransportEventSub, "revoke")
	a.NotNil(r)
}
