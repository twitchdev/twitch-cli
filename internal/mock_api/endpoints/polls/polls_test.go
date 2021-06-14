// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package polls

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/twitchdev/twitch-cli/test_setup"
	"github.com/twitchdev/twitch-cli/test_setup/test_server"
)

func TestPolls(t *testing.T) {
	a := test_setup.SetupTestEnv(t)
	ts := test_server.SetupTestServer(Polls{})

	// get
	req, _ := http.NewRequest(http.MethodGet, ts.URL+Polls{}.Path(), nil)
	q := req.URL.Query()
	req.URL.RawQuery = q.Encode()
	resp, err := http.DefaultClient.Do(req)
	a.Nil(err)
	a.Equal(401, resp.StatusCode)

	q.Set("broadcaster_id", "1")
	req, _ = http.NewRequest(http.MethodGet, ts.URL+Polls{}.Path(), nil)
	req.URL.RawQuery = q.Encode()
	resp, err = http.DefaultClient.Do(req)
	a.Nil(err, err)
	a.Equal(200, resp.StatusCode)

	q.Set("id", "1")
	req, _ = http.NewRequest(http.MethodGet, ts.URL+Polls{}.Path(), nil)
	req.URL.RawQuery = q.Encode()
	resp, err = http.DefaultClient.Do(req)
	a.Nil(err)
	a.Equal(200, resp.StatusCode)

	// post
	body := PostPollsBody{
		BroadcasterID: "1",
		Title:         "2",
		Choices: []PostPollsBodyChoice{
			{Title: "3"},
			{Title: "4"},
		},
		BitsVotingEnabled:          false,
		ChannelPointsVotingEnabled: false,
		Duration:                   300,
	}

	b, _ := json.Marshal(body)
	req, _ = http.NewRequest(http.MethodPost, ts.URL+Polls{}.Path(), bytes.NewBuffer(b))
	req.URL.RawQuery = q.Encode()
	resp, err = http.DefaultClient.Do(req)
	a.Nil(err, err)
	a.Equal(200, resp.StatusCode)

	body.BroadcasterID = "2"
	b, _ = json.Marshal(body)
	req, _ = http.NewRequest(http.MethodPost, ts.URL+Polls{}.Path(), bytes.NewBuffer(b))
	req.URL.RawQuery = q.Encode()
	resp, err = http.DefaultClient.Do(req)
	a.Nil(err, err)
	a.Equal(401, resp.StatusCode)

	body.BroadcasterID = "1"
	body.Title = ""
	b, _ = json.Marshal(body)
	req, _ = http.NewRequest(http.MethodPost, ts.URL+Polls{}.Path(), bytes.NewBuffer(b))
	req.URL.RawQuery = q.Encode()
	resp, err = http.DefaultClient.Do(req)
	a.Nil(err, err)
	a.Equal(400, resp.StatusCode)

	body.Title = "2"
	body.Choices = body.Choices[:1]
	b, _ = json.Marshal(body)
	req, _ = http.NewRequest(http.MethodPost, ts.URL+Polls{}.Path(), bytes.NewBuffer(b))
	req.URL.RawQuery = q.Encode()
	resp, err = http.DefaultClient.Do(req)
	a.Nil(err, err)
	a.Equal(400, resp.StatusCode)

	body.Duration = 14
	body.Choices = append(body.Choices, PostPollsBodyChoice{Title: "4"})
	b, _ = json.Marshal(body)
	req, _ = http.NewRequest(http.MethodPost, ts.URL+Polls{}.Path(), bytes.NewBuffer(b))
	req.URL.RawQuery = q.Encode()
	resp, err = http.DefaultClient.Do(req)
	a.Nil(err, err)
	a.Equal(400, resp.StatusCode)

	//patch
	body2 := PatchPollsBody{
		BroadcasterID: "1",
		ID:            "testID",
		Status:        "ARCHIVED",
	}

	b, _ = json.Marshal(body2)
	req, _ = http.NewRequest(http.MethodPatch, ts.URL+Polls{}.Path(), bytes.NewBuffer(b))
	req.URL.RawQuery = q.Encode()
	resp, err = http.DefaultClient.Do(req)
	a.Nil(err, err)
	a.Equal(200, resp.StatusCode)

	body2.BroadcasterID = "2"
	b, _ = json.Marshal(body2)
	req, _ = http.NewRequest(http.MethodPatch, ts.URL+Polls{}.Path(), bytes.NewBuffer(b))
	req.URL.RawQuery = q.Encode()
	resp, err = http.DefaultClient.Do(req)
	a.Nil(err, err)
	a.Equal(401, resp.StatusCode)

	body2.BroadcasterID = "1"
	body2.ID = ""
	b, _ = json.Marshal(body2)
	req, _ = http.NewRequest(http.MethodPatch, ts.URL+Polls{}.Path(), bytes.NewBuffer(b))
	req.URL.RawQuery = q.Encode()
	resp, err = http.DefaultClient.Do(req)
	a.Nil(err, err)
	a.Equal(400, resp.StatusCode)

	body2.ID = "testingID"
	body2.Status = "potato"
	b, _ = json.Marshal(body2)
	req, _ = http.NewRequest(http.MethodPatch, ts.URL+Polls{}.Path(), bytes.NewBuffer(b))
	req.URL.RawQuery = q.Encode()
	resp, err = http.DefaultClient.Do(req)
	a.Nil(err, err)
	a.Equal(400, resp.StatusCode)
}
