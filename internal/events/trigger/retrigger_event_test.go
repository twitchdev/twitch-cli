// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package trigger

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/twitchdev/twitch-cli/internal/models"
	"github.com/twitchdev/twitch-cli/test_setup"
)

func TestRefireEvent(t *testing.T) {
	a := test_setup.SetupTestEnv(t)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusAccepted)

		_, err := ioutil.ReadAll(r.Body)
		a.Nil(err)
	}))
	defer ts.Close()

	params := *&TriggerParameters{
		Event:          "gift",
		Transport:      models.TransportEventSub,
		IsAnonymous:    false,
		FromUser:       "",
		ToUser:         "",
		GiftUser:       "",
		Status:         "",
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

	json, err := RefireEvent(body.Subscription.ID, params)
	a.Nil(err)
	a.Equal(response, json)
}
