// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package trigger

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/twitchdev/twitch-cli/internal/models"
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

		body, err := io.ReadAll(r.Body)
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
}

func TestForwardEventSignatureUsesTimestampHeader(t *testing.T) {
	a := test_setup.SetupTestEnv(t)

	secret := "potaytoes"
	id := "test-message-id"

	// Valid RFC3339Nano timestamps that are not in Go's canonical output form
	timestamps := []string{
		"2023-01-01T00:00:00.000Z",
		"2023-01-01T00:00:00.120Z",
		"2023-01-01T00:00:00+00:00",
	}

	for _, timestamp := range timestamps {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			body, err := io.ReadAll(r.Body)
			a.Nil(err)

			// Signature is HMAC-SHA256 over message id + timestamp header + body, as documented by Twitch
			mac := hmac.New(sha256.New, []byte(secret))
			mac.Write([]byte(r.Header.Get("Twitch-Eventsub-Message-Id")))
			mac.Write([]byte(r.Header.Get("Twitch-Eventsub-Message-Timestamp")))
			mac.Write(body)

			a.Equal(timestamp, r.Header.Get("Twitch-Eventsub-Message-Timestamp"))
			a.Equal(fmt.Sprintf("sha256=%x", mac.Sum(nil)), r.Header.Get("Twitch-Eventsub-Message-Signature"), timestamp)
			w.WriteHeader(http.StatusAccepted)
		}))

		resp, err := ForwardEvent(ForwardParamters{
			ID:                  id,
			ForwardAddress:      ts.URL,
			JSON:                []byte(`{"test":true}`),
			Transport:           models.TransportWebhook,
			Timestamp:           timestamp,
			Secret:              secret,
			Event:               "channel.follow",
			Type:                EventSubMessageTypeNotification,
			SubscriptionVersion: "2",
		})
		a.Nil(err)
		a.Equal(http.StatusAccepted, resp.StatusCode)
		resp.Body.Close()
		ts.Close()
	}
}
