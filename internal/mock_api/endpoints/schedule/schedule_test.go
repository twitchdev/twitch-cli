// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package schedule

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/twitchdev/twitch-cli/internal/database"
	"github.com/twitchdev/twitch-cli/internal/util"
	"github.com/twitchdev/twitch-cli/test_setup"
	"github.com/twitchdev/twitch-cli/test_setup/test_server"
)

type RewardResponse struct {
	Data []database.ChannelPointsReward `json:"data"`
}

var (
	segment database.ScheduleSegment
)

func TestMain(m *testing.M) {
	test_setup.SetupTestEnv(&testing.T{})

	db, err := database.NewConnection(true)
	if err != nil {
		log.Fatal(err)
	}
	f := false
	s := database.ScheduleSegment{
		ID:          util.RandomGUID(),
		UserID:      "1",
		Title:       "from_unit_tests",
		IsRecurring: true,
		IsVacation:  false,
		StartTime:   time.Now().UTC().Format(time.RFC3339),
		EndTime:     time.Now().UTC().Add(24 * time.Hour).Format(time.RFC3339),
		IsCanceled:  &f,
	}
	err = db.NewQuery(nil, 100).InsertSchedule(s)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.NewQuery(nil, 100).GetSchedule(database.ScheduleSegment{UserID: "1"}, time.Date(0, 0, 0, 0, 0, 0, 0, time.UTC))

	segment = s
	db.DB.Close()

	os.Exit(m.Run())
}
func TestSchedule(t *testing.T) {
	a := test_setup.SetupTestEnv(t)

	ts := test_server.SetupTestServer(Schedule{})

	// get
	req, _ := http.NewRequest(http.MethodGet, ts.URL+Schedule{}.Path(), nil)
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

	q.Set("broadcaster_id", "2")
	req.URL.RawQuery = q.Encode()
	resp, err = http.DefaultClient.Do(req)
	a.Nil(err)
	a.Equal(200, resp.StatusCode)

	q.Set("broadcaster_id", "1")
	q.Set("id", segment.ID)
	req.URL.RawQuery = q.Encode()
	resp, err = http.DefaultClient.Do(req)
	a.Nil(err)
	a.Equal(200, resp.StatusCode)

	q.Set("utc_offset", "60")
	req.URL.RawQuery = q.Encode()
	resp, err = http.DefaultClient.Do(req)
	a.Nil(err)
	a.Equal(200, resp.StatusCode)

	q.Set("start_time", "test")
	req.URL.RawQuery = q.Encode()
	resp, err = http.DefaultClient.Do(req)
	a.Nil(err)
	a.Equal(400, resp.StatusCode)

	q.Set("start_time", segment.StartTime)
	req.URL.RawQuery = q.Encode()
	resp, err = http.DefaultClient.Do(req)
	a.Nil(err)
	a.Equal(200, resp.StatusCode)
}

func TestICal(t *testing.T) {
	a := test_setup.SetupTestEnv(t)

	ts := test_server.SetupTestServer(ScheduleICal{})
	req, _ := http.NewRequest(http.MethodGet, ts.URL+Schedule{}.Path(), nil)
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
}

