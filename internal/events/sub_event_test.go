// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package trigger

import (
	"encoding/json"
	"testing"

	"github.com/twitchdev/twitch-cli/internal/models"
)

func TestEventsubSubscribe(t *testing.T) {
	params := *&SubscribeParams{
		FromUser:  fromUser,
		ToUser:    toUser,
		Transport: "eventsub",
	}

	r, err := GenerateSubBody(params)
	if err != nil {
		t.Error(err)
	}

	var body models.SubEventSubResponse
	if err = json.Unmarshal(r.JSON, &body); err != nil {
		t.Error("Error unmarshalling JSON")
	}

	if body.Event.BroadcasterUserID != toUser {
		t.Errorf("Expected to user %v, got %v", toUser, body.Event.BroadcasterUserID)
	}

	if body.Event.UserID != fromUser {
		t.Errorf("Expected from user %v, got %v", r.ToUser, body.Event.UserID)
	}

}

func TestWebusbSubscribe(t *testing.T) {
	params := *&SubscribeParams{
		FromUser:  fromUser,
		ToUser:    toUser,
		Transport: "websub",
	}

	r, err := GenerateSubBody(params)
	if err != nil {
		t.Error(err)
	}

	var body models.SubWebSubResponse
	if err = json.Unmarshal(r.JSON, &body); err != nil {
		t.Error("Error unmarshalling JSON")
	}

	if body.Data[0].EventData.BroadcasterID != toUser {
		t.Errorf("Expected to user %v, got %v", toUser, r.ToUser)
	}

	if body.Data[0].EventData.UserID != fromUser {
		t.Errorf("Expected from user %v, expected empty string", r.FromUser)
	}

	if body.Data[0].EventData.IsGift {
		t.Error("Marked as git sub when not a gift sub")
	}
}

func TestWebsubGifts(t *testing.T) {
	params := *&SubscribeParams{
		FromUser:  fromUser,
		ToUser:    toUser,
		Transport: "websub",
		IsGift:    true,
	}

	r, err := GenerateSubBody(params)
	if err != nil {
		t.Error(err)
	}

	var body models.SubWebSubResponse
	if err = json.Unmarshal(r.JSON, &body); err != nil {
		t.Error("Error marshalling JSON")
	}

	if body.Data[0].EventData.BroadcasterID != toUser {
		t.Errorf("Expected to user %v, got %v", toUser, r.ToUser)
	}

	if body.Data[0].EventData.UserID != fromUser {
		t.Errorf("Expected from user %v, expected empty string", r.FromUser)
	}

	if !body.Data[0].EventData.IsGift {
		t.Error("Failed to mark as git sub")
	}

}
