// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package endpoints

import (
	"github.com/twitchdev/twitch-cli/internal/mock_api"
	"github.com/twitchdev/twitch-cli/internal/mock_api/endpoints/bits"
	"github.com/twitchdev/twitch-cli/internal/mock_api/endpoints/categories"
	"github.com/twitchdev/twitch-cli/internal/mock_api/endpoints/channel_points"
	"github.com/twitchdev/twitch-cli/internal/mock_api/endpoints/channels"
	"github.com/twitchdev/twitch-cli/internal/mock_api/endpoints/chat"
	"github.com/twitchdev/twitch-cli/internal/mock_api/endpoints/clips"
	"github.com/twitchdev/twitch-cli/internal/mock_api/endpoints/drops"
	"github.com/twitchdev/twitch-cli/internal/mock_api/endpoints/hype_train"
	"github.com/twitchdev/twitch-cli/internal/mock_api/endpoints/moderation"
	"github.com/twitchdev/twitch-cli/internal/mock_api/endpoints/search"
	"github.com/twitchdev/twitch-cli/internal/mock_api/endpoints/streams"
	"github.com/twitchdev/twitch-cli/internal/mock_api/endpoints/subscriptions"
	"github.com/twitchdev/twitch-cli/internal/mock_api/endpoints/teams"
	"github.com/twitchdev/twitch-cli/internal/mock_api/endpoints/users"
	"github.com/twitchdev/twitch-cli/internal/mock_api/endpoints/videos"
)

func All() []mock_api.MockEndpoint {
	return []mock_api.MockEndpoint{
		bits.BitsLeaderboard{},
		bits.Cheermotes{},
		categories.Games{},
		categories.TopGames{},
		channel_points.Redemption{},
		channel_points.Reward{},
		channels.CommercialEndpoint{},
		channels.Editors{},
		channels.InformationEndpoint{},
		chat.ChannelBadges{},
		chat.GlobalBadges{},
		clips.Clips{},
		drops.DropsEntitlements{},
		hype_train.HypeTrainEvents{},
		moderation.AutomodHeld{},
		moderation.AutomodStatus{},
		moderation.BannedEvents{},
		moderation.Bans{},
		moderation.ModeratorEvents{},
		moderation.Moderators{},
		search.SearchCategories{},
		search.SearchChannels{},
		streams.AllTags{},
		streams.FollowedStreams{},
		streams.Markers{},
		streams.StreamKey{},
		streams.Streams{},
		streams.StreamTags{},
		subscriptions.BroadcasterSubscriptions{},
		subscriptions.UserSubscriptions{},
		teams.ChannelTeams{},
		teams.Teams{},
		users.Blocks{},
		users.FollowsEndpoint{},
		users.UsersEndpoint{},
		videos.Videos{},
	}
}
