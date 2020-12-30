// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package trigger

import (
	"encoding/json"
	"testing"

	"github.com/twitchdev/twitch-cli/internal/models"
)

func TestEventSubFollow(t *testing.T) {
	params := *&FollowParams{
		FromUser:  fromUser,
		ToUser:    toUser,
		Transport: "eventsub",
	}

	r, err := GenerateFollowBody(params)
	if err != nil {
		t.Error(err)
	}
	var body models.FollowEventSubResponse
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

func TestWebsubFollow(t *testing.T) {
	params := *&FollowParams{
		FromUser:  fromUser,
		ToUser:    toUser,
		Transport: "websub",
	}

	r, err := GenerateFollowBody(params)
	if err != nil {
		t.Error(err)
	}

	var body models.FollowWebSubResponse
	if err = json.Unmarshal(r.JSON, &body); err != nil {
		t.Error("Error unmarshalling JSON")
	}

	if body.Data[0].ToID != toUser {
		t.Errorf("Expected to user %v, got %v", toUser, body.Data[0].ToID)
	}

	if body.Data[0].FromID != fromUser {
		t.Errorf("Expected from user %v, got %v", r.ToUser, body.Data[0].FromID)
	}
}
