// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package channels

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"testing"

	"github.com/twitchdev/twitch-cli/internal/database"
	"github.com/twitchdev/twitch-cli/internal/util"
	"github.com/twitchdev/twitch-cli/test_setup"
	"github.com/twitchdev/twitch-cli/test_setup/test_server"
)

func TestMain(m *testing.M) {
	test_setup.SetupTestEnv(&testing.T{})

	// adding mock data
	db, _ := database.NewConnection(true)
	q := db.NewQuery(nil, 100)
	q.InsertStream(database.Stream{ID: util.RandomGUID(), UserID: "1", StreamType: "live", ViewerCount: 0}, false)
	db.DB.Close()

	os.Exit(m.Run())
}
func TestCommercial(t *testing.T) {
	a := test_setup.SetupTestEnv(t)

	ts := test_server.SetupTestServer(CommercialEndpoint{})

	thirty := 30
	body := CommercialEndpointRequest{
		Length:        &thirty,
		BroadcasterID: "1",
	}

	b, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodPost, ts.URL+CommercialEndpoint{}.Path(), bytes.NewBuffer(b))
	resp, err := http.DefaultClient.Do(req)
	a.Nil(err)
	a.Equal(200, resp.StatusCode)

	ten := 10
	body.Length = &ten
	b, _ = json.Marshal(body)
	req, _ = http.NewRequest(http.MethodPost, ts.URL+CommercialEndpoint{}.Path(), bytes.NewBuffer(b))
	resp, err = http.DefaultClient.Do(req)
	a.Nil(err)
	a.Equal(400, resp.StatusCode)

	body.Length = &thirty
	body.BroadcasterID = "2"
	b, _ = json.Marshal(body)
	req, _ = http.NewRequest(http.MethodPost, ts.URL+CommercialEndpoint{}.Path(), bytes.NewBuffer(b))
	resp, err = http.DefaultClient.Do(req)
	a.Nil(err)
	a.Equal(401, resp.StatusCode)

	body.BroadcasterID = ""
	b, _ = json.Marshal(body)
	req, _ = http.NewRequest(http.MethodPost, ts.URL+CommercialEndpoint{}.Path(), bytes.NewBuffer(b))
	resp, err = http.DefaultClient.Do(req)
	a.Nil(err)
	a.Equal(400, resp.StatusCode)
}

func TestEditors(t *testing.T) {
	a := test_setup.SetupTestEnv(t)
	ts := test_server.SetupTestServer(Editors{})

	req, _ := http.NewRequest(http.MethodGet, ts.URL+Editors{}.Path(), nil)
	q := req.URL.Query()
	req.URL.RawQuery = q.Encode()
	resp, err := http.DefaultClient.Do(req)
	a.Nil(err)
	a.Equal(400, resp.StatusCode)

	q.Set("broadcaster_id", "2")
	req.URL.RawQuery = q.Encode()
	resp, err = http.DefaultClient.Do(req)
	a.Nil(err)
	a.Equal(401, resp.StatusCode)

	q.Set("broadcaster_id", "1")
	req.URL.RawQuery = q.Encode()
	resp, err = http.DefaultClient.Do(req)
	a.Nil(err)
	a.Equal(200, resp.StatusCode)
}

func TestVIPs(t *testing.T) {
	a := test_setup.SetupTestEnv(t)
	ts := test_server.SetupTestServer(Vips{})

	// get
	req, _ := http.NewRequest(http.MethodGet, ts.URL+Vips{}.Path(), nil)
	q := req.URL.Query()
	req.URL.RawQuery = q.Encode()
	resp, err := http.DefaultClient.Do(req)
	a.Nil(err)
	a.Equal(400, resp.StatusCode)

	q.Set("broadcaster_id", "1")
	q.Set("user_id", "3")
	req.URL.RawQuery = q.Encode()
	resp, err = http.DefaultClient.Do(req)
	a.Nil(err)
	a.Equal(200, resp.StatusCode)

	// post
	req, _ = http.NewRequest(http.MethodPost, ts.URL+Vips{}.Path(), nil)
	q.Set("broadcaster_id", "1")
	q.Set("user_id", "99")
	req.URL.RawQuery = q.Encode()
	resp, err = http.DefaultClient.Do(req)
	a.Nil(err)
	a.Equal(204, resp.StatusCode)

	req, _ = http.NewRequest(http.MethodPost, ts.URL+Vips{}.Path(), nil)
	q.Set("user_id", "-1")
	req.URL.RawQuery = q.Encode()
	resp, err = http.DefaultClient.Do(req)
	a.Nil(err)
	a.Equal(404, resp.StatusCode)

	req, _ = http.NewRequest(http.MethodPost, ts.URL+Vips{}.Path(), nil)
	q.Del("user_id")
	req.URL.RawQuery = q.Encode()
	resp, err = http.DefaultClient.Do(req)
	a.Nil(err)
	a.Equal(400, resp.StatusCode)
}

