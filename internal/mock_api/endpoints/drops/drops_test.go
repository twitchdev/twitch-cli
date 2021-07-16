// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package drops

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

var entitlement database.DropsEntitlement

func TestMain(m *testing.M) {
	test_setup.SetupTestEnv(&testing.T{})

	db, err := database.NewConnection()
	if err != nil {
		log.Fatal(err)
	}
	e := database.DropsEntitlement{
		ID:        util.RandomGUID(),
		UserID:    "1",
		BenefitID: "1234",
		GameID:    "1",
		Timestamp: util.GetTimestamp().Format(time.RFC3339),
		Status:    "CLAIMED",
	}

	err = db.NewQuery(nil, 100).InsertDropsEntitlement(e)
	if err != nil {
		log.Fatal(err)
	}
	entitlement = e

	db.DB.Close()

	os.Exit(m.Run())
}
func TestDropsEntitlements(t *testing.T) {
	a := test_setup.SetupTestEnv(t)
	ts := test_server.SetupTestServer(DropsEntitlements{})

	// get
	req, _ := http.NewRequest(http.MethodGet, ts.URL+DropsEntitlements{}.Path(), nil)
	q := req.URL.Query()
	req.URL.RawQuery = q.Encode()
	resp, err := http.DefaultClient.Do(req)
	a.Nil(err)
	a.Equal(200, resp.StatusCode)

	// patch
	// patch tests
	body := PatchEntitlementsBody{
		FulfillmentStatus: "FULFILLED",
		EntitlementIDs: []string{
			entitlement.ID,
			"potato",
		},
	}

	b, _ := json.Marshal(body)
	req, _ = http.NewRequest(http.MethodPatch, ts.URL+DropsEntitlements{}.Path(), bytes.NewBuffer(b))
	q = req.URL.Query()
	req.URL.RawQuery = q.Encode()
	resp, err = http.DefaultClient.Do(req)
	a.Nil(err)
	a.NotNil(resp)
	a.Equal(200, resp.StatusCode)

	body.FulfillmentStatus = "potato"
	b, _ = json.Marshal(body)
	req, _ = http.NewRequest(http.MethodPatch, ts.URL+DropsEntitlements{}.Path(), bytes.NewBuffer(b))
	q = req.URL.Query()
	req.URL.RawQuery = q.Encode()
	resp, err = http.DefaultClient.Do(req)
	a.Nil(err)
	a.NotNil(resp)
	a.Equal(400, resp.StatusCode)
}
