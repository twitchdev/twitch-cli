// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package trigger

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/twitchdev/twitch-cli/test_setup"
)

func TestForwardEventEventsub(t *testing.T) {
	a := test_setup.SetupTestEnv(t)

	secret := "potaytoes"

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusAccepted)

		a.NotEmpty(r.Header.Get("Twitch-Eventsub-Message-Retry"))
		a.NotEmpty(r.Header.Get("Twitch-Eventsub-Subscription-Version"))
		a.NotEmpty(r.Header.Get("Twitch-Eventsub-Message-Type"))
		a.NotEmpty(r.Header.Get("Twitch-Eventsub-Subscription-Type"))
		a.NotEmpty(r.Header.Get("Twitch-Eventsub-Message-Id"))

		body, err := ioutil.ReadAll(r.Body)
		a.Nil(err)
		a.NotNil(body)

		mac := hmac.New(sha256.New, []byte(secret))
		timestamp, err := time.Parse(time.RFC3339Nano, r.Header.Get("Twitch-Eventsub-Message-Timestamp"))
		a.Nil(err)

		id := r.Header.Get("Twitch-Eventsub-Message-Id")

		mac.Write(timestamp.AppendFormat([]byte(id), time.RFC3339Nano))
		mac.Write(body)

		hash := fmt.Sprintf("sha256=%x", mac.Sum(nil))
		a.Equal(hash, r.Header.Get("Twitch-Eventsub-Message-Signature"))
	}))
	defer ts.Close()

	// TODO update test
}

func TestForwardEventWebsub(t *testing.T) {
	a := test_setup.SetupTestEnv(t)

	secret := "potaytoes"

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusAccepted)

		body, err := ioutil.ReadAll(r.Body)
		a.Nil(err)
		a.NotNil(body)

		mac := hmac.New(sha256.New, []byte(secret))

		mac.Write(body)

		hash := fmt.Sprintf("sha256=%x", mac.Sum(nil))
		a.Equal(hash, r.Header.Get("X-Hub-Signature"))
	}))
	defer ts.Close()

	// TODO update test

}
