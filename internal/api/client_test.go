// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package api

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/twitchdev/twitch-cli/test_setup"
	"golang.org/x/time/rate"
)

func TestNewClient(t *testing.T) {
	a := test_setup.SetupTestEnv(t)

	var ok = "{\"status\":\"ok\"}"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(ok))

		_, err := io.ReadAll(r.Body)
		a.Nil(err)

	}))

	rl := rate.NewLimiter(rate.Every(time.Minute), 800)
	c := NewClient(rl)

	req, _ := http.NewRequest(http.MethodGet, ts.URL, nil)
	resp, err := c.Do(req)
	a.Nil(err)

	body, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	a.NoError(err)
	a.Equal(ok, string(body), "Body mismatch")

	req, _ = http.NewRequest(http.MethodGet, "potato", nil)
	_, err = c.Do(req)
	a.NotNil(err)
}