func TestSegment(t *testing.T) {
	a := test_setup.SetupTestEnv(t)

	ts := test_server.SetupTestServer(ScheduleSegment{})
	tr := true

	// post tests
	body := SegmentPatchAndPostBody{
		Title:       "hello",
		Timezone:    "America/Los_Angeles",
		StartTime:   time.Now().Format(time.RFC3339),
		IsRecurring: &tr,
		Duration:    "60",
	}

	b, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodPost, ts.URL+ScheduleSegment{}.Path(), bytes.NewBuffer(b))
	q := req.URL.Query()
	q.Set("broadcaster_id", "1")
	req.URL.RawQuery = q.Encode()
	resp, err := http.DefaultClient.Do(req)
	a.Nil(err)
	a.NotNil(resp)
	a.Equal(200, resp.StatusCode)

	body.Title = ""
	b, _ = json.Marshal(body)
	req, _ = http.NewRequest(http.MethodPost, ts.URL+ScheduleSegment{}.Path(), bytes.NewBuffer(b))
	q.Set("broadcaster_id", "1")
	req.URL.RawQuery = q.Encode()
	resp, err = http.DefaultClient.Do(req)
	a.Nil(err)
	a.Equal(200, resp.StatusCode)

	b, _ = json.Marshal(body)
	req, _ = http.NewRequest(http.MethodPost, ts.URL+ScheduleSegment{}.Path(), bytes.NewBuffer(b))
	q.Set("broadcaster_id", "2")
	req.URL.RawQuery = q.Encode()
	resp, err = http.DefaultClient.Do(req)
	a.Nil(err)
	a.Equal(401, resp.StatusCode)

	body.Title = "testing"
	body.Timezone = ""
	b, _ = json.Marshal(body)
	req, _ = http.NewRequest(http.MethodPost, ts.URL+ScheduleSegment{}.Path(), bytes.NewBuffer(b))
	q.Set("broadcaster_id", "1")
	req.URL.RawQuery = q.Encode()
	resp, err = http.DefaultClient.Do(req)
	a.Nil(err)
	a.Equal(400, resp.StatusCode)

	body.Timezone = "test"
	b, _ = json.Marshal(body)
	req, _ = http.NewRequest(http.MethodPost, ts.URL+ScheduleSegment{}.Path(), bytes.NewBuffer(b))
	q.Set("broadcaster_id", "1")
	req.URL.RawQuery = q.Encode()
	resp, err = http.DefaultClient.Do(req)
	a.Nil(err)
	a.Equal(400, resp.StatusCode)

	// patch
	// no id
	b, _ = json.Marshal(body)
	req, _ = http.NewRequest(http.MethodPatch, ts.URL+ScheduleSegment{}.Path(), bytes.NewBuffer(b))
	q.Set("broadcaster_id", "1")
	q.Del("id")
	req.URL.RawQuery = q.Encode()
	resp, err = http.DefaultClient.Do(req)
	a.Nil(err)
	a.Equal(400, resp.StatusCode)

	//mismatch bid and token
	b, _ = json.Marshal(body)
	req, _ = http.NewRequest(http.MethodPatch, ts.URL+ScheduleSegment{}.Path(), bytes.NewBuffer(b))
	q.Set("broadcaster_id", "2")
	req.URL.RawQuery = q.Encode()
	resp, err = http.DefaultClient.Do(req)
	a.Nil(err)
	a.Equal(401, resp.StatusCode)

	// good request
	body.Title = "patched_title"
	body.Timezone = "America/Los_Angeles"
	b, _ = json.Marshal(body)
	req, _ = http.NewRequest(http.MethodPatch, ts.URL+ScheduleSegment{}.Path(), bytes.NewBuffer(b))
	q.Set("broadcaster_id", "1")
	q.Set("id", segment.ID)
	req.URL.RawQuery = q.Encode()
	resp, err = http.DefaultClient.Do(req)
	a.Nil(err)
	a.Equal(200, resp.StatusCode)

	// delete
	req, _ = http.NewRequest(http.MethodDelete, ts.URL+ScheduleSegment{}.Path(), nil)
	q.Set("broadcaster_id", "1")
	req.URL.RawQuery = q.Encode()
	resp, err = http.DefaultClient.Do(req)
	a.Nil(err)
	a.Equal(204, resp.StatusCode)

	q.Set("id", segment.ID)
	req.URL.RawQuery = q.Encode()
	resp, err = http.DefaultClient.Do(req)
	a.Nil(err)
	a.Equal(204, resp.StatusCode)

	q.Set("broadcaster_id", "2")
	req.URL.RawQuery = q.Encode()
	resp, err = http.DefaultClient.Do(req)
	a.Nil(err)
	a.Equal(401, resp.StatusCode)
}

