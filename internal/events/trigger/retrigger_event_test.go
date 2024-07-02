// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package trigger

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/twitchdev/twitch-cli/internal/models"
	"github.com/twitchdev/twitch-cli/test_setup"
	"github.com/twitchdev/twitch-cli/internal/util"
)

func TestRefireEvent(t *testing.T) {
	a := test_setup.SetupTestEnv(t)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusAccepted)

		_, err := io.ReadAll(r.Body)
		a.Nil(err)
	}))
	defer ts.Close()

	var eventMessageID = util.RandomGUID();

	params := TriggerParameters{
		Event:          "gift",
		EventMessageID: eventMessageID,
		Transport:      models.TransportWebhook,
		IsAnonymous:    false,
		FromUser:       "",
		ToUser:         "",
		GiftUser:       "",
		EventStatus:    "",
		ItemID:         "",
		Cost:           0,
		ForwardAddress: ts.URL,
		Secret:         "potato",
		Verbose:        false,
		Count:          0,
	}

	response, err := Fire(params)
	a.Nil(err)
	log.Print(err)
	var body models.SubEventSubResponse
	err = json.Unmarshal([]byte(response), &body)
	a.Nil(err)

	json, err := RefireEvent(eventMessageID, params)
	a.Nil(err)
	a.Equal(response, json)
}
