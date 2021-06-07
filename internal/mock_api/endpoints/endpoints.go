// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package endpoints

import (
	"github.com/twitchdev/twitch-cli/internal/mock_api"
	"github.com/twitchdev/twitch-cli/internal/mock_api/endpoints/bits"
	"github.com/twitchdev/twitch-cli/internal/mock_api/endpoints/channel_points"
	"github.com/twitchdev/twitch-cli/internal/mock_api/endpoints/channels"
	"github.com/twitchdev/twitch-cli/internal/mock_api/endpoints/chat"
	"github.com/twitchdev/twitch-cli/internal/mock_api/endpoints/clips"
	"github.com/twitchdev/twitch-cli/internal/mock_api/endpoints/subscriptions"
	"github.com/twitchdev/twitch-cli/internal/mock_api/endpoints/users"
)

func All() []mock_api.MockEndpoint {
	return []mock_api.MockEndpoint{
		bits.BitsLeaderboard{},
		bits.Cheermotes{},
		channel_points.Reward{},
		channel_points.Redemption{},
		channels.CommercialEndpoint{},
		channels.Editors{},
		channels.InformationEndpoint{},
		chat.ChannelBadges{},
		chat.GlobalBadges{},
		clips.Clips{},
		users.UsersEndpoint{},
		users.FollowsEndpoint{},
		subscriptions.Endpoint{},
	}
}
