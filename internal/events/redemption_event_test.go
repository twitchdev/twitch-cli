// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package trigger

import (
	"encoding/json"
	"testing"

	"github.com/twitchdev/twitch-cli/internal/models"
	"github.com/twitchdev/twitch-cli/internal/util"
)

func TestEventsubRedemption(t *testing.T) {
	a := util.SetupTestEnv(t)

	params := *&RedemptionParams{
		Transport: TransportEventSub,
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
	a.Nil(err)

	var body models.RedemptionEventSubResponse
	err = json.Unmarshal(r.JSON, &body)
	a.Nil(err)

	a.Equal(toUser, body.Event.BroadcasterUserID, "Expected to user %v, got %v", toUser, body.Event.BroadcasterUserID)
	a.Equal(fromUser, body.Event.UserID, "Expected from user %v, got %v", r.ToUser, body.Event.UserID)
	a.Equal(params.Status, body.Event.Status)
	a.Equal(params.Cost, body.Event.Reward.Cost)
	a.Equal(params.RewardID, body.Event.Reward.ID)

	params = *&RedemptionParams{
		Transport: TransportEventSub,
	}
	r, err = GenerateRedemptionBody(params)
	a.Nil(err)

	err = json.Unmarshal(r.JSON, &body)
	a.Nil(err)

	a.NotNil(body.Event.BroadcasterUserID)
	a.NotNil(body.Event.UserID)
	a.NotNil(body.Event.Reward.ID)
}
func TestWebsubRedemption(t *testing.T) {
	a := util.SetupTestEnv(t)

	params := *&RedemptionParams{
		Transport: TransportWebSub,
	}

	_, err := GenerateRedemptionBody(params)
	a.NotNil(err, "Expected error (Channel Points unsupported on websub")
}
