// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package trigger

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/twitchdev/twitch-cli/internal/models"
	"github.com/twitchdev/twitch-cli/internal/request"
)

type ForwardParamters struct {
	ID                  string
	ForwardAddress      string
	JSON                []byte
	Transport           string
	Timestamp           string
	Secret              string
	Event               string
	EventMessageID      string
	Method              string
	Type                string
	SubscriptionVersion string
}

type header struct {
	HeaderName  string
	HeaderValue string
}

const (
	EventSubMessageTypeNotification = "notification"
	EventSubMessageTypeVerification = "webhook_callback_verification"
	EventSubMessageTypeRevocation   = "revocation"
)

var notificationHeaders = map[string][]header{
	models.TransportWebhook: {
		{
			HeaderName:  `Twitch-Eventsub-Message-Retry`,
			HeaderValue: `0`,
		},
	},
}

func ForwardEvent(p ForwardParamters) (*http.Response, error) {
	method := http.MethodPost
	if p.Method != "" {
		method = p.Method
	}

	req, err := request.NewRequest(method, p.ForwardAddress, bytes.NewBuffer(p.JSON))
	if err != nil {
		return &http.Response{}, err
	}

	req.Header.Set("Content-Type", "application/json")
	for _, header := range notificationHeaders[p.Transport] {
		req.Header.Add(header.HeaderName, header.HeaderValue)
	}

	switch p.Transport {
	case models.TransportWebhook:
		req.Header.Set("Twitch-Eventsub-Message-Id", p.ID)
		req.Header.Set("Twitch-Eventsub-Subscription-Type", p.Event)
		req.Header.Set("Twitch-Eventsub-Subscription-Version", p.SubscriptionVersion)
		switch p.Type {
		case EventSubMessageTypeNotification:
			req.Header.Add("Twitch-Eventsub-Message-Type", EventSubMessageTypeNotification)
		case EventSubMessageTypeVerification:
			req.Header.Add("Twitch-Eventsub-Message-Type", EventSubMessageTypeVerification)
		case EventSubMessageTypeRevocation:
			req.Header.Add("Twitch-Eventsub-Message-Type", EventSubMessageTypeRevocation)
		}
	}

	// Twitch only supports IPv4 currently, so we will force this TCP connection to only use IPv4
	var dialer net.Dialer
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
		return dialer.DialContext(ctx, "tcp4", addr)
	}

	if p.Secret != "" {
		getSignatureHeader(req, p.ID, p.Secret, p.Transport, p.Timestamp, p.JSON)
	}

	client := &http.Client{
		Timeout: time.Second * 10,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Transport: transport,
	}
	resp, err := client.Do(req)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

func getSignatureHeader(req *http.Request, id string, secret string, transport string, timestamp string, payload []byte) {
	mac := hmac.New(sha256.New, []byte(secret))
	ts, _ := time.Parse(time.RFC3339Nano, timestamp)

	switch transport {
	case models.TransportWebhook:
		req.Header.Set("Twitch-Eventsub-Message-Timestamp", timestamp)
		prefix := ts.AppendFormat([]byte(id), time.RFC3339Nano)
		mac.Write(prefix)
		mac.Write(payload)
		req.Header.Set("Twitch-Eventsub-Message-Signature", fmt.Sprintf("sha256=%x", mac.Sum(nil)))
	}
}
