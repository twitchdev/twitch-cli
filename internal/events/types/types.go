// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package types

import (
	"errors"

	"github.com/twitchdev/twitch-cli/internal/events"
	"github.com/twitchdev/twitch-cli/internal/events/types/authorization_revoke"
	"github.com/twitchdev/twitch-cli/internal/events/types/channel_points_redemption"
	"github.com/twitchdev/twitch-cli/internal/events/types/channel_points_reward"
	"github.com/twitchdev/twitch-cli/internal/events/types/cheer"
	"github.com/twitchdev/twitch-cli/internal/events/types/extension_transaction"
	"github.com/twitchdev/twitch-cli/internal/events/types/follow"
	"github.com/twitchdev/twitch-cli/internal/events/types/moderator_change"
	"github.com/twitchdev/twitch-cli/internal/events/types/raid"
	"github.com/twitchdev/twitch-cli/internal/events/types/stream_change"
	"github.com/twitchdev/twitch-cli/internal/events/types/streamdown"
	"github.com/twitchdev/twitch-cli/internal/events/types/streamup"
	"github.com/twitchdev/twitch-cli/internal/events/types/subscribe"
)

func All() []events.MockEvent {
	return []events.MockEvent{
		authorization_revoke.Event{},
		channel_points_redemption.Event{},
		channel_points_reward.Event{},
		cheer.Event{},
		extension_transaction.Event{},
		follow.Event{},
		raid.Event{},
		subscribe.Event{},
		stream_change.Event{},
		streamup.Event{},
		streamdown.Event{},
		moderator_change.Event{},
	}
}

func GetByTriggerAndTransport(trigger string, transport string) (events.MockEvent, error) {
	for _, e := range All() {
		if e.ValidTrigger(trigger) == true && e.ValidTransport(transport) == true {
			return e, nil
		}
	}

	return nil, errors.New("Invalid event/transport combination")
}
