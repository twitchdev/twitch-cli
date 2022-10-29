// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package trigger

import (
	"fmt"

	"github.com/twitchdev/twitch-cli/internal/database"
	"github.com/twitchdev/twitch-cli/internal/events/types"
)

func RefireEvent(id string, p TriggerParameters) (string, error) {
	db, err := database.NewConnection()
	if err != nil {
		return "", err
	}
	res, err := db.NewQuery(nil, 100).GetEventByID(id)
	if err != nil {
		return "", err
	}

	p.Transport = res.Transport

	e, err := types.GetByTriggerAndTransport(res.Event, p.Transport)
	if err != nil {
		return "", err
	}

	topic := e.GetTopic(p.Transport, res.Event)
	if topic == "" && e.GetEventSubAlias(res.Event) != "" {
		topic = res.Event
	}

	if p.ForwardAddress != "" {
		resp, err := ForwardEvent(ForwardParamters{
			ID:                  id,
			Transport:           res.Transport,
			Timestamp:           p.Timestamp,
			ForwardAddress:      p.ForwardAddress,
			Secret:              p.Secret,
			JSON:                []byte(res.JSON),
			Event:               topic,
			Type:                EventSubMessageTypeNotification,
			SubscriptionVersion: e.SubscriptionVersion(),
		})
		defer resp.Body.Close()

		if err != nil {
			return "", err
		}
		fmt.Printf("[%v] Endpoint received refired event.", resp.StatusCode)
	}

	return res.JSON, nil
}
