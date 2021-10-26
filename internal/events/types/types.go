// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package types

import (
	"errors"

	"github.com/twitchdev/twitch-cli/internal/events"
	"github.com/twitchdev/twitch-cli/internal/events/types/authorization"
	"github.com/twitchdev/twitch-cli/internal/events/types/ban"
	"github.com/twitchdev/twitch-cli/internal/events/types/channel_points_redemption"
	"github.com/twitchdev/twitch-cli/internal/events/types/channel_points_reward"
	"github.com/twitchdev/twitch-cli/internal/events/types/cheer"
	"github.com/twitchdev/twitch-cli/internal/events/types/drop"
	"github.com/twitchdev/twitch-cli/internal/events/types/extension_transaction"
	"github.com/twitchdev/twitch-cli/internal/events/types/follow"
	"github.com/twitchdev/twitch-cli/internal/events/types/gift"
	"github.com/twitchdev/twitch-cli/internal/events/types/hype_train"
	"github.com/twitchdev/twitch-cli/internal/events/types/moderator_change"
	"github.com/twitchdev/twitch-cli/internal/events/types/poll"
	"github.com/twitchdev/twitch-cli/internal/events/types/prediction"
	"github.com/twitchdev/twitch-cli/internal/events/types/raid"
	"github.com/twitchdev/twitch-cli/internal/events/types/stream_change"
	"github.com/twitchdev/twitch-cli/internal/events/types/streamdown"
	"github.com/twitchdev/twitch-cli/internal/events/types/streamup"
	"github.com/twitchdev/twitch-cli/internal/events/types/subscribe"
	"github.com/twitchdev/twitch-cli/internal/events/types/subscription_message"
	"github.com/twitchdev/twitch-cli/internal/models"
)

func All() []events.MockEvent {
	return []events.MockEvent{
		authorization.Event{},
		ban.Event{},
		channel_points_redemption.Event{},
		channel_points_reward.Event{},
		cheer.Event{},
		drop.Event{},
		extension_transaction.Event{},
		follow.Event{},
		gift.Event{},
		hype_train.Event{},
		moderator_change.Event{},
		poll.Event{},
		prediction.Event{},
		raid.Event{},
		stream_change.Event{},
		streamup.Event{},
		streamdown.Event{},
		subscribe.Event{},
		subscription_message.Event{},
	}
}

func GetByTriggerAndTransport(trigger string, transport string) (events.MockEvent, error) {
	for _, e := range All() {
		if transport == models.TransportEventSub {
			newTrigger := e.GetEventSubAlias(trigger)
			if newTrigger != "" {
				trigger = newTrigger
			}
		}
		if e.ValidTrigger(trigger) == true && e.ValidTransport(transport) == true {
			return e, nil
		}
	}

	return nil, errors.New("Invalid event")
}
