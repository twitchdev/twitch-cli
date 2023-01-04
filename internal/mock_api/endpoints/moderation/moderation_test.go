// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package moderation

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/twitchdev/twitch-cli/test_setup"
	"github.com/twitchdev/twitch-cli/test_setup/test_server"
)

func TestModerators(t *testing.T) {
	a := test_setup.SetupTestEnv(t)
	ts := test_server.SetupTestServer(Moderators{})

	req, _ := http.NewRequest(http.MethodGet, ts.URL+Moderators{}.Path(), nil)
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

	q.Set("user_id", "2")
	req.URL.RawQuery = q.Encode()
	resp, err = http.DefaultClient.Do(req)
	a.Nil(err)
	a.Equal(200, resp.StatusCode)
}

func TestAutoModHeld(t *testing.T) {
	a := test_setup.SetupTestEnv(t)
	ts := test_server.SetupTestServer(AutomodHeld{})
	// post
	body := PostAutomodHeldBody{
		UserID:    "1",
		MessageID: "123",
		Action:    "ALLOW",
	}

	b, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodPost, ts.URL+AutomodHeld{}.Path(), bytes.NewBuffer(b))
	q := req.URL.Query()
	req.URL.RawQuery = q.Encode()
	resp, err := http.DefaultClient.Do(req)
	a.Nil(err)
	a.Equal(204, resp.StatusCode)

	body.UserID = "2"
	b, _ = json.Marshal(body)
	req, _ = http.NewRequest(http.MethodPost, ts.URL+AutomodHeld{}.Path(), bytes.NewBuffer(b))
	req.URL.RawQuery = q.Encode()
	resp, err = http.DefaultClient.Do(req)
	a.Nil(err)
	a.Equal(401, resp.StatusCode)

	body.UserID = "1"
	body.MessageID = ""
	b, _ = json.Marshal(body)
	req, _ = http.NewRequest(http.MethodPost, ts.URL+AutomodHeld{}.Path(), bytes.NewBuffer(b))
	req.URL.RawQuery = q.Encode()
	resp, err = http.DefaultClient.Do(req)
	a.Nil(err)
	a.Equal(400, resp.StatusCode)

	body.MessageID = "1234"
	body.Action = "POTATO"
	b, _ = json.Marshal(body)
	req, _ = http.NewRequest(http.MethodPost, ts.URL+AutomodHeld{}.Path(), bytes.NewBuffer(b))
	req.URL.RawQuery = q.Encode()
	resp, err = http.DefaultClient.Do(req)
	a.Nil(err)
	a.Equal(400, resp.StatusCode)
}

func TestAutoModStatus(t *testing.T) {
	a := test_setup.SetupTestEnv(t)
	ts := test_server.SetupTestServer(AutomodStatus{})
	// post
	body := PostAutomodStatusBody{
		Data: []PostAutomodStatusBodyData{{
			UserID:    "1",
			MessageID: "123",
		}},
	}

	b, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodPost, ts.URL+AutomodStatus{}.Path(), bytes.NewBuffer(b))
	q := req.URL.Query()
	req.URL.RawQuery = q.Encode()
	resp, err := http.DefaultClient.Do(req)
	a.Nil(err)
	a.Equal(401, resp.StatusCode)

	b, _ = json.Marshal(body)
	req, _ = http.NewRequest(http.MethodPost, ts.URL+AutomodHeld{}.Path(), bytes.NewBuffer(b))
	q.Set("broadcaster_id", "1")
	req.URL.RawQuery = q.Encode()
	resp, err = http.DefaultClient.Do(req)
	a.Nil(err)
	a.Equal(400, resp.StatusCode)

	body.Data[0].MessageText = "ALLOW"
	b, _ = json.Marshal(body)
	req, _ = http.NewRequest(http.MethodPost, ts.URL+AutomodHeld{}.Path(), bytes.NewBuffer(b))
	q.Set("broadcaster_id", "1")
	req.URL.RawQuery = q.Encode()
	resp, err = http.DefaultClient.Do(req)
	a.Nil(err)
	a.Equal(200, resp.StatusCode)
}

func TestBanned(t *testing.T) {
	a := test_setup.SetupTestEnv(t)
	ts := test_server.SetupTestServer(Bans{})

	req, _ := http.NewRequest(http.MethodGet, ts.URL+Bans{}.Path(), nil)
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

	q.Set("user_id", "2")
	req.URL.RawQuery = q.Encode()
	resp, err = http.DefaultClient.Do(req)
	a.Nil(err)
	a.Equal(200, resp.StatusCode)
}
