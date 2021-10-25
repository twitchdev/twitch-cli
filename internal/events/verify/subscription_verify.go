// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package verify

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/fatih/color"
	"github.com/twitchdev/twitch-cli/internal/events/trigger"
	"github.com/twitchdev/twitch-cli/internal/events/types"
	"github.com/twitchdev/twitch-cli/internal/models"
	"github.com/twitchdev/twitch-cli/internal/util"
)

type VerifyParameters struct {
	Transport      string
	Event          string
	ForwardAddress string
	Secret         string
}

type VerifyResponse struct {
	IsStatusValid    bool
	IsChallengeValid bool
	Body             string
}

func VerifyWebhookSubscription(p VerifyParameters) (VerifyResponse, error) {
	r := VerifyResponse{}

	challenge := util.RandomGUID()

	event, err := types.GetByTriggerAndTransport(p.Event, p.Transport)
	if err != nil {
		return VerifyResponse{}, err
	}

	if p.Transport == models.TransportEventSub {
		newTrigger := event.GetEventbusAlias(p.Event)
		if newTrigger != "" {
			p.Event = newTrigger
		}
	}

	body, err := generateWebhookSubscriptionBody(p.Transport, event.GetTopic(p.Transport, p.Event), challenge, p.ForwardAddress)
	if err != nil {
		return VerifyResponse{}, err
	}

	r.Body = string(body.JSON)

	if p.ForwardAddress != "" {
		requestMethod := http.MethodPost
		u, err := url.Parse(p.ForwardAddress)
		if err != nil {
			return VerifyResponse{}, err
		}

		if p.Transport == models.TransportWebSub {
			q := u.Query()
			q.Add("hub.challenge", challenge)
			// this isn't per spec, however for the purposes of verifying whether a service is responding properly, it'll do
			q.Add("hub.topic", event.GetTopic(p.Transport, p.Event))
			q.Add("hub.mode", "subscribe")
			u.RawQuery = q.Encode()
			requestMethod = http.MethodGet
		}

		resp, err := trigger.ForwardEvent(trigger.ForwardParamters{
			ID:             body.ID,
			Event:          event.GetTopic(p.Transport, p.Event),
			JSON:           body.JSON,
			Transport:      p.Transport,
			Secret:         p.Secret,
			Method:         requestMethod,
			ForwardAddress: u.String(),
			Type:           trigger.EventSubMessageTypeVerification,
		})
		if err != nil {
			return VerifyResponse{}, err
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return VerifyResponse{}, err
		}

		respChallenge := string(body)
		if respChallenge == challenge {
			color.New().Add(color.FgGreen).Println(fmt.Sprintf(`✔ Valid response. Received challenge %s in body`, challenge))
			r.IsChallengeValid = true
		} else {
			color.New().Add(color.FgRed).Println(fmt.Sprintf(`✗ Invalid response. Received %s as body, expected %s`, respChallenge, challenge))
			r.IsChallengeValid = false
		}

		if resp.Header.Get("Content-Type") == "text/plain" {
			color.New().Add(color.FgGreen).Println(fmt.Sprintf(`✔ Valid content-type header. Received type %v`, resp.Header.Get("Content-Type")))
		} else {
			color.New().Add(color.FgRed).Println(fmt.Sprintf(`✗ Invalid content-type header. Received type %v`, resp.Header.Get("Content-Type")))
		}

		if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
			color.New().Add(color.FgGreen).Println(fmt.Sprintf(`✔ Valid status code. Received status %v`, resp.StatusCode))
			r.IsStatusValid = true
		} else {
			color.New().Add(color.FgRed).Println(fmt.Sprintf(`✗ Invalid status code. Received %v, expected a 2XX status`, resp.StatusCode))
			r.IsStatusValid = false
		}
	}

	return r, nil
}

func generateWebhookSubscriptionBody(transport string, event string, challenge string, callback string) (trigger.TriggerResponse, error) {
	var res []byte
	var err error
	id := util.RandomGUID()
	ts := util.GetTimestamp().Format(time.RFC3339Nano)
	switch transport {
	case models.TransportEventSub:
		body := models.EventsubSubscriptionVerification{
			Challenge: challenge,
			Subscription: models.EventsubSubscription{
				ID:      id,
				Status:  "webhook_callback_verification_pending",
				Type:    event,
				Version: "1",
				Condition: models.EventsubCondition{
					BroadcasterUserID: util.RandomUserID(),
				},
				Transport: models.EventsubTransport{
					Method:   "webhook",
					Callback: callback,
				},
				CreatedAt: ts,
			},
		}
		res, err = json.Marshal(body)
		if err != nil {
			return trigger.TriggerResponse{}, err
		}
	case models.TransportWebSub:

	default:
		res = []byte("")
	}
	return trigger.TriggerResponse{
		ID:        id,
		JSON:      res,
		Timestamp: ts,
	}, nil
}
