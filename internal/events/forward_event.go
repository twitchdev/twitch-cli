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
}

type header struct {
	HeaderName  string
	HeaderValue string
}

var notificationHeaders = map[string][]header{
	TransportEventSub: {
		{
			HeaderName:  `Twitch-Eventsub-Message-Retry`,
			HeaderValue: `0`,
		},
		{
			HeaderName:  `Twitch-Eventsub-Message-Type`,
			HeaderValue: `notification`,
		},
		{
			HeaderName:  `Twitch-Eventsub-Subscription-Version`,
			HeaderValue: `test`,
		},
	},
	TransportWebSub: {
		{
			HeaderName:  `Twitch-Notification-Timestamp`,
			HeaderValue: util.GetTimestamp().Format(time.RFC3339),
		},
		{
			HeaderName:  `Twitch-Notification-Retry`,
			HeaderValue: `0`,
		},
	},
}

func forwardEvent(p ForwardParamters) (int, error) {
	req, err := request.NewRequest(http.MethodPost, p.ForwardAddress, bytes.NewBuffer(p.JSON))
	if err != nil {
		return 0, err
	}

	req.Header.Set("Content-Type", "application/json")
	for _, header := range notificationHeaders[p.Transport] {
		req.Header.Add(header.HeaderName, header.HeaderValue)
	}

	switch p.Transport {
	case TransportEventSub:
		req.Header.Set("Twitch-Eventsub-Message-Id", p.ID)
		req.Header.Set("Twitch-Eventsub-Subscription-Type", p.Event)
	case TransportWebSub:
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
		return 0, err
	}

	defer resp.Body.Close()

	return resp.StatusCode, nil
}

func getSignatureHeader(req *http.Request, id string, secret string, transport string, payload []byte) {
	mac := hmac.New(sha256.New, []byte(secret))
	ts := util.GetTimestamp()
	switch transport {
	case TransportEventSub:
		req.Header.Set("Twitch-Eventsub-Message-Timestamp", ts.Format(time.RFC3339))
		prefix := ts.AppendFormat([]byte(id), time.RFC3339)
		mac.Write(prefix)
		mac.Write(payload)
		req.Header.Set("Twitch-Eventsub-Message-Signature", fmt.Sprintf("sha256=%x", mac.Sum(nil)))
	case TransportWebSub:
		mac.Write(payload)
		req.Header.Set("X-Hub-Signature", fmt.Sprintf("sha256=%x", mac.Sum(nil)))
	}
}
