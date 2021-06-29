// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package streams

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/twitchdev/twitch-cli/test_setup"
	"github.com/twitchdev/twitch-cli/test_setup/test_server"
)

func TestAllTags(t *testing.T) {
	a := test_setup.SetupTestEnv(t)
	ts := test_server.SetupTestServer(AllTags{})

	// get
	req, _ := http.NewRequest(http.MethodGet, ts.URL+AllTags{}.Path(), nil)
	q := req.URL.Query()
	req.URL.RawQuery = q.Encode()
	resp, err := http.DefaultClient.Do(req)
	a.Nil(err)
	a.Equal(200, resp.StatusCode)

	q.Set("tag_id", "1234")
	req.URL.RawQuery = q.Encode()
	resp, err = http.DefaultClient.Do(req)
	a.Nil(err)
	a.Equal(200, resp.StatusCode)
}

func TestFollowedStreams(t *testing.T) {
	a := test_setup.SetupTestEnv(t)
	ts := test_server.SetupTestServer(FollowedStreams{})

	// get
	req, _ := http.NewRequest(http.MethodGet, ts.URL+FollowedStreams{}.Path(), nil)
	q := req.URL.Query()
	req.URL.RawQuery = q.Encode()
	resp, err := http.DefaultClient.Do(req)
	a.Nil(err)
	a.Equal(401, resp.StatusCode)

	q.Set("user_id", "1")
	req.URL.RawQuery = q.Encode()
	resp, err = http.DefaultClient.Do(req)
	a.Nil(err)
	a.Equal(200, resp.StatusCode)
}

func TestMarkers(t *testing.T) {
	a := test_setup.SetupTestEnv(t)
	ts := test_server.SetupTestServer(Markers{})

	// get
	req, _ := http.NewRequest(http.MethodGet, ts.URL+Markers{}.Path(), nil)
	q := req.URL.Query()
	req.URL.RawQuery = q.Encode()
	resp, err := http.DefaultClient.Do(req)
	a.Nil(err)
	a.Equal(400, resp.StatusCode)

	q.Set("user_id", "1")
	req.URL.RawQuery = q.Encode()
	resp, err = http.DefaultClient.Do(req)
	a.Nil(err)
	a.Equal(200, resp.StatusCode)

	q.Set("video_id", "1")
	req.URL.RawQuery = q.Encode()
	resp, err = http.DefaultClient.Do(req)
	a.Nil(err)
	a.Equal(400, resp.StatusCode)

	q.Del("video_id")
	q.Set("user_id", "2")
	req.URL.RawQuery = q.Encode()
	resp, err = http.DefaultClient.Do(req)
	a.Nil(err)
	a.Equal(400, resp.StatusCode)

	//post
	post := MarkerPostBody{
		UserID:      "1",
		Description: "1234",
	}
	b, _ := json.Marshal(post)
	req, _ = http.NewRequest(http.MethodPost, ts.URL+Markers{}.Path(), bytes.NewBuffer(b))
	resp, err = http.DefaultClient.Do(req)
	a.Nil(err)

	post.UserID = "2"
	b, _ = json.Marshal(post)
	req, _ = http.NewRequest(http.MethodPost, ts.URL+Markers{}.Path(), bytes.NewBuffer(b))
	resp, err = http.DefaultClient.Do(req)
	a.Nil(err)
	a.Equal(400, resp.StatusCode)
}

func TestStreamTags(t *testing.T) {
	a := test_setup.SetupTestEnv(t)
	ts := test_server.SetupTestServer(StreamTags{})

	// get
	req, _ := http.NewRequest(http.MethodGet, ts.URL+StreamTags{}.Path(), nil)
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
	put := PutBodyStreamTags{
		TagIDs: []string{"1234"},
	}
	b, _ := json.Marshal(put)
	req, _ = http.NewRequest(http.MethodPut, ts.URL+StreamTags{}.Path(), bytes.NewBuffer(b))
	q.Del("broadcaster_id")
	req.URL.RawQuery = q.Encode()
	resp, err = http.DefaultClient.Do(req)
	a.Nil(err)
	a.Equal(401, resp.StatusCode)

	b, _ = json.Marshal(put)
	req, _ = http.NewRequest(http.MethodPut, ts.URL+StreamTags{}.Path(), bytes.NewBuffer(b))
	q.Set("broadcaster_id", "1")
	req.URL.RawQuery = q.Encode()
	resp, err = http.DefaultClient.Do(req)
	a.Nil(err)
	a.Equal(400, resp.StatusCode)
}

func TestStreamKey(t *testing.T) {
	a := test_setup.SetupTestEnv(t)
	ts := test_server.SetupTestServer(StreamKey{})

	// get
	req, _ := http.NewRequest(http.MethodGet, ts.URL+StreamKey{}.Path(), nil)
	q := req.URL.Query()
	req.URL.RawQuery = q.Encode()
	resp, err := http.DefaultClient.Do(req)
	a.Nil(err)
	a.Equal(401, resp.StatusCode)

	// since the tag ID isn't known, a 400 is expected
	q.Set("broadcaster_id", "1")
	req.URL.RawQuery = q.Encode()
	resp, err = http.DefaultClient.Do(req)
	a.Nil(err)
	a.Equal(200, resp.StatusCode)
}

func TestStreams(t *testing.T) {
	a := test_setup.SetupTestEnv(t)
	ts := test_server.SetupTestServer(Streams{})

	// get
	req, _ := http.NewRequest(http.MethodGet, ts.URL+Streams{}.Path(), nil)
	q := req.URL.Query()
	req.URL.RawQuery = q.Encode()
	resp, err := http.DefaultClient.Do(req)
	a.Nil(err)
	a.Equal(200, resp.StatusCode)

	q.Set("first", "1")
	req.URL.RawQuery = q.Encode()
	resp, err = http.DefaultClient.Do(req)
	a.Nil(err)
	a.Equal(200, resp.StatusCode)

	q.Set("after", "ttt")
	req.URL.RawQuery = q.Encode()
	resp, err = http.DefaultClient.Do(req)
	a.Nil(err)
	a.Equal(200, resp.StatusCode)

	q.Del("after")
	q.Set("before", "ttt")
	req.URL.RawQuery = q.Encode()
	resp, err = http.DefaultClient.Do(req)
	a.Nil(err)
	a.Equal(200, resp.StatusCode)

	q.Del("after")
	q.Del("before")
	q.Set("language", "en")
	q.Set("user_id", "1")
	q.Set("user_login", "1")
	req.URL.RawQuery = q.Encode()
	resp, err = http.DefaultClient.Do(req)
	a.Nil(err)
	a.Equal(200, resp.StatusCode)
}
