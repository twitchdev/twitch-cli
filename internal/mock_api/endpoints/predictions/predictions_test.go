// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package predictions

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/twitchdev/twitch-cli/test_setup"
	"github.com/twitchdev/twitch-cli/test_setup/test_server"
)

func TestPredictions(t *testing.T) {
	a := test_setup.SetupTestEnv(t)
	ts := test_server.SetupTestServer(Predictions{})

	// get
	req, _ := http.NewRequest(http.MethodGet, ts.URL+Predictions{}.Path(), nil)
	q := req.URL.Query()
	req.URL.RawQuery = q.Encode()
	resp, err := http.DefaultClient.Do(req)
	a.Nil(err)
	a.Equal(401, resp.StatusCode)

	q.Set("broadcaster_id", "1")
	req, _ = http.NewRequest(http.MethodGet, ts.URL+Predictions{}.Path(), nil)
	req.URL.RawQuery = q.Encode()
	resp, err = http.DefaultClient.Do(req)
	a.Nil(err, err)
	a.Equal(200, resp.StatusCode)

	q.Set("id", "1")
	req, _ = http.NewRequest(http.MethodGet, ts.URL+Predictions{}.Path(), nil)
	req.URL.RawQuery = q.Encode()
	resp, err = http.DefaultClient.Do(req)
	a.Nil(err)
	a.Equal(200, resp.StatusCode)

	// post
	post := PostPredictionsBody{
		BroadcasterID:    "1",
		Title:            "2",
		Outcomes:         []PostPredictionsBodyOutcomes{{Title: "3"}, {Title: "4"}},
		PredictionWindow: 100,
	}

	b, _ := json.Marshal(post)
	req, _ = http.NewRequest(http.MethodPost, ts.URL+Predictions{}.Path(), bytes.NewBuffer(b))
	req.URL.RawQuery = q.Encode()
	resp, err = http.DefaultClient.Do(req)
	a.Nil(err, err)
	a.Equal(200, resp.StatusCode)

	post.BroadcasterID = "2"
	b, _ = json.Marshal(post)
	req, _ = http.NewRequest(http.MethodPost, ts.URL+Predictions{}.Path(), bytes.NewBuffer(b))
	req.URL.RawQuery = q.Encode()
	resp, err = http.DefaultClient.Do(req)
	a.Nil(err, err)
	a.Equal(401, resp.StatusCode)

	post.BroadcasterID = "1"
	post.Title = ""
	b, _ = json.Marshal(post)
	req, _ = http.NewRequest(http.MethodPost, ts.URL+Predictions{}.Path(), bytes.NewBuffer(b))
	req.URL.RawQuery = q.Encode()
	resp, err = http.DefaultClient.Do(req)
	a.Nil(err, err)
	a.Equal(400, resp.StatusCode)

	post.Title = "123"
	post.PredictionWindow = 0
	b, _ = json.Marshal(post)
	req, _ = http.NewRequest(http.MethodPost, ts.URL+Predictions{}.Path(), bytes.NewBuffer(b))
	req.URL.RawQuery = q.Encode()
	resp, err = http.DefaultClient.Do(req)
	a.Nil(err, err)
	a.Equal(400, resp.StatusCode)

	post.PredictionWindow = 100
	post.Outcomes = append(post.Outcomes, PostPredictionsBodyOutcomes{Title: "6"})
	b, _ = json.Marshal(post)
	req, _ = http.NewRequest(http.MethodPost, ts.URL+Predictions{}.Path(), bytes.NewBuffer(b))
	req.URL.RawQuery = q.Encode()
	resp, err = http.DefaultClient.Do(req)
	a.Nil(err, err)
	a.Equal(400, resp.StatusCode)

	// patch
	patch := PatchPredictionsBody{
		BroadcasterID:    "1",
		ID:               "test",
		Status:           "RESOLVED",
		WinningOutcomeID: "1234",
	}

	b, _ = json.Marshal(patch)
	req, _ = http.NewRequest(http.MethodPatch, ts.URL+Predictions{}.Path(), bytes.NewBuffer(b))
	req.URL.RawQuery = q.Encode()
	resp, err = http.DefaultClient.Do(req)
	a.Nil(err, err)
	a.Equal(200, resp.StatusCode)

	patch.BroadcasterID = "2"
	b, _ = json.Marshal(patch)
	req, _ = http.NewRequest(http.MethodPatch, ts.URL+Predictions{}.Path(), bytes.NewBuffer(b))
	req.URL.RawQuery = q.Encode()
	resp, err = http.DefaultClient.Do(req)
	a.Nil(err, err)
	a.Equal(401, resp.StatusCode)

	patch.BroadcasterID = "1"
	patch.Status = ""
	b, _ = json.Marshal(patch)
	req, _ = http.NewRequest(http.MethodPatch, ts.URL+Predictions{}.Path(), bytes.NewBuffer(b))
	req.URL.RawQuery = q.Encode()
	resp, err = http.DefaultClient.Do(req)
	a.Nil(err, err)
	a.Equal(400, resp.StatusCode)

	patch.Status = "RESOLVED"
	patch.ID = ""
	b, _ = json.Marshal(patch)
	req, _ = http.NewRequest(http.MethodPatch, ts.URL+Predictions{}.Path(), bytes.NewBuffer(b))
	req.URL.RawQuery = q.Encode()
	resp, err = http.DefaultClient.Do(req)
	a.Nil(err, err)
	a.Equal(400, resp.StatusCode)

	patch.ID = "RESOLVED"
	patch.WinningOutcomeID = ""
	b, _ = json.Marshal(patch)
	req, _ = http.NewRequest(http.MethodPatch, ts.URL+Predictions{}.Path(), bytes.NewBuffer(b))
	req.URL.RawQuery = q.Encode()
	resp, err = http.DefaultClient.Do(req)
	a.Nil(err, err)
	a.Equal(400, resp.StatusCode)
}
