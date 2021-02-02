// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package events

import (
	"encoding/json"
	"testing"

	"github.com/twitchdev/twitch-cli/internal/models"
	"github.com/twitchdev/twitch-cli/internal/util"
)

func TestEventsubCheer(t *testing.T) {
	a := util.SetupTestEnv(t)
	params := *&CheerParams{
		FromUser:  fromUser,
		ToUser:    toUser,
		Transport: TransportEventSub,
	}

	r, err := GenerateCheerBody(params)
	a.Nil(err, "Error generating body.")

	var body models.CheerEventSubResponse

	err = json.Unmarshal(r.JSON, &body)
	a.Nil(err, "Error unmarshalling JSON")

	a.Equal(toUser, body.Event.BroadcasterUserID, "Expected to user %v, got %v", toUser, body.Event.BroadcasterUserID)
	a.Equal(fromUser, body.Event.UserID, "Expected from user %v, got %v", r.ToUser, body.Event.UserID)

	params = *&CheerParams{
		Transport: TransportEventSub,
	}

	r, err = GenerateCheerBody(params)
	a.Nil(err, "Error generating body.")

	err = json.Unmarshal(r.JSON, &body)
	a.Nil(err, "Error unmarshalling JSON")

	a.NotNil(body.Event.BroadcasterUserID, "BroadcasterUserID empty")
	a.NotNil(body.Event.UserID, "UserID empty")
}

func TestEventsubAnonymousCheer(t *testing.T) {
	a := util.SetupTestEnv(t)

	params := *&CheerParams{
		FromUser:    fromUser,
		ToUser:      toUser,
		Transport:   TransportEventSub,
		IsAnonymous: true,
	}

	r, err := GenerateCheerBody(params)
	a.Nil(err)

	var body models.CheerEventSubResponse
	err = json.Unmarshal(r.JSON, &body)
	a.Nil(err)

	a.Equal(toUser, body.Event.BroadcasterUserID, "Expected to user %v, got %v", toUser, body.Event.BroadcasterUserID)
	a.Equal("", body.Event.UserID, "Expected empty from user, got %v", body.Event.UserID)
}

func TestWebsubCheer(t *testing.T) {
	a := util.SetupTestEnv(t)

	params := *&CheerParams{
		Transport: TransportWebSub,
	}

	_, err := GenerateCheerBody(params)
	a.NotNil(err, "Expected error (Cheer unsupported on websub")
}
