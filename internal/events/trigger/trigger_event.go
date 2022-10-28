// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package trigger

import (
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/twitchdev/twitch-cli/internal/database"
	"github.com/twitchdev/twitch-cli/internal/events"
	"github.com/twitchdev/twitch-cli/internal/events/types"
	"github.com/twitchdev/twitch-cli/internal/models"
	"github.com/twitchdev/twitch-cli/internal/util"
)

// TriggerParameters defines the parameters used to emit an event.
type TriggerParameters struct {
	Event              string
	Transport          string
	IsAnonymous        bool
	FromUser           string
	ToUser             string
	GiftUser           string
	EventStatus        string
	SubscriptionStatus string
	ItemID             string
	Cost               int64
	ForwardAddress     string
	Secret             string
	Verbose            bool
	Count              int
	Description        string
	ItemName           string
	GameID             string
	Timestamp          string
}

type TriggerResponse struct {
	ID        string
	JSON      []byte
	FromUser  string
	ToUser    string
	Timestamp string
}

// Fire emits an event using the TriggerParameters defined above.
func Fire(p TriggerParameters) (string, error) {
	var resp events.MockEventResponse
	var err error

	if p.ToUser == "" {
		p.ToUser = util.RandomUserID()
	}

	if p.FromUser == "" {
		p.FromUser = util.RandomUserID()
	}

	if p.GameID == "" {
		p.GameID = fmt.Sprint(util.RandomInt(10 * 1000))
	}

	if p.Timestamp == "" {
		p.Timestamp = util.GetTimestamp().Format(time.RFC3339Nano)
	} else {
		// Verify custom timestamp
		_, err := time.Parse(time.RFC3339Nano, p.Timestamp)
		if err != nil {
			return "", fmt.Errorf(
				`Discarding event: Invalid timestamp provided.
Please follow RFC3339Nano, which is used by Twitch as seen here:
https://dev.twitch.tv/docs/eventsub/handling-webhook-events#processing-an-event`)
		}
	}

	eventParamaters := events.MockEventParameters{
		ID:                 util.RandomGUID(),
		Trigger:            p.Event,
		Transport:          p.Transport,
		FromUserID:         p.FromUser,
		FromUserName:       "testFromUser",
		ToUserID:           p.ToUser,
		ToUserName:         "testBroadcaster",
		IsAnonymous:        p.IsAnonymous,
		Cost:               p.Cost,
		EventStatus:        p.EventStatus,
		ItemID:             p.ItemID,
		Description:        p.Description,
		ItemName:           p.ItemName,
		GameID:             p.GameID,
		SubscriptionStatus: p.SubscriptionStatus,
		Timestamp:          p.Timestamp,
	}

	e, err := types.GetByTriggerAndTransport(p.Event, p.Transport)
	if err != nil {
		return "", err
	}

	if eventParamaters.Transport == models.TransportEventSub {
		newTrigger := e.GetEventSubAlias(p.Event)
		if newTrigger != "" {
			eventParamaters.Trigger = newTrigger // overwrite the existing trigger with the "correct" one
		}
	}

	resp, err = e.GenerateEvent(eventParamaters)
	if err != nil {
		return "", err
	}

	db, err := database.NewConnection()
	if err != nil {
		return "", err
	}

	err = db.NewQuery(nil, 100).InsertIntoDB(database.EventCacheParameters{
		ID:        resp.ID,
		Event:     p.Event,
		JSON:      string(resp.JSON),
		FromUser:  resp.FromUser,
		ToUser:    resp.ToUser,
		Transport: p.Transport,
		Timestamp: p.Timestamp,
	})
	if err != nil {
		return "", err
	}
	topic := e.GetTopic(p.Transport, p.Event)
	if topic == "" && e.GetEventSubAlias(p.Event) != "" {
		topic = p.Event
	}

	messageType := EventSubMessageTypeNotification
	// Set to "revocation" if SubscriptionStatus is not set to "enabled"
	// We don't have to worry about "webhook_callback_verification" in this bit of code, since it's an entirely different command. All this code is from "event trigger".
	if !strings.EqualFold(p.SubscriptionStatus, "enabled") {
		messageType = EventSubMessageTypeRevocation
	}

	if p.ForwardAddress != "" {
		resp, err := ForwardEvent(ForwardParamters{
			ID:                  resp.ID,
			Transport:           p.Transport,
			Timestamp:           p.Timestamp,
			JSON:                resp.JSON,
			Secret:              p.Secret,
			ForwardAddress:      p.ForwardAddress,
			Event:               topic,
			Type:                messageType,
			SubscriptionVersion: e.SubscriptionVersion(),
		})
		if err != nil {
			return "", err
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}

		respTrigger := string(body)
		if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
			color.New().Add(color.FgGreen).Println(fmt.Sprintf(`✔ Request Sent. Received Status Code: %v`, resp.StatusCode))
			color.New().Add(color.FgGreen).Println(fmt.Sprintf(`✔ Server Said: %s`, respTrigger))
		} else {
			color.New().Add(color.FgRed).Println(fmt.Sprintf(`✗ Invalid response. Received Status Code: %v`, resp.StatusCode))
			color.New().Add(color.FgRed).Println(fmt.Sprintf(`✗ Server Said: %s`, respTrigger))
		}
	}

	return string(resp.JSON), nil
}
