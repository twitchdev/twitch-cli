// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package trigger

import (
	"encoding/json"
	"testing"

	"github.com/twitchdev/twitch-cli/internal/models"
)

func TestEventsubRedemption(t *testing.T) {
	params := *&RedemptionParams{
		Transport: "eventsub",
		Type:      "channel.channel_points_custom_reward_redemption.add",
		ToUser:    toUser,
		FromUser:  fromUser,
		Title:     "Test Title",
		Prompt:    "Test Prompt",
		Status:    "tested",
		RewardID:  "12345678-1234-abcd-5678-000000000000",
		Cost:      1337,
	}

	r, err := GenerateRedemptionBody(params)
	if err != nil {
		t.Error(err)
	}
	var body models.RedemptionEventSubResponse

	if err = json.Unmarshal(r.JSON, &body); err != nil {
		t.Error("Error unmarshalling JSON")
	}

	if body.Event.BroadcasterUserID != toUser {
		t.Errorf("Expected to user %v, got %v", toUser, body.Event.BroadcasterUserID)
	}

	if body.Event.UserID != fromUser {
		t.Errorf("Expected from user %v, got %v", r.ToUser, body.Event.UserID)
	}

	if body.Event.Status != "tested" {
		t.Errorf("Expected status tested, got %v", body.Event.Status)
	}

	if body.Event.Reward.Cost != 1337 {
		t.Errorf("Expected reward cost 1337, got %v", body.Event.Reward.Cost)
	}

	if body.Event.Reward.ID != "12345678-1234-abcd-5678-000000000000" {
		t.Errorf("Expected reward cost 12345678-1234-abcd-5678-00000000000, got %v", body.Event.Reward.ID)
	}
}
