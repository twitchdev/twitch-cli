// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package trigger

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/twitchdev/twitch-cli/internal/models"
	"github.com/twitchdev/twitch-cli/test_setup"
)

func TestFire(t *testing.T) {
	a := test_setup.SetupTestEnv(t)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusAccepted)

		_, err := ioutil.ReadAll(r.Body)
		a.Nil(err)
	}))
	defer ts.Close()

	params := *&TriggerParameters{
		Event:          "gift",
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

	res, err := Fire(params)
	a.Nil(err)
	a.NotEmpty(res)

	params = *&TriggerParameters{
		Event:          "cheer",
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
	res, err = Fire(params)
	a.Nil(err)
	a.NotEmpty(res)

	params = *&TriggerParameters{
		Event:          "follow",
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
		Version:        "2",
	}
	res, err = Fire(params)
	a.Nil(err)
	a.NotEmpty(res)

	params = *&TriggerParameters{
		Event:          "cheer",
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
	res, err = Fire(params)
	a.Nil(err)
	a.NotEmpty(res)

	params = *&TriggerParameters{
		Event:          "add-redemption",
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
	res, err = Fire(params)
	a.Nil(err)
	a.NotEmpty(res)

	params = *&TriggerParameters{
		Event:          "add-reward",
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
	res, err = Fire(params)
	a.Nil(err)
	a.NotEmpty(res)

	params = *&TriggerParameters{
		Event:          "transaction",
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
	res, err = Fire(params)
	a.Nil(err)
	a.NotEmpty(res)

	params = *&TriggerParameters{
		Event:          "add-reward",
		Transport:      "potato",
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
	res, err = Fire(params)
	a.NotNil(err)
}
