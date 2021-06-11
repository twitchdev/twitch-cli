// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package endpoints

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/twitchdev/twitch-cli/internal/mock_api/endpoints/bits"
	"github.com/twitchdev/twitch-cli/internal/mock_api/endpoints/categories"
	"github.com/twitchdev/twitch-cli/test_setup"
	"github.com/twitchdev/twitch-cli/test_setup/test_server"
)

// used to share the same db connection and to run them in parallel
func TestCheermotes(t *testing.T) {
	a := test_setup.SetupTestEnv(t)

	ts := test_server.SetupTestServer(bits.Cheermotes{})

	req, _ := http.NewRequest(http.MethodGet, ts.URL+bits.Cheermotes{}.Path(), nil)

	resp, err := http.DefaultClient.Do(req)
	a.Nil(err)
	a.NotNil(resp)
}

func TestBitsLeaderboard(t *testing.T) {
	a := test_setup.SetupTestEnv(t)

	ts := test_server.SetupTestServer(bits.BitsLeaderboard{})

	req, _ := http.NewRequest(http.MethodGet, ts.URL+bits.BitsLeaderboard{}.Path(), nil)
	q := req.URL.Query()
	q.Set("period", "day")
	req.URL.RawQuery = q.Encode()
	resp, err := http.DefaultClient.Do(req)
	a.Nil(err, fmt.Sprint(err))
	a.Equal(http.StatusOK, resp.StatusCode)

	req, _ = http.NewRequest(http.MethodGet, ts.URL+bits.BitsLeaderboard{}.Path(), nil)
	q.Set("period", "week")
	req.URL.RawQuery = q.Encode()
	resp, err = http.DefaultClient.Do(req)
	a.Nil(err, fmt.Sprint(err))
	a.Equal(http.StatusOK, resp.StatusCode)

	req, _ = http.NewRequest(http.MethodGet, ts.URL+bits.BitsLeaderboard{}.Path(), nil)
	q.Set("period", "month")
	req.URL.RawQuery = q.Encode()
	resp, err = http.DefaultClient.Do(req)
	a.Nil(err, fmt.Sprint(err))
	a.Equal(http.StatusOK, resp.StatusCode)

	req, _ = http.NewRequest(http.MethodGet, ts.URL+bits.BitsLeaderboard{}.Path(), nil)
	q.Set("period", "week")
	req.URL.RawQuery = q.Encode()
	resp, err = http.DefaultClient.Do(req)
	a.Nil(err, fmt.Sprint(err))
	a.Equal(http.StatusOK, resp.StatusCode)

	req, _ = http.NewRequest(http.MethodGet, ts.URL+bits.BitsLeaderboard{}.Path(), nil)
	q.Set("period", "potato")
	req.URL.RawQuery = q.Encode()
	resp, err = http.DefaultClient.Do(req)
	a.Nil(err, fmt.Sprint(err))
	a.Equal(http.StatusBadRequest, resp.StatusCode)

	req, _ = http.NewRequest(http.MethodGet, ts.URL+bits.BitsLeaderboard{}.Path(), nil)
	q.Set("count", "1")
	q.Set("period", "")
	req.URL.RawQuery = q.Encode()
	resp, err = http.DefaultClient.Do(req)
	a.Nil(err, fmt.Sprint(err))
	a.Equal(http.StatusOK, resp.StatusCode)

	req, _ = http.NewRequest(http.MethodGet, ts.URL+bits.BitsLeaderboard{}.Path(), nil)
	q.Set("count", "potato")
	req.URL.RawQuery = q.Encode()
	resp, err = http.DefaultClient.Do(req)
	a.Nil(err, fmt.Sprint(err))
	a.Equal(http.StatusBadRequest, resp.StatusCode)

	req, _ = http.NewRequest(http.MethodGet, ts.URL+bits.BitsLeaderboard{}.Path(), nil)
	q.Set("count", "1000")
	req.URL.RawQuery = q.Encode()
	resp, err = http.DefaultClient.Do(req)
	a.Nil(err, fmt.Sprint(err))
	a.Equal(http.StatusBadRequest, resp.StatusCode)

	req, _ = http.NewRequest(http.MethodGet, ts.URL+bits.BitsLeaderboard{}.Path(), nil)
	q.Set("started_at", "2021")
	req.URL.RawQuery = q.Encode()
	resp, err = http.DefaultClient.Do(req)
	a.Nil(err, fmt.Sprint(err))
	a.Equal(http.StatusBadRequest, resp.StatusCode)
}

func TestGames(t *testing.T) {
	a := test_setup.SetupTestEnv(t)

	ts := test_server.SetupTestServer(categories.Games{})

	req, _ := http.NewRequest(http.MethodGet, ts.URL+categories.Games{}.Path(), nil)
	q := req.URL.Query()
	resp, err := http.DefaultClient.Do(req)
	a.Nil(err)
	a.NotNil(resp)

	req, _ = http.NewRequest(http.MethodGet, ts.URL+categories.Games{}.Path(), nil)
	q.Set("id", "1")
	req.URL.RawQuery = q.Encode()
	resp, err = http.DefaultClient.Do(req)
	a.Nil(err)
	a.Equal(http.StatusOK, resp.StatusCode)

	req, _ = http.NewRequest(http.MethodGet, ts.URL+categories.Games{}.Path(), nil)
	q.Set("name", "day")
	req.URL.RawQuery = q.Encode()
	resp, err = http.DefaultClient.Do(req)
	a.Nil(err)
	a.Equal(http.StatusOK, resp.StatusCode)
}

func TestTopGames(t *testing.T) {
	a := test_setup.SetupTestEnv(t)

	ts := test_server.SetupTestServer(categories.TopGames{})

	req, _ := http.NewRequest(http.MethodGet, ts.URL+categories.TopGames{}.Path(), nil)
	q := req.URL.Query()
	q.Set("name", "day")
	req.URL.RawQuery = q.Encode()
	resp, err := http.DefaultClient.Do(req)
	a.Nil(err)
	a.Equal(http.StatusOK, resp.StatusCode)

}