func TestSettings(t *testing.T) {
	a := test_setup.SetupTestEnv(t)

	ts := test_server.SetupTestServer(ScheduleSettings{})
	tr := true
	f := false

	// patch tests
	body := PatchSettingsBody{
		Timezone:          "America/Los_Angeles",
		VacationStartTime: time.Now().Format(time.RFC3339),
		VacationEndTime:   segment.EndTime,
		IsVacationEnabled: &f,
	}

	b, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodPatch, ts.URL+ScheduleSettings{}.Path(), bytes.NewBuffer(b))
	q := req.URL.Query()
	req.URL.RawQuery = q.Encode()
	resp, err := http.DefaultClient.Do(req)
	a.Nil(err)
	a.NotNil(resp)
	a.Equal(401, resp.StatusCode)

	b, _ = json.Marshal(body)
	req, _ = http.NewRequest(http.MethodPatch, ts.URL+ScheduleSettings{}.Path(), bytes.NewBuffer(b))
	q.Set("broadcaster_id", "1")
	req.URL.RawQuery = q.Encode()
	resp, err = http.DefaultClient.Do(req)
	a.Nil(err)
	a.Equal(204, resp.StatusCode)

	body.IsVacationEnabled = &tr
	b, _ = json.Marshal(body)
	req, _ = http.NewRequest(http.MethodPatch, ts.URL+ScheduleSettings{}.Path(), bytes.NewBuffer(b))
	req.URL.RawQuery = q.Encode()
	resp, err = http.DefaultClient.Do(req)
	a.Nil(err)
	a.Equal(204, resp.StatusCode)

	b, _ = json.Marshal(body)
	req, _ = http.NewRequest(http.MethodPatch, ts.URL+ScheduleSettings{}.Path(), bytes.NewBuffer(b))
	q.Set("broadcaster_id", "1")
	req.URL.RawQuery = q.Encode()
	resp, err = http.DefaultClient.Do(req)
	a.Nil(err)
	a.Equal(400, resp.StatusCode)

	body.IsVacationEnabled = &f
	b, _ = json.Marshal(body)
	req, _ = http.NewRequest(http.MethodPatch, ts.URL+ScheduleSettings{}.Path(), bytes.NewBuffer(b))
	req.URL.RawQuery = q.Encode()
	resp, err = http.DefaultClient.Do(req)
	a.Nil(err)
	a.Equal(204, resp.StatusCode)

	body.IsVacationEnabled = &tr
	body.VacationStartTime = "123"
	b, _ = json.Marshal(body)
	req, _ = http.NewRequest(http.MethodPatch, ts.URL+ScheduleSettings{}.Path(), bytes.NewBuffer(b))
	req.URL.RawQuery = q.Encode()
	resp, err = http.DefaultClient.Do(req)
	a.Nil(err)
	a.Equal(400, resp.StatusCode)

	body.VacationStartTime = segment.StartTime
	body.VacationEndTime = "123"
	b, _ = json.Marshal(body)
	req, _ = http.NewRequest(http.MethodPatch, ts.URL+ScheduleSettings{}.Path(), bytes.NewBuffer(b))
	req.URL.RawQuery = q.Encode()
	resp, err = http.DefaultClient.Do(req)
	a.Nil(err)
	a.Equal(400, resp.StatusCode)

	body.VacationEndTime = segment.EndTime
	body.Timezone = "1"
	b, _ = json.Marshal(body)
	req, _ = http.NewRequest(http.MethodPatch, ts.URL+ScheduleSettings{}.Path(), bytes.NewBuffer(b))
	req.URL.RawQuery = q.Encode()
	resp, err = http.DefaultClient.Do(req)
	a.Nil(err)
	a.Equal(400, resp.StatusCode)
}
