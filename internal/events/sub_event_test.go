// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package events

import (
	"encoding/json"
	"testing"

	"github.com/twitchdev/twitch-cli/internal/models"
	"github.com/twitchdev/twitch-cli/internal/util"
)

func TestEventsubSubscribe(t *testing.T) {
	a := util.SetupTestEnv(t)

	params := *&SubscribeParams{
		FromUser:  fromUser,
		ToUser:    toUser,
		Transport: TransportEventSub,
	}

	r, err := GenerateSubBody(params)
	a.Nil(err)

	var body models.SubEventSubResponse
	err = json.Unmarshal(r.JSON, &body)
	a.Nil(err)

	a.Equal(toUser, body.Event.BroadcasterUserID, "Expected to user %v, got %v", toUser, body.Event.BroadcasterUserID)
	a.Equal(fromUser, body.Event.UserID, "Expected from user %v, got %v", r.ToUser, body.Event.UserID)
}

func TestWebusbSubscribe(t *testing.T) {
	a := util.SetupTestEnv(t)

	params := *&SubscribeParams{
		FromUser:  fromUser,
		ToUser:    toUser,
		Transport: TransportWebSub,
	}

	r, err := GenerateSubBody(params)
	a.Nil(err)

	var body models.SubWebSubResponse
	err = json.Unmarshal(r.JSON, &body)
	a.Nil(err)

	a.Equal(toUser, body.Data[0].EventData.BroadcasterID, "Expected to user %v, got %v", toUser, body.Data[0].EventData.BroadcasterID)
	a.Equal(fromUser, body.Data[0].EventData.UserID, "Expected from user %v, got %v", fromUser, body.Data[0].EventData.UserID)

	a.Equal(false, body.Data[0].EventData.IsGift)
}

func TestWebsubGifts(t *testing.T) {
	a := util.SetupTestEnv(t)

	params := *&SubscribeParams{
		ToUser:          toUser,
		FromUser:        fromUser,
		Transport:       TransportWebSub,
		IsGift:          true,
		IsAnonymousGift: true,
	}

	r, err := GenerateSubBody(params)
	a.Nil(err)

	var body models.SubWebSubResponse
	err = json.Unmarshal(r.JSON, &body)
	a.Nil(err)

	a.Equal(toUser, body.Data[0].EventData.BroadcasterID, "Expected to user %v, got %v", toUser, body.Data[0].EventData.BroadcasterID)
	a.Equal(fromUser, body.Data[0].EventData.UserID, "Expected from user %v, got %v", fromUser, body.Data[0].EventData.UserID)
	a.Equal("274598607", body.Data[0].EventData.GifterID)

	a.Equal(true, body.Data[0].EventData.IsGift)
}
