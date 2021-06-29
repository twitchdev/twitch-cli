// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package channel_points

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
	rewardID     string
	redemptionID string
)

func TestMain(m *testing.M) {
	test_setup.SetupTestEnv(&testing.T{})

	db, err := database.NewConnection()
	if err != nil {
		log.Fatal(err)
	}
	bTrue := true
	reward := database.ChannelPointsReward{
		ID:                         util.RandomGUID(),
		BroadcasterID:              "1",
		BackgroundColor:            "#fff",
		IsEnabled:                  &bTrue,
		Cost:                       100,
		Title:                      "from_unit_tests",
		RewardPrompt:               "",
		IsUserInputRequired:        false,
		IsPaused:                   false,
		IsInStock:                  false,
		ShouldRedemptionsSkipQueue: false,
	}
	err = db.NewQuery(nil, 100).InsertChannelPointsReward(reward)
	if err != nil {
		log.Fatal(err)
	}

	redemption := database.ChannelPointsRedemption{
		ID:               util.RandomGUID(),
		BroadcasterID:    "1",
		UserID:           "2",
		RedemptionStatus: "UNFULFILLED",
		RewardID:         reward.ID,
		RedeemedAt:       util.GetTimestamp().Format(time.RFC3339),
	}
	err = db.NewQuery(nil, 100).InsertChannelPointsRedemption(redemption)
	dbr, err := db.NewQuery(nil, 100).GetChannelPointsRedemption(database.ChannelPointsRedemption{BroadcasterID: "1", RewardID: reward.ID}, "")
	log.Printf("%v %#v", err, dbr.Data)

	rewardID = dbr.Data.([]database.ChannelPointsRedemption)[0].RewardID
	redemptionID = dbr.Data.([]database.ChannelPointsRedemption)[0].ID
	db.DB.Close()

	println(rewardID, redemptionID)
	os.Exit(m.Run())
}
func TestRedemption(t *testing.T) {
	a := test_setup.SetupTestEnv(t)

	ts := test_server.SetupTestServer(Redemption{})

	req, _ := http.NewRequest(http.MethodGet, ts.URL+Redemption{}.Path(), nil)
	q := req.URL.Query()
	q.Set("broadcaster_id", "1")
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
	q.Set("reward_id", rewardID)
	req.URL.RawQuery = q.Encode()
	resp, err = http.DefaultClient.Do(req)
	a.Nil(err)
	a.Equal(200, resp.StatusCode)

	q.Set("sort", "potato")
	req.URL.RawQuery = q.Encode()
	resp, err = http.DefaultClient.Do(req)
	a.Nil(err)
	a.Equal(400, resp.StatusCode)

	q.Set("sort", "OLDEST")
	req.URL.RawQuery = q.Encode()
	resp, err = http.DefaultClient.Do(req)
	a.Nil(err)
	a.Equal(200, resp.StatusCode)

	q.Set("id", "1234")
	req.URL.RawQuery = q.Encode()
	resp, err = http.DefaultClient.Do(req)
	a.Nil(err)
	a.Equal(200, resp.StatusCode)

	// patch

	body := PatchRedemptionBody{
		Status: "FULFILLED",
	}

	b, _ := json.Marshal(body)
	req, _ = http.NewRequest(http.MethodPatch, ts.URL+Redemption{}.Path(), bytes.NewBuffer(b))
	req.URL.RawQuery = q.Encode()
	resp, err = http.DefaultClient.Do(req)
	a.Nil(err)
	a.Equal(404, resp.StatusCode)

	b, _ = json.Marshal(body)
	req, _ = http.NewRequest(http.MethodPatch, ts.URL+Redemption{}.Path(), bytes.NewBuffer(b))
	q.Del("id")
	req.URL.RawQuery = q.Encode()
	resp, err = http.DefaultClient.Do(req)
	a.Nil(err)
	a.Equal(400, resp.StatusCode)

	b, _ = json.Marshal(body)
	req, _ = http.NewRequest(http.MethodPatch, ts.URL+Redemption{}.Path(), bytes.NewBuffer(b))
	q.Set("reward_id", rewardID)
	q.Set("id", redemptionID)
	req.URL.RawQuery = q.Encode()
	resp, err = http.DefaultClient.Do(req)
	a.Nil(err)
	a.Equal(200, resp.StatusCode)

	body.Status = "potato"
	b, _ = json.Marshal(body)
	req, _ = http.NewRequest(http.MethodPatch, ts.URL+Redemption{}.Path(), bytes.NewBuffer(b))
	req.URL.RawQuery = q.Encode()
	resp, err = http.DefaultClient.Do(req)
	a.Nil(err)
	a.Equal(400, resp.StatusCode)

}

