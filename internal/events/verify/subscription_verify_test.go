// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package verify

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/twitchdev/twitch-cli/internal/models"
	"github.com/twitchdev/twitch-cli/test_setup"
)

func TestSubscriptionVerify(t *testing.T) {
	a := test_setup.SetupTestEnv(t)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var challenge string
		w.WriteHeader(http.StatusAccepted)

		body, err := ioutil.ReadAll(r.Body)
		a.Nil(err)

		if r.Method == http.MethodPost {
			var verification models.EventsubSubscriptionVerification
			err := json.Unmarshal(body, &verification)
			a.Nil(err)
			challenge = verification.Challenge

			if r.URL.Path == "/badendpoint" {
				challenge = "badresponse"
			}
		} else if r.Method == http.MethodGet {
			q := r.URL.Query()

			challenge = q.Get("hub.challenge")
		}

		w.Header().Add("Content-Type", "text/plain")
		w.Write([]byte(challenge))
	}))
	defer ts.Close()

	p := VerifyParameters{
		Transport:      models.TransportWebhook,
		Event:          "subscribe",
		ForwardAddress: ts.URL,
		Secret:         "potatoes",
	}
	res, err := VerifyWebhookSubscription(p)
	a.Nil(err)
	a.Equal(res.IsChallengeValid, true)

	p.Transport = "notarealtransport"
	_, err = VerifyWebhookSubscription(p)
	a.NotNil(err)

	p = VerifyParameters{
		ForwardAddress: ts.URL + "/badendpoint",
		Transport:      models.TransportWebhook,
		Event:          "subscribe",
		Secret:         "potatoes",
	}
	res, err = VerifyWebhookSubscription(p)
	a.Nil(err)
	a.Equal(res.IsChallengeValid, false)
}
