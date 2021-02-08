// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package trigger

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"net/http"
	"time"

	"github.com/twitchdev/twitch-cli/internal/models"
	"github.com/twitchdev/twitch-cli/internal/request"
	"github.com/twitchdev/twitch-cli/internal/util"
)

type ForwardParamters struct {
	ID             string
	ForwardAddress string
	JSON           []byte
	Transport      string
	Secret         string
	Event          string
	Method         string
	Type           string
}

type header struct {
	HeaderName  string
	HeaderValue string
}

const (
	EventSubMessageTypeNotification = "notification"
	EventSubMessageTypeVerification = "webhook_callback_verification"
)

var notificationHeaders = map[string][]header{
	models.TransportEventSub: {
		{
			HeaderName:  `Twitch-Eventsub-Message-Retry`,
			HeaderValue: `0`,
		},
		{
			HeaderName:  `Twitch-Eventsub-Subscription-Version`,
			HeaderValue: `test`,
		},
	},
	models.TransportWebSub: {
		{
			HeaderName:  `Twitch-Notification-Timestamp`,
			HeaderValue: util.GetTimestamp().Format(time.RFC3339Nano),
		},
		{
			HeaderName:  `Twitch-Notification-Retry`,
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
	case models.TransportEventSub:
		req.Header.Set("Twitch-Eventsub-Message-Id", p.ID)
		req.Header.Set("Twitch-Eventsub-Subscription-Type", p.Event)
		switch p.Type {
		case EventSubMessageTypeNotification:
			req.Header.Add("Twitch-Eventsub-Message-Type", EventSubMessageTypeNotification)
		case EventSubMessageTypeVerification:
			req.Header.Add("Twitch-Eventsub-Message-Type", EventSubMessageTypeVerification)
		}
	case models.TransportWebSub:
		req.Header.Set("Twitch-Notification-Id", p.ID)
	}

	if p.Secret != "" {
		getSignatureHeader(req, p.ID, p.Secret, p.Transport, p.JSON)
	}

	client := &http.Client{
		Timeout: time.Second * 10,
	}
	resp, err := client.Do(req)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

func getSignatureHeader(req *http.Request, id string, secret string, transport string, payload []byte) {
	mac := hmac.New(sha256.New, []byte(secret))
	ts := util.GetTimestamp()
	switch transport {
	case models.TransportEventSub:
		req.Header.Set("Twitch-Eventsub-Message-Timestamp", ts.Format(time.RFC3339Nano))
		prefix := ts.AppendFormat([]byte(id), time.RFC3339Nano)
		mac.Write(prefix)
		mac.Write(payload)
		req.Header.Set("Twitch-Eventsub-Message-Signature", fmt.Sprintf("sha256=%x", mac.Sum(nil)))
	case models.TransportWebSub:
		mac.Write(payload)
		req.Header.Set("X-Hub-Signature", fmt.Sprintf("sha256=%x", mac.Sum(nil)))
	}
}
