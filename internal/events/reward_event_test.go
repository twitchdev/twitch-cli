// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package events

import (
	"encoding/json"
	"testing"

	"github.com/twitchdev/twitch-cli/internal/models"
	"github.com/twitchdev/twitch-cli/internal/util"
)

func TestEventsubReward(t *testing.T) {
	a := util.SetupTestEnv(t)

	params := *&RewardParams{
		Transport: TransportEventSub,
		Type:      "channel.channel_points_custom_reward.add",
		ToUser:    toUser,
		Title:     "Test Title",
		Prompt:    "Test Prompt",
		Cost:      1337,
	}

	r, err := GenerateRewardBody(params)
	a.Nil(err)

	var body models.RewardEventSubResponse
	err = json.Unmarshal(r.JSON, &body)
	a.Nil(err)

	a.Equal(toUser, body.Event.BroadcasterUserID, "Expected to user %v, got %v", toUser, body.Event.BroadcasterUserID)
	a.Equal(params.Cost, body.Event.Cost, "Expected cost %v, got %v", params.Cost, body.Event.Cost)

	params = *&RewardParams{
		Transport: TransportEventSub,
	}

	r, err = GenerateRewardBody(params)
	a.Nil(err)

	err = json.Unmarshal(r.JSON, &body)
	a.Nil(err)

	a.NotNil(body.Event.BroadcasterUserID)
	a.NotNil(body.Event.Cost)
	a.NotNil(body.Event.ID)
}

func TestWebsubReward(t *testing.T) {
	a := util.SetupTestEnv(t)

	params := *&RewardParams{
		Transport: TransportWebSub,
	}

	_, err := GenerateRewardBody(params)
	a.NotNil(err)

}
