// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package endpoints

import (
	"github.com/twitchdev/twitch-cli/internal/mock_api"
	"github.com/twitchdev/twitch-cli/internal/mock_api/endpoints/bits"
	"github.com/twitchdev/twitch-cli/internal/mock_api/endpoints/categories"
	"github.com/twitchdev/twitch-cli/internal/mock_api/endpoints/ccl"
	"github.com/twitchdev/twitch-cli/internal/mock_api/endpoints/channel_points"
	"github.com/twitchdev/twitch-cli/internal/mock_api/endpoints/channels"
	"github.com/twitchdev/twitch-cli/internal/mock_api/endpoints/charity"
	"github.com/twitchdev/twitch-cli/internal/mock_api/endpoints/chat"
	"github.com/twitchdev/twitch-cli/internal/mock_api/endpoints/clips"
	"github.com/twitchdev/twitch-cli/internal/mock_api/endpoints/drops"
	"github.com/twitchdev/twitch-cli/internal/mock_api/endpoints/goals"
	"github.com/twitchdev/twitch-cli/internal/mock_api/endpoints/hype_train"
	"github.com/twitchdev/twitch-cli/internal/mock_api/endpoints/moderation"
	"github.com/twitchdev/twitch-cli/internal/mock_api/endpoints/polls"
	"github.com/twitchdev/twitch-cli/internal/mock_api/endpoints/predictions"
	"github.com/twitchdev/twitch-cli/internal/mock_api/endpoints/raids"
	"github.com/twitchdev/twitch-cli/internal/mock_api/endpoints/schedule"
	"github.com/twitchdev/twitch-cli/internal/mock_api/endpoints/search"
	"github.com/twitchdev/twitch-cli/internal/mock_api/endpoints/streams"
	"github.com/twitchdev/twitch-cli/internal/mock_api/endpoints/subscriptions"
	"github.com/twitchdev/twitch-cli/internal/mock_api/endpoints/teams"
	"github.com/twitchdev/twitch-cli/internal/mock_api/endpoints/users"
	"github.com/twitchdev/twitch-cli/internal/mock_api/endpoints/videos"
	"github.com/twitchdev/twitch-cli/internal/mock_api/endpoints/whispers"
)

func All() []mock_api.MockEndpoint {
	return []mock_api.MockEndpoint{
		bits.BitsLeaderboard{},
		bits.Cheermotes{},
		categories.Games{},
		categories.TopGames{},
		ccl.ContentClassificationLabels{},
		channel_points.Redemption{},
		channel_points.Reward{},
		channels.CommercialEndpoint{},
		channels.Editors{},
		channels.InformationEndpoint{},
		channels.Vips{},
		charity.CharityCampaign{},
		charity.CharityDonations{},
		chat.Announcements{},
		chat.ChannelBadges{},
		chat.ChannelEmotes{},
		chat.Chatters{},
		chat.Color{},
		chat.EmoteSets{},
		chat.GlobalBadges{},
		chat.GlobalEmotes{},
		chat.Settings{},
		chat.Shoutouts{},
		clips.Clips{},
		drops.DropsEntitlements{},
		goals.Goals{},
		hype_train.HypeTrainEvents{},
		moderation.AutomodHeld{},
		moderation.AutomodStatus{},
		moderation.Banned{},
		moderation.Bans{},
		moderation.Chat{},
		moderation.Moderators{},
		moderation.ShieldMode{},
		polls.Polls{},
		predictions.Predictions{},
		raids.Raids{},
		schedule.Schedule{},
		schedule.ScheduleICal{},
		schedule.ScheduleSegment{},
		schedule.ScheduleSettings{},
		search.SearchCategories{},
		search.SearchChannels{},
		streams.FollowedStreams{},
		streams.Markers{},
		streams.StreamKey{},
		streams.Streams{},
		subscriptions.BroadcasterSubscriptions{},
		subscriptions.UserSubscriptions{},
		teams.ChannelTeams{},
		teams.Teams{},
		users.Blocks{},
		users.FollowsEndpoint{},
		users.UsersEndpoint{},
		videos.Videos{},
		whispers.Whispers{},
	}
}

// All these endpoints return 410 Gone
func Gone() map[string][]string {
	return map[string][]string{
		"/tags/streams": {
			"GET",
		},
		"/streams/tags": {
			"GET",
			"POST",
			"PUT",
		},
	}
}
