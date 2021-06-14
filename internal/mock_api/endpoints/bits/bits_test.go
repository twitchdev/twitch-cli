// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package bits

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/twitchdev/twitch-cli/test_setup"
	"github.com/twitchdev/twitch-cli/test_setup/test_server"
)

func TestCheermotes(t *testing.T) {
	a := test_setup.SetupTestEnv(t)

	ts := test_server.SetupTestServer(Cheermotes{})

	req, _ := http.NewRequest(http.MethodGet, ts.URL+Cheermotes{}.Path(), nil)

	resp, err := http.DefaultClient.Do(req)
	a.Nil(err)
	a.NotNil(resp)
}

func TestBitsLeaderboard(t *testing.T) {
	a := test_setup.SetupTestEnv(t)

	ts := test_server.SetupTestServer(BitsLeaderboard{})

	req, _ := http.NewRequest(http.MethodGet, ts.URL+BitsLeaderboard{}.Path(), nil)
	q := req.URL.Query()
	q.Set("period", "day")
	req.URL.RawQuery = q.Encode()
	resp, err := http.DefaultClient.Do(req)
	a.Nil(err, fmt.Sprint(err))
	a.Equal(http.StatusOK, resp.StatusCode)

	req, _ = http.NewRequest(http.MethodGet, ts.URL+BitsLeaderboard{}.Path(), nil)
	q.Set("period", "week")
	req.URL.RawQuery = q.Encode()
	resp, err = http.DefaultClient.Do(req)
	a.Nil(err, fmt.Sprint(err))
	a.Equal(http.StatusOK, resp.StatusCode)

	req, _ = http.NewRequest(http.MethodGet, ts.URL+BitsLeaderboard{}.Path(), nil)
	q.Set("period", "month")
	req.URL.RawQuery = q.Encode()
	resp, err = http.DefaultClient.Do(req)
	a.Nil(err, fmt.Sprint(err))
	a.Equal(http.StatusOK, resp.StatusCode)

	req, _ = http.NewRequest(http.MethodGet, ts.URL+BitsLeaderboard{}.Path(), nil)
	q.Set("period", "week")
	req.URL.RawQuery = q.Encode()
	resp, err = http.DefaultClient.Do(req)
	a.Nil(err, fmt.Sprint(err))
	a.Equal(http.StatusOK, resp.StatusCode)

	req, _ = http.NewRequest(http.MethodGet, ts.URL+BitsLeaderboard{}.Path(), nil)
	q.Set("period", "potato")
	req.URL.RawQuery = q.Encode()
	resp, err = http.DefaultClient.Do(req)
	a.Nil(err, fmt.Sprint(err))
	a.Equal(http.StatusBadRequest, resp.StatusCode)

	req, _ = http.NewRequest(http.MethodGet, ts.URL+BitsLeaderboard{}.Path(), nil)
	q.Set("count", "1")
	q.Set("period", "")
	req.URL.RawQuery = q.Encode()
	resp, err = http.DefaultClient.Do(req)
	a.Nil(err, fmt.Sprint(err))
	a.Equal(http.StatusOK, resp.StatusCode)

	req, _ = http.NewRequest(http.MethodGet, ts.URL+BitsLeaderboard{}.Path(), nil)
	q.Set("count", "potato")
	req.URL.RawQuery = q.Encode()
	resp, err = http.DefaultClient.Do(req)
	a.Nil(err, fmt.Sprint(err))
	a.Equal(http.StatusBadRequest, resp.StatusCode)

	req, _ = http.NewRequest(http.MethodGet, ts.URL+BitsLeaderboard{}.Path(), nil)
	q.Set("count", "1000")
	req.URL.RawQuery = q.Encode()
	resp, err = http.DefaultClient.Do(req)
	a.Nil(err, fmt.Sprint(err))
	a.Equal(http.StatusBadRequest, resp.StatusCode)

	req, _ = http.NewRequest(http.MethodGet, ts.URL+BitsLeaderboard{}.Path(), nil)
	q.Set("started_at", "2021")
	req.URL.RawQuery = q.Encode()
	resp, err = http.DefaultClient.Do(req)
	a.Nil(err, fmt.Sprint(err))
	a.Equal(http.StatusBadRequest, resp.StatusCode)
}
