// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package trigger

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/rpc"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/twitchdev/twitch-cli/internal/database"
	"github.com/twitchdev/twitch-cli/internal/events"
	"github.com/twitchdev/twitch-cli/internal/events/mock_ws"
	"github.com/twitchdev/twitch-cli/internal/events/types"
	"github.com/twitchdev/twitch-cli/internal/models"
	"github.com/twitchdev/twitch-cli/internal/util"
)

// TriggerParameters defines the parameters used to emit an event.
type TriggerParameters struct {
	Event               string
	Transport           string
	IsAnonymous         bool
	FromUser            string
	ToUser              string
	GiftUser            string
	EventStatus         string
	SubscriptionStatus  string
	ItemID              string
	Cost                int64
	ForwardAddress      string
	Secret              string
	Verbose             bool
	Count               int
	Description         string
	ItemName            string
	GameID              string
	Tier                string
	Timestamp           string
	EventID             string // Also serves as subscription ID. See https://github.com/twitchdev/twitch-cli/issues/184
	CharityCurrentValue int
	CharityTargetValue  int
	WebSocketClient     string
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

	switch p.Tier {
	case "":
		p.Tier = "1000"
	case "1000", "2000", "3000":
		// do nothing, these are valid values
	default:
		return "", fmt.Errorf(
			"Discarding event: Invalid tier provided.\n" +
				"Valid values are 1000, 2000 or 3000")
	}

	if p.EventID == "" {
		p.EventID = util.RandomGUID()
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
		ID:                  p.EventID,
		Trigger:             p.Event,
		Transport:           p.Transport,
		FromUserID:          p.FromUser,
		FromUserName:        "testFromUser",
		ToUserID:            p.ToUser,
		ToUserName:          "testBroadcaster",
		IsAnonymous:         p.IsAnonymous,
		Cost:                p.Cost,
		EventStatus:         p.EventStatus,
		ItemID:              p.ItemID,
		Description:         p.Description,
		ItemName:            p.ItemName,
		GameID:              p.GameID,
		Tier:                p.Tier,
		SubscriptionStatus:  p.SubscriptionStatus,
		Timestamp:           p.Timestamp,
		CharityCurrentValue: p.CharityCurrentValue,
		CharityTargetValue:  p.CharityTargetValue,
	}

	e, err := types.GetByTriggerAndTransport(p.Event, p.Transport)
	if err != nil {
		return "", err
	}

	newTrigger := e.GetEventSubAlias(p.Event)
	if newTrigger != "" {
		eventParamaters.Trigger = newTrigger // overwrite the existing trigger with the "correct" one
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

	if p.ForwardAddress != "" && strings.EqualFold(p.Transport, "webhook") { // Forwarding to an address requires Webhook, as its done via HTTP
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

	// Forward to WebSocket server via RPC
	if strings.EqualFold(p.Transport, "websocket") {
		client, err := rpc.DialHTTP("tcp", ":44747")
		if err != nil {
			return "", errors.New("Failed to dial RPC handler for WebSocket server. Is it online?\nError: " + err.Error())
		}

		var reply bool

		// Modify transport
		modifiedTransportJSON := models.EventsubResponse{}
		err = json.Unmarshal([]byte(resp.JSON), &modifiedTransportJSON)
		if err != nil {
			return "", errors.New("Unexpected error unmarshling JSON before forwarding to WebSocket server: " + err.Error())
		}
		modifiedTransportJSON.Subscription.Transport.Method = "websocket"
		modifiedTransportJSON.Subscription.Transport.Callback = ""
		modifiedTransportJSON.Subscription.Transport.SessionID = "TBD-by-WebSocket-Server"
		rawModifiedTransportJSON, _ := json.Marshal(modifiedTransportJSON)
		resp.JSON = rawModifiedTransportJSON

		if strings.HasPrefix(p.Event, "mock.websocket") {
			// Trigger server command to emulate a specific feature, such as websocket reconnect emulation
			args := &mock_ws.RPCArgs{
				Body: string(resp.JSON),
			}

			err = client.Call("WebSocketServerRPC.ServerCommand", args, &reply)
		} else {
			// Trigger any EventSub subscription that's available over 1st party WebSocket connections
			args := &mock_ws.RPCArgs{
				Body:       string(resp.JSON),
				ClientName: p.WebSocketClient,
			}

			err = client.Call("WebSocketServerRPC.RemoteFireEventSub", args, &reply)
		}

		// Error checking for all possible RPC executions
		if err != nil {
			return "", errors.New("Failed to send via RPC to WebSocket server: " + err.Error())
		}

		if reply {
			color.New().Add(color.FgGreen).Println(fmt.Sprintf(`✔ Forwarded for use in mock EventSub WebSocket server`))
		} else {
			color.New().Add(color.FgRed).Println(fmt.Sprintf(`✗ EventSub WebSocket server failed to process event. See server for more details.`))
		}
	}

	return string(resp.JSON), nil
}