func TestInformation(t *testing.T) {
	a := test_setup.SetupTestEnv(t)
	ts := test_server.SetupTestServer(InformationEndpoint{})

	// get
	req, _ := http.NewRequest(http.MethodGet, ts.URL+InformationEndpoint{}.Path(), nil)
	q := req.URL.Query()
	req.URL.RawQuery = q.Encode()
	resp, err := http.DefaultClient.Do(req)
	a.Nil(err)
	a.Equal(400, resp.StatusCode)

	q.Set("broadcaster_id", "1")
	req.URL.RawQuery = q.Encode()
	resp, err = http.DefaultClient.Do(req)
	a.Nil(err)
	a.Equal(200, resp.StatusCode)

	// patch
	body := PatchInformationEndpointRequest{
		GameID:              "",
		BroadcasterLanguage: "en",
		Title:               "1234",
	}

	b, _ := json.Marshal(body)
	req, _ = http.NewRequest(http.MethodPatch, ts.URL+InformationEndpoint{}.Path(), bytes.NewBuffer(b))
	q.Del("broadcaster_id")
	req.URL.RawQuery = q.Encode()
	resp, err = http.DefaultClient.Do(req)
	a.Nil(err)
	a.Equal(400, resp.StatusCode)

	req, _ = http.NewRequest(http.MethodPatch, ts.URL+InformationEndpoint{}.Path(), bytes.NewBuffer(b))
	q.Set("broadcaster_id", "2")
	req.URL.RawQuery = q.Encode()
	resp, err = http.DefaultClient.Do(req)
	a.Nil(err, resp.StatusCode)
	a.Equal(401, resp.StatusCode)

	req, _ = http.NewRequest(http.MethodPatch, ts.URL+InformationEndpoint{}.Path(), bytes.NewBuffer(b))
	q.Set("broadcaster_id", "1")
	req.URL.RawQuery = q.Encode()
	resp, err = http.DefaultClient.Do(req)
	a.Nil(err, err)
	a.Equal(204, resp.StatusCode)

	body.GameID = "not a real gameid"
	b, _ = json.Marshal(body)
	req, _ = http.NewRequest(http.MethodPatch, ts.URL+InformationEndpoint{}.Path(), bytes.NewBuffer(b))
	req.URL.RawQuery = q.Encode()
	resp, err = http.DefaultClient.Do(req)
	a.Nil(err)
	a.Equal(400, resp.StatusCode)
}

func TestFollowed(t *testing.T) {
	a := test_setup.SetupTestEnv(t)
	ts := test_server.SetupTestServer(FollowedEndpoint{})

	// get
	req, _ := http.NewRequest(http.MethodGet, ts.URL+FollowedEndpoint{}.Path(), nil)
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

	q.Set("broadcaster_id", "2")
	req.URL.RawQuery = q.Encode()
	resp, err = http.DefaultClient.Do(req)
	a.Nil(err)
	a.Equal(200, resp.StatusCode)
}

func TestFollowers(t *testing.T) {
	a := test_setup.SetupTestEnv(t)
	ts := test_server.SetupTestServer(FollowersEndpoint{})

	// get
	req, _ := http.NewRequest(http.MethodGet, ts.URL+FollowersEndpoint{}.Path(), nil)
	q := req.URL.Query()
	req.URL.RawQuery = q.Encode()
	resp, err := http.DefaultClient.Do(req)
	a.Nil(err)
	a.Equal(400, resp.StatusCode)

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
