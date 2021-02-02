// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package events

import (
	"encoding/json"
	"testing"

	"github.com/twitchdev/twitch-cli/internal/models"
	"github.com/twitchdev/twitch-cli/internal/util"
)

func TestEventSubFollow(t *testing.T) {
	a := util.SetupTestEnv(t)

	params := *&FollowParams{
		FromUser:  fromUser,
		ToUser:    toUser,
		Transport: TransportEventSub,
	}

	r, err := GenerateFollowBody(params)
	a.Nil(err)

	var body models.FollowEventSubResponse
	err = json.Unmarshal(r.JSON, &body)
	a.Nil(err)

	a.Equal(toUser, body.Event.BroadcasterUserID, "Expected to user %v, got %v", toUser, body.Event.BroadcasterUserID)
	a.Equal(fromUser, body.Event.UserID, "Expected from user %v, got %v", r.ToUser, body.Event.UserID)

	params = *&FollowParams{
		Transport: TransportEventSub,
	}

	r, err = GenerateFollowBody(params)

	err = json.Unmarshal(r.JSON, &body)
	a.Nil(err, "Error unmarshalling JSON")

	a.NotNil(body.Event.BroadcasterUserID, "BroadcasterUserID empty")
	a.NotNil(body.Event.UserID, "UserID empty")
}

func TestWebsubFollow(t *testing.T) {
	a := util.SetupTestEnv(t)

	params := *&FollowParams{
		FromUser:  fromUser,
		ToUser:    toUser,
		Transport: TransportWebSub,
	}

	r, err := GenerateFollowBody(params)
	a.Nil(err)

	var body models.FollowWebSubResponse
	err = json.Unmarshal(r.JSON, &body)
	a.Nil(err)

	a.Equal(toUser, body.Data[0].ToID, "Expected to user %v, got %v", toUser, body.Data[0].ToID)
	a.Equal(fromUser, body.Data[0].FromID, "Expected from user %v, got %v", r.ToUser, body.Data[0].FromID)
}
