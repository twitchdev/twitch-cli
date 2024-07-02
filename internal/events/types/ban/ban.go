// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package ban

import (
	"encoding/json"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/twitchdev/twitch-cli/internal/events"
	"github.com/twitchdev/twitch-cli/internal/models"
	"github.com/twitchdev/twitch-cli/internal/util"
)

var transportsSupported = map[string]bool{
	models.TransportWebhook:   true,
	models.TransportWebSocket: true,
}

var triggerSupported = []string{"ban"}

var triggerMapping = map[string]map[string]string{
	models.TransportWebhook: {
		"ban": "channel.ban",
	},
	models.TransportWebSocket: {
		"ban": "channel.ban",
	},
}

type Event struct{}

func (e Event) GenerateEvent(params events.MockEventParameters) (events.MockEventResponse, error) {
	var event []byte
	var err error

	switch params.Transport {
	case models.TransportWebhook, models.TransportWebSocket:
		ban := models.BanEventSubEvent{
			UserID:               params.FromUserID,
			UserLogin:            params.FromUserName,
			UserName:             params.FromUserName,
			BroadcasterUserID:    params.ToUserID,
			BroadcasterUserLogin: params.ToUserName,
			BroadcasterUserName:  params.ToUserName,
			ModeratorUserId:      util.RandomUserID(),
			ModeratorUserLogin:   "CLIModerator",
			ModeratorUserName:    "CLIModerator",
		}

		reason := "This is a test event"

		// This event supports --timestamp historically, but is overridden by the newer --ban-start
		bannedAt := params.Timestamp
		if params.BanStartTimestamp != "" {
			bannedAt = params.BanStartTimestamp
		}

		var endsAt *string = nil
		var isPermanent bool

		if params.BanEndTimestamp == "" {
			// Default to perma ban
			isPermanent = true
		} else {
			r1 := regexp.MustCompile("^[0-9]+$")
			r2 := regexp.MustCompile("^(?:(?P<Days>[0-9]+)[dD])?(?:(?P<Hours>[0-9]+)[hH])?(?:(?P<Minutes>[0-9]+)[mM])?(?:(?P<Seconds>[0-9]+)[sS])?$")

			if r1.MatchString(params.BanEndTimestamp) {
				// Similar format to /timeout <user> <seconds>
				// twitch event trigger channel.ban --ban-end=600
				seconds, _ := strconv.Atoi(r1.FindAllString(params.BanEndTimestamp, -1)[0])
				tNow, _ := time.Parse(time.RFC3339Nano, params.Timestamp)
				tLater := tNow.Add(time.Duration(seconds) * time.Second).Format(time.RFC3339Nano)
				endsAt = &tLater
				isPermanent = false

			} else if r2.MatchString(params.BanEndTimestamp) {
				// Relative time specified by shorthands. e.g. 90d10h30m45s
				// Can include or exclude any of those, but they have to be in the same order as above
				values := r2.FindStringSubmatch(params.BanEndTimestamp)
				days, _ := strconv.Atoi(values[r2.SubexpIndex("Days")])
				hours, _ := strconv.Atoi(values[r2.SubexpIndex("Hours")])
				minutes, _ := strconv.Atoi(values[r2.SubexpIndex("Minutes")])
				seconds, _ := strconv.Atoi(values[r2.SubexpIndex("Seconds")])

				tNow, _ := time.Parse(time.RFC3339Nano, params.Timestamp)
				tLater := tNow.Add(time.Duration(days*24) * time.Hour).
					Add(time.Duration(hours) * time.Hour).
					Add(time.Duration(minutes) * time.Minute).
					Add(time.Duration(seconds) * time.Second).
					Format(time.RFC3339Nano)
				endsAt = &tLater
				isPermanent = false

			} else {
				// Timeout with user provided timestamp
				endsAt = &params.BanEndTimestamp
				isPermanent = false

			}
		}

		ban.Reason = reason
		ban.BannedAt = bannedAt
		ban.EndsAt = endsAt
		ban.IsPermanent = isPermanent

		body := models.EventsubResponse{
			Subscription: models.EventsubSubscription{
				ID:      params.SubscriptionID,
				Status:  params.SubscriptionStatus,
				Type:    triggerMapping[params.Transport][params.Trigger],
				Version: e.SubscriptionVersion(),
				Condition: models.EventsubCondition{
					BroadcasterUserID: params.ToUserID,
				},
				Transport: models.EventsubTransport{
					Method:   "webhook",
					Callback: "null",
				},
				Cost:      0,
				CreatedAt: params.Timestamp,
			},
			Event: ban,
		}

		event, err = json.Marshal(body)
		if err != nil {
			return events.MockEventResponse{}, err
		}

		// Delete event info if Subscription.Status is not set to "enabled"
		if !strings.EqualFold(params.SubscriptionStatus, "enabled") {
			var i interface{}
			if err := json.Unmarshal([]byte(event), &i); err != nil {
				return events.MockEventResponse{}, err
			}
			if m, ok := i.(map[string]interface{}); ok {
				delete(m, "event") // Matches JSON key defined in body variable above
			}

			event, err = json.Marshal(i)
			if err != nil {
				return events.MockEventResponse{}, err
			}
		}
	default:
		return events.MockEventResponse{}, nil
	}

	return events.MockEventResponse{
		ID:       params.EventMessageID,
		JSON:     event,
		FromUser: params.FromUserID,
		ToUser:   params.ToUserID,
	}, nil
}

func (e Event) ValidTransport(t string) bool {
	return transportsSupported[t]
}

func (e Event) ValidTrigger(t string) bool {
	for _, ts := range triggerSupported {
		if ts == t {
			return true
		}
	}
	return false
}

func (e Event) GetTopic(transport string, trigger string) string {
	return triggerMapping[transport][trigger]
}
func (e Event) GetAllTopicsByTransport(transport string) []string {
	allTopics := []string{}
	for _, topic := range triggerMapping[transport] {
		allTopics = append(allTopics, topic)
	}
	return allTopics
}
func (e Event) GetEventSubAlias(t string) string {
	// check for aliases
	for trigger, topic := range triggerMapping[models.TransportWebhook] {
		if topic == t {
			return trigger
		}
	}
	return ""
}

func (e Event) SubscriptionVersion() string {
	return "1"
}
