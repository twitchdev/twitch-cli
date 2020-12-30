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
)

func TestForwardEventEventsub(t *testing.T) {
	secret := "potaytoes"

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusAccepted)

		if r.Header.Get("Twitch-Eventsub-Message-Retry") == "" || r.Header.Get("Twitch-Eventsub-Subscription-Version") == "" || r.Header.Get("Twitch-Eventsub-Message-Type") == "" || r.Header.Get("Twitch-Eventsub-Subscription-Type") == "" || r.Header.Get("Twitch-Eventsub-Message-Id") == "" {
			t.Error("Missing Eventsub headers")
		}

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Error(err)
		}

		if body == nil {
			t.Errorf("Expected request to have body, got nil")
		}

		mac := hmac.New(sha256.New, []byte(secret))
		timestamp, err := time.Parse(time.RFC3339, r.Header.Get("Twitch-Eventsub-Message-Timestamp"))
		if err != nil {
			t.Error(err)
		}

		id := r.Header.Get("Twitch-Eventsub-Message-Id")

		mac.Write(timestamp.AppendFormat([]byte(id), time.RFC3339))
		mac.Write(body)

		hash := fmt.Sprintf("sha256=%x", mac.Sum(nil))

		if hash != r.Header.Get("Twitch-Eventsub-Message-Signature") {
			t.Error("Signature verification failure")
		}

	}))
	defer ts.Close()

	sParams := SubscribeParams{
		Transport: "eventsub",
		Type:      "channel.subscribe",
	}

	event, err := GenerateSubBody(sParams)
	if err != nil {
		t.Error(err)
	}

	fParams := ForwardParamters{
		ID:             event.ID,
		ForwardAddress: ts.URL,
		JSON:           event.JSON,
		Transport:      "eventsub",
		Event:          sParams.Type,
		Secret:         secret,
	}

	_, err = forwardEvent(fParams)
	if err != nil {
		t.Error(err)
	}
}

func TestForwardEventWebsub(t *testing.T) {
	secret := "potaytoes"

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusAccepted)

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Error(err)
		}

		if body == nil {
			t.Errorf("Expected request to have body, got nil")
		}

		mac := hmac.New(sha256.New, []byte(secret))

		mac.Write(body)

		hash := fmt.Sprintf("sha256=%x", mac.Sum(nil))

		if hash != r.Header.Get("X-Hub-Signature") {
			t.Errorf("Signature verification failure, got %v and expected %v", r.Header.Get("X-Hub-Signature"), hash)
		}

	}))
	defer ts.Close()

	sParams := SubscribeParams{
		Transport: "webusb",
		Type:      "subscribe",
	}

	event, err := GenerateSubBody(sParams)
	if err != nil {
		t.Error(err)
	}

	fParams := ForwardParamters{
		ID:             event.ID,
		ForwardAddress: ts.URL,
		JSON:           event.JSON,
		Transport:      "websub",
		Event:          sParams.Type,
		Secret:         secret,
	}

	_, err = forwardEvent(fParams)
	if err != nil {
		t.Error(err)
	}
}
