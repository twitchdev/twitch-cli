// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package events

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/twitchdev/twitch-cli/internal/util"
)

func TestFire(t *testing.T) {
	a := util.SetupTestEnv(t)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusAccepted)

		_, err := ioutil.ReadAll(r.Body)
		a.Nil(err)
	}))
	defer ts.Close()

	params := *&TriggerParameters{
		Event:          "gift",
		Transport:      TransportEventSub,
		IsAnonymous:    false,
		FromUser:       "",
		ToUser:         "",
		GiftUser:       "",
		Status:         "",
		ItemId:         "",
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
		Transport:      TransportEventSub,
		IsAnonymous:    false,
		FromUser:       "",
		ToUser:         "",
		GiftUser:       "",
		Status:         "",
		ItemId:         "",
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
		Transport:      TransportEventSub,
		IsAnonymous:    false,
		FromUser:       "",
		ToUser:         "",
		GiftUser:       "",
		Status:         "",
		ItemId:         "",
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
		Event:          "cheer",
		Transport:      TransportEventSub,
		IsAnonymous:    false,
		FromUser:       "",
		ToUser:         "",
		GiftUser:       "",
		Status:         "",
		ItemId:         "",
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
		Transport:      TransportEventSub,
		IsAnonymous:    false,
		FromUser:       "",
		ToUser:         "",
		GiftUser:       "",
		Status:         "",
		ItemId:         "",
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
		Transport:      TransportEventSub,
		IsAnonymous:    false,
		FromUser:       "",
		ToUser:         "",
		GiftUser:       "",
		Status:         "",
		ItemId:         "",
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
		Transport:      TransportWebSub,
		IsAnonymous:    false,
		FromUser:       "",
		ToUser:         "",
		GiftUser:       "",
		Status:         "",
		ItemId:         "",
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
		Transport:      TransportEventSub,
		IsAnonymous:    false,
		FromUser:       "",
		ToUser:         "",
		GiftUser:       "",
		Status:         "",
		ItemId:         "",
		Cost:           0,
		ForwardAddress: ts.URL,
		Secret:         "potato",
		Verbose:        false,
		Count:          0,
	}
	res, err = Fire(params)
	a.NotNil(err)

	params = *&TriggerParameters{
		Event:          "add-reward",
		Transport:      "potato",
		IsAnonymous:    false,
		FromUser:       "",
		ToUser:         "",
		GiftUser:       "",
		Status:         "",
		ItemId:         "",
		Cost:           0,
		ForwardAddress: ts.URL,
		Secret:         "potato",
		Verbose:        false,
		Count:          0,
	}
	res, err = Fire(params)
	a.NotNil(err)
}
func TestValidTriggers(t *testing.T) {
	a := util.SetupTestEnv(t)

	t1 := ValidTriggers()
	a.NotEmpty(t1)
}

func TestValidTransports(t *testing.T) {
	a := util.SetupTestEnv(t)

	t1 := ValidTransports()
	a.NotEmpty(t1)
}
