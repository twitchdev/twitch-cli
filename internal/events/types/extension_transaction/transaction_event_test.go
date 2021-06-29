// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package extension_transaction

import (
	"encoding/json"
	"testing"

	"github.com/twitchdev/twitch-cli/internal/events"
	"github.com/twitchdev/twitch-cli/internal/models"
	"github.com/twitchdev/twitch-cli/test_setup"
)

var fromUser = "1234"
var toUser = "4567"

func TestWebusbTransaction(t *testing.T) {
	a := test_setup.SetupTestEnv(t)

	params := events.MockEventParameters{
		FromUserID: fromUser,
		ToUserID:   toUser,
		Transport:  models.TransportWebSub,
		Trigger:    "transaction",
	}

	r, err := Event{}.GenerateEvent(params)
	a.Nil(err)

	var body models.TransactionWebSubResponse
	err = json.Unmarshal(r.JSON, &body)
	a.Nil(err)

	a.Equal(toUser, body.Data[0].BroadcasterID, "Expected to user %v, got %v", toUser, body.Data[0].BroadcasterID)
	a.Equal(fromUser, body.Data[0].UserID, "Expected from user %v, got %v", r.ToUser, body.Data[0].UserID)

}

func TestEventsubTransaction(t *testing.T) {
	a := test_setup.SetupTestEnv(t)

	params := events.MockEventParameters{
		FromUserID: fromUser,
		ToUserID:   toUser,
		Transport:  models.TransportEventSub,
		Trigger:    "transaction",
	}

	_, err := Event{}.GenerateEvent(params)
	a.NotNil(err)
}

func TestFakeTransport(t *testing.T) {
	a := test_setup.SetupTestEnv(t)

	params := events.MockEventParameters{
		FromUserID: fromUser,
		ToUserID:   toUser,
		Transport:  "fake_transport",
		Trigger:    "transaction",
	}

	r, err := Event{}.GenerateEvent(params)
	a.Nil(err)
	a.Empty(r)
}

func TestValidTrigger(t *testing.T) {
	a := test_setup.SetupTestEnv(t)

	r := Event{}.ValidTrigger("transaction")
	a.Equal(true, r)

	r = Event{}.ValidTrigger("not_trigger_keyword")
	a.Equal(false, r)
}

func TestValidTransport(t *testing.T) {
	a := test_setup.SetupTestEnv(t)

	r := Event{}.ValidTransport(models.TransportWebSub)
	a.Equal(true, r)

	r = Event{}.ValidTransport("noteventsub")
	a.Equal(false, r)
}
func TestGetTopic(t *testing.T) {
	a := test_setup.SetupTestEnv(t)

	r := Event{}.GetTopic(models.TransportWebSub, "transaction")
	a.NotNil(r)
}
