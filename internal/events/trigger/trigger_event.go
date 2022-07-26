// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package trigger

import (
	"fmt"
	"io/ioutil"
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
	Event          string
	Transport      string
	IsAnonymous    bool
	FromUser       string
	ToUser         string
	GiftUser       string
	Status         string
	ItemID         string
	Cost           int64
	ForwardAddress string
	Secret         string
	Verbose        bool
	Count          int
	Description    string
	ItemName       string
	GameID         string
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
	eventParamaters := events.MockEventParameters{
		ID:           util.RandomGUID(),
		Trigger:      p.Event,
		Transport:    p.Transport,
		FromUserID:   p.FromUser,
		FromUserName: "testFromUser",
		ToUserID:     p.ToUser,
		ToUserName:   "testBroadcaster",
		IsAnonymous:  p.IsAnonymous,
		Cost:         p.Cost,
		Status:       p.Status,
		ItemID:       p.ItemID,
		Description:  p.Description,
		ItemName:     p.ItemName,
		GameID:       p.GameID,
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
		Timestamp: util.GetTimestamp().Format(time.RFC3339Nano),
	})
	if err != nil {
		return "", err
	}

	if p.ForwardAddress != "" {
		resp, err := ForwardEvent(ForwardParamters{
			ID:             resp.ID,
			Transport:      p.Transport,
			JSON:           resp.JSON,
			Secret:         p.Secret,
			ForwardAddress: p.ForwardAddress,
			Event:          p.Event,
			Type:           EventSubMessageTypeNotification,
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
