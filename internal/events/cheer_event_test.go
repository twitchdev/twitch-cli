// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package trigger

import (
	"encoding/json"
	"testing"

	"github.com/twitchdev/twitch-cli/internal/models"
)

func TestEventsubCheer(t *testing.T) {
	params := *&CheerParams{
		FromUser:  fromUser,
		ToUser:    toUser,
		Transport: "eventsub",
	}

	r, err := GenerateCheerBody(params)
	if err != nil {
		t.Error(err)
	}
	var body models.CheerEventSubResponse

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

func TestEventsubAnonymousCheer(t *testing.T) {
	params := *&CheerParams{
		FromUser:    fromUser,
		ToUser:      toUser,
		Transport:   "eventsub",
		IsAnonymous: true,
	}

	r, err := GenerateCheerBody(params)
	if err != nil {
		t.Error(err)
	}
	var body models.CheerEventSubResponse

	if err = json.Unmarshal(r.JSON, &body); err != nil {
		t.Error("Error unmarshalling JSON")
	}

	if body.Event.BroadcasterUserID != toUser {
		t.Errorf("Expected to user %v, got %v", toUser, body.Event.BroadcasterUserID)
	}

	if body.Event.UserID != "" {
		t.Errorf("Expected from user %v, got %v", r.ToUser, body.Event.UserID)
	}
}
