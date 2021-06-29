// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package categories

import (
	"net/http"
	"testing"

	"github.com/twitchdev/twitch-cli/test_setup"
	"github.com/twitchdev/twitch-cli/test_setup/test_server"
)

func TestGames(t *testing.T) {
	a := test_setup.SetupTestEnv(t)

	ts := test_server.SetupTestServer(Games{})

	req, _ := http.NewRequest(http.MethodGet, ts.URL+Games{}.Path(), nil)
	q := req.URL.Query()
	resp, err := http.DefaultClient.Do(req)
	a.Nil(err)
	a.NotNil(resp)

	req, _ = http.NewRequest(http.MethodGet, ts.URL+Games{}.Path(), nil)
	q.Set("id", "1")
	req.URL.RawQuery = q.Encode()
	resp, err = http.DefaultClient.Do(req)
	a.Nil(err)
	a.Equal(http.StatusOK, resp.StatusCode)

	req, _ = http.NewRequest(http.MethodGet, ts.URL+Games{}.Path(), nil)
	q.Set("name", "day")
	req.URL.RawQuery = q.Encode()
	resp, err = http.DefaultClient.Do(req)
	a.Nil(err)
	a.Equal(http.StatusOK, resp.StatusCode)
}

func TestTopGames(t *testing.T) {
	a := test_setup.SetupTestEnv(t)

	ts := test_server.SetupTestServer(TopGames{})

	req, _ := http.NewRequest(http.MethodGet, ts.URL+TopGames{}.Path(), nil)
	q := req.URL.Query()
	q.Set("name", "day")
	req.URL.RawQuery = q.Encode()
	resp, err := http.DefaultClient.Do(req)
	a.Nil(err)
	a.Equal(http.StatusOK, resp.StatusCode)

}
