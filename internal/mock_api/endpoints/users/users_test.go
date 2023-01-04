// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package users

import (
	"net/http"
	"testing"

	"github.com/twitchdev/twitch-cli/test_setup"
	"github.com/twitchdev/twitch-cli/test_setup/test_server"
)

func TestUsers(t *testing.T) {
	a := test_setup.SetupTestEnv(t)
	ts := test_server.SetupTestServer(UsersEndpoint{})

	// get
	req, _ := http.NewRequest(http.MethodGet, ts.URL+UsersEndpoint{}.Path(), nil)
	q := req.URL.Query()
	req.URL.RawQuery = q.Encode()
	resp, err := http.DefaultClient.Do(req)
	a.Nil(err)
	a.Equal(200, resp.StatusCode)

	q.Set("id", "1")
	q.Set("login", "1")
	q.Add("id", "2")
	req.URL.RawQuery = q.Encode()
	resp, err = http.DefaultClient.Do(req)
	a.Nil(err)
	a.Equal(200, resp.StatusCode)

	// put
	req, _ = http.NewRequest(http.MethodPut, ts.URL+UsersEndpoint{}.Path(), nil)
	q = req.URL.Query()
	req.URL.RawQuery = q.Encode()
	resp, err = http.DefaultClient.Do(req)
	a.Nil(err)
	a.Equal(400, resp.StatusCode)

	q.Set("description", "potato")
	req.URL.RawQuery = q.Encode()
	resp, err = http.DefaultClient.Do(req)
	a.Nil(err)
	a.Equal(200, resp.StatusCode)
}

func TestBlocks(t *testing.T) {
	a := test_setup.SetupTestEnv(t)
	ts := test_server.SetupTestServer(Blocks{})

	// get
	req, _ := http.NewRequest(http.MethodGet, ts.URL+Blocks{}.Path(), nil)
	q := req.URL.Query()
	req.URL.RawQuery = q.Encode()
	resp, err := http.DefaultClient.Do(req)
	a.Nil(err)
	a.Equal(401, resp.StatusCode)

	q.Set("broadcaster_id", "1")
	req.URL.RawQuery = q.Encode()
	resp, err = http.DefaultClient.Do(req)
	a.Nil(err)
	a.Equal(200, resp.StatusCode)

	// put
	req, _ = http.NewRequest(http.MethodPut, ts.URL+Blocks{}.Path(), nil)
	q = req.URL.Query()
	req.URL.RawQuery = q.Encode()
	resp, err = http.DefaultClient.Do(req)
	a.Nil(err)
	a.Equal(400, resp.StatusCode)

	q.Set("target_user_id", "1")
	req.URL.RawQuery = q.Encode()
	resp, err = http.DefaultClient.Do(req)
	a.Nil(err)
	a.Equal(400, resp.StatusCode)

	q.Set("target_user_id", "2")
	q.Set("reason", "other")
	q.Set("source_context", "chat")
	req.URL.RawQuery = q.Encode()
	resp, err = http.DefaultClient.Do(req)
	a.Nil(err)
	a.Equal(204, resp.StatusCode)

	q.Set("reason", "other")
	q.Set("source_context", "potato")
	req.URL.RawQuery = q.Encode()
	resp, err = http.DefaultClient.Do(req)
	a.Nil(err)
	a.Equal(400, resp.StatusCode)

	q.Del("source_context")
	q.Set("reason", "potato")
	req.URL.RawQuery = q.Encode()
	resp, err = http.DefaultClient.Do(req)
	a.Nil(err)
	a.Equal(400, resp.StatusCode)

	// delete
	req, _ = http.NewRequest(http.MethodDelete, ts.URL+Blocks{}.Path(), nil)
	q = req.URL.Query()
	req.URL.RawQuery = q.Encode()
	resp, err = http.DefaultClient.Do(req)
	a.Nil(err)
	a.Equal(400, resp.StatusCode)

	q.Set("target_user_id", "2")
	req.URL.RawQuery = q.Encode()
	resp, err = http.DefaultClient.Do(req)
	a.Nil(err)
	a.Equal(204, resp.StatusCode)
}

func TestFollows(t *testing.T) {
	a := test_setup.SetupTestEnv(t)
	ts := test_server.SetupTestServer(FollowsEndpoint{})

	// get
	req, _ := http.NewRequest(http.MethodGet, ts.URL+FollowsEndpoint{}.Path(), nil)
	q := req.URL.Query()
	req.URL.RawQuery = q.Encode()
	resp, err := http.DefaultClient.Do(req)
	a.Nil(err)
	a.Equal(400, resp.StatusCode)

	q.Set("to_id", "1")
	req.URL.RawQuery = q.Encode()
	resp, err = http.DefaultClient.Do(req)
	a.Nil(err)
	a.Equal(200, resp.StatusCode)
}
