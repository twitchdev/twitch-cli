// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package trigger

import (
	"encoding/json"
	"testing"

	"github.com/twitchdev/twitch-cli/internal/models"
)

func TestEventsubReward(t *testing.T) {
	params := *&RewardParams{
		Transport: "eventsub",
		Type:      "channel.channel_points_custom_reward.add",
		ToUser:    toUser,
		Title:     "Test Title",
		Prompt:    "Test Prompt",
		Cost:      1337,
	}

	r, err := GenerateRewardBody(params)
	if err != nil {
		t.Error(err)
	}
	var body models.RewardEventSubResponse

	if err = json.Unmarshal(r.JSON, &body); err != nil {
		t.Error("Error unmarshalling JSON")
	}

	if body.Event.BroadcasterUserID != toUser {
		t.Errorf("Expected to user %v, got %v", toUser, body.Event.BroadcasterUserID)
	}

	if body.Event.Cost != 1337 {
		t.Errorf("Expected reward cost 1337, got %v", body.Event.Cost)
	}
}
