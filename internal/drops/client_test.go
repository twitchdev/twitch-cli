// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package drops

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/twitchdev/twitch-cli/internal/util"
	"golang.org/x/time/rate"
)

func TestNewClient(t *testing.T) {
	a := util.SetupTestEnv(t)

	var ok = "{\"status\":\"ok\"}"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(ok))

		_, err := ioutil.ReadAll(r.Body)
		a.Nil(err)

	}))

	rl := rate.NewLimiter(rate.Every(10*time.Second), 50)
	c := NewClient(rl)

	req, _ := http.NewRequest(http.MethodGet, ts.URL, nil)
	resp, err := c.Do(req)
	a.Nil(err)

	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	a.Equal(ok, string(body), "Body mismatch")

	req, _ = http.NewRequest(http.MethodGet, "potato", nil)
	resp, err = c.Do(req)
	a.NotNil(err)
}