func TestRewards(t *testing.T) {
	a := test_setup.SetupTestEnv(t)

	ts := test_server.SetupTestServer(Reward{})

	// post tests
	oneHundred := 100
	bTrue := true
	body := PatchAndPostRewardBody{
		Title:     "hello",
		Cost:      &oneHundred,
		IsEnabled: &bTrue,
	}

	b, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodPost, ts.URL+Reward{}.Path(), bytes.NewBuffer(b))
	q := req.URL.Query()
	q.Set("broadcaster_id", "1")
	req.URL.RawQuery = q.Encode()
	resp, err := http.DefaultClient.Do(req)
	a.Nil(err)
	a.NotNil(resp)

	body.Title = ""
	b, _ = json.Marshal(body)
	req, _ = http.NewRequest(http.MethodPost, ts.URL+Reward{}.Path(), bytes.NewBuffer(b))
	q.Set("broadcaster_id", "1")
	req.URL.RawQuery = q.Encode()
	resp, err = http.DefaultClient.Do(req)
	a.Nil(err)
	a.Equal(400, resp.StatusCode)

	b, _ = json.Marshal(body)
	req, _ = http.NewRequest(http.MethodPost, ts.URL+Reward{}.Path(), bytes.NewBuffer(b))
	q.Set("broadcaster_id", "2")
	req.URL.RawQuery = q.Encode()
	resp, err = http.DefaultClient.Do(req)
	a.Nil(err)
	a.Equal(401, resp.StatusCode)

	body.Title = "testing"
	body.IsEnabled = nil
	body.StreamMaxEnabled = true
	b, _ = json.Marshal(body)
	req, _ = http.NewRequest(http.MethodPost, ts.URL+Reward{}.Path(), bytes.NewBuffer(b))
	q.Set("broadcaster_id", "1")
	req.URL.RawQuery = q.Encode()
	resp, err = http.DefaultClient.Do(req)
	a.Nil(err)
	a.Equal(400, resp.StatusCode)

	body.StreamMaxEnabled = false
	body.GlobalCooldownEnabled = true
	b, _ = json.Marshal(body)
	req, _ = http.NewRequest(http.MethodPost, ts.URL+Reward{}.Path(), bytes.NewBuffer(b))
	req.URL.RawQuery = q.Encode()
	resp, err = http.DefaultClient.Do(req)
	a.Nil(err)
	a.Equal(400, resp.StatusCode)

	body.GlobalCooldownEnabled = false
	body.StreamUserMaxEnabled = true
	b, _ = json.Marshal(body)
	req, _ = http.NewRequest(http.MethodPost, ts.URL+Reward{}.Path(), bytes.NewBuffer(b))
	req.URL.RawQuery = q.Encode()
	resp, err = http.DefaultClient.Do(req)
	a.Nil(err)
	a.Equal(400, resp.StatusCode)

	body.StreamUserMaxEnabled = false

	// get
	// good cases
	req, _ = http.NewRequest(http.MethodGet, ts.URL+Reward{}.Path(), nil)
	req.URL.RawQuery = q.Encode()
	resp, err = http.DefaultClient.Do(req)
	a.Nil(err)
	a.Equal(200, resp.StatusCode)

	q.Set("id", rewardID)
	req.URL.RawQuery = q.Encode()
	resp, err = http.DefaultClient.Do(req)
	a.Nil(err)
	a.Equal(200, resp.StatusCode)

	//bad
	q.Del("broadcaster_id")
	req.URL.RawQuery = q.Encode()
	resp, err = http.DefaultClient.Do(req)
	a.Nil(err)
	a.Equal(401, resp.StatusCode)

	// patch
	// no id
	b, _ = json.Marshal(body)
	req, _ = http.NewRequest(http.MethodPatch, ts.URL+Reward{}.Path(), bytes.NewBuffer(b))
	q.Set("broadcaster_id", "1")
	q.Del("id")
	req.URL.RawQuery = q.Encode()
	resp, err = http.DefaultClient.Do(req)
	a.Nil(err)
	a.Equal(400, resp.StatusCode)

	//mismatch bid and token
	b, _ = json.Marshal(body)
	req, _ = http.NewRequest(http.MethodPatch, ts.URL+Reward{}.Path(), bytes.NewBuffer(b))
	q.Set("broadcaster_id", "2")
	req.URL.RawQuery = q.Encode()
	resp, err = http.DefaultClient.Do(req)
	a.Nil(err)
	a.Equal(401, resp.StatusCode)

	// good request
	body.Title = "patched_title"
	b, _ = json.Marshal(body)
	req, _ = http.NewRequest(http.MethodPatch, ts.URL+Reward{}.Path(), bytes.NewBuffer(b))
	q.Set("broadcaster_id", "1")
	q.Set("id", rewardID)
	req.URL.RawQuery = q.Encode()
	resp, err = http.DefaultClient.Do(req)
	a.Nil(err)
	a.Equal(200, resp.StatusCode)

	// bad body
	body.Cost = nil
	b, _ = json.Marshal(body)
	req, _ = http.NewRequest(http.MethodPatch, ts.URL+Reward{}.Path(), bytes.NewBuffer(b))
	q.Set("broadcaster_id", "1")
	q.Set("id", rewardID)
	req.URL.RawQuery = q.Encode()
	resp, err = http.DefaultClient.Do(req)
	a.Nil(err, err)
	a.Equal(400, resp.StatusCode)

	// enabled flag testing below
	body.Cost = &oneHundred
	body.StreamUserMaxEnabled = true
	b, _ = json.Marshal(body)
	req, _ = http.NewRequest(http.MethodPatch, ts.URL+Reward{}.Path(), bytes.NewBuffer(b))
	q.Set("broadcaster_id", "1")
	q.Set("id", rewardID)
	req.URL.RawQuery = q.Encode()
	resp, err = http.DefaultClient.Do(req)
	a.Nil(err, err)
	a.Equal(400, resp.StatusCode)

	body.StreamUserMaxEnabled = false
	body.GlobalCooldownEnabled = true
	b, _ = json.Marshal(body)
	req, _ = http.NewRequest(http.MethodPatch, ts.URL+Reward{}.Path(), bytes.NewBuffer(b))
	q.Set("broadcaster_id", "1")
	q.Set("id", rewardID)
	req.URL.RawQuery = q.Encode()
	resp, err = http.DefaultClient.Do(req)
	a.Nil(err, err)
	a.Equal(400, resp.StatusCode)

	body.GlobalCooldownEnabled = false
	body.StreamMaxEnabled = true
	b, _ = json.Marshal(body)
	req, _ = http.NewRequest(http.MethodPatch, ts.URL+Reward{}.Path(), bytes.NewBuffer(b))
	q.Set("broadcaster_id", "1")
	q.Set("id", rewardID)
	req.URL.RawQuery = q.Encode()
	resp, err = http.DefaultClient.Do(req)
	a.Nil(err, err)
	a.Equal(400, resp.StatusCode)
	body.StreamMaxEnabled = false

	// delete
	req, _ = http.NewRequest(http.MethodDelete, ts.URL+Reward{}.Path(), bytes.NewBuffer(b))
	q.Del("broadcaster_id")
	q.Del("id")
	req.URL.RawQuery = q.Encode()
	resp, err = http.DefaultClient.Do(req)
	a.Nil(err, err)
	a.Equal(401, resp.StatusCode)

	req, _ = http.NewRequest(http.MethodDelete, ts.URL+Reward{}.Path(), bytes.NewBuffer(b))
	q.Set("broadcaster_id", "1")
	q.Del("id")
	req.URL.RawQuery = q.Encode()
	resp, err = http.DefaultClient.Do(req)
	a.Nil(err, err)
	a.Equal(400, resp.StatusCode)

	req, _ = http.NewRequest(http.MethodDelete, ts.URL+Reward{}.Path(), bytes.NewBuffer(b))
	q.Set("broadcaster_id", "1")
	q.Set("id", util.RandomGUID())
	req.URL.RawQuery = q.Encode()
	resp, err = http.DefaultClient.Do(req)
	a.Nil(err, err)
	a.Equal(404, resp.StatusCode)

	req, _ = http.NewRequest(http.MethodDelete, ts.URL+Reward{}.Path(), bytes.NewBuffer(b))
	q.Set("broadcaster_id", "1")
	q.Set("id", rewardID)
	req.URL.RawQuery = q.Encode()
	println(q.Encode())
	resp, err = http.DefaultClient.Do(req)
	a.Nil(err, err)
	a.Equal(204, resp.StatusCode)
}
