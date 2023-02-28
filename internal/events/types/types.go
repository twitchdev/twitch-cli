// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package types

import (
	"errors"
	"sort"

	"github.com/twitchdev/twitch-cli/internal/events"
	"github.com/twitchdev/twitch-cli/internal/events/types/authorization"
	"github.com/twitchdev/twitch-cli/internal/events/types/ban"
	"github.com/twitchdev/twitch-cli/internal/events/types/channel_points_redemption"
	"github.com/twitchdev/twitch-cli/internal/events/types/channel_points_reward"
	"github.com/twitchdev/twitch-cli/internal/events/types/charity"
	"github.com/twitchdev/twitch-cli/internal/events/types/cheer"
	"github.com/twitchdev/twitch-cli/internal/events/types/drop"
	"github.com/twitchdev/twitch-cli/internal/events/types/extension_transaction"
	"github.com/twitchdev/twitch-cli/internal/events/types/follow"
	"github.com/twitchdev/twitch-cli/internal/events/types/gift"
	"github.com/twitchdev/twitch-cli/internal/events/types/goal"
	"github.com/twitchdev/twitch-cli/internal/events/types/hype_train"
	"github.com/twitchdev/twitch-cli/internal/events/types/moderator_change"
	"github.com/twitchdev/twitch-cli/internal/events/types/poll"
	"github.com/twitchdev/twitch-cli/internal/events/types/prediction"
	"github.com/twitchdev/twitch-cli/internal/events/types/raid"
	"github.com/twitchdev/twitch-cli/internal/events/types/shield_mode"
	"github.com/twitchdev/twitch-cli/internal/events/types/shoutout"
	"github.com/twitchdev/twitch-cli/internal/events/types/stream_change"
	"github.com/twitchdev/twitch-cli/internal/events/types/streamdown"
	"github.com/twitchdev/twitch-cli/internal/events/types/streamup"
	"github.com/twitchdev/twitch-cli/internal/events/types/subscribe"
	"github.com/twitchdev/twitch-cli/internal/events/types/subscription_message"
	user_update "github.com/twitchdev/twitch-cli/internal/events/types/user"
	"github.com/twitchdev/twitch-cli/internal/models"
)

func AllEvents() []events.MockEvent {
	return []events.MockEvent{
		authorization.Event{},
		ban.Event{},
		channel_points_redemption.Event{},
		channel_points_reward.Event{},
		charity.Event{},
		cheer.Event{},
		drop.Event{},
		extension_transaction.Event{},
		follow.Event{},
		gift.Event{},
		goal.Event{},
		hype_train.Event{},
		moderator_change.Event{},
		poll.Event{},
		prediction.Event{},
		raid.Event{},
		shield_mode.Event{},
		shoutout.Event{},
		stream_change.Event{},
		streamup.Event{},
		streamdown.Event{},
		subscribe.Event{},
		subscription_message.Event{},
		user_update.Event{},
	}
}

func AllEventTopics() []string {
	allEvents := []string{}

	for _, e := range AllEvents() {
		for _, topic := range e.GetAllTopicsByTransport(models.TransportEventSub) {
			allEvents = append(allEvents, topic)
		}
	}

	// Sort the topics alphabetically
	sort.Strings(allEvents)

	return allEvents
}

func GetByTriggerAndTransport(trigger string, transport string) (events.MockEvent, error) {
	for _, e := range AllEvents() {
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
