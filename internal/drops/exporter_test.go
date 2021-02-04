// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package drops

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/spf13/viper"
	"github.com/twitchdev/twitch-cli/internal/models"
	"github.com/twitchdev/twitch-cli/internal/util"
)

func TestExportEntitlements(t *testing.T) {
	a := util.SetupTestEnv(t)
	viper.Set("clientid", "1111")
	viper.Set("clientsecret", "2222")
	viper.Set("accesstoken", "4567")
	viper.Set("refreshtoken", "123")
	viper.Set("tokenexpiration", "0")

	var okModel = &models.DropsEntitlementsResponse{
		Data: []models.DropsEntitlementsData{
			{
				ID:        "1234",
				BenefitID: "234",
				GameID:    "34",
				UserID:    "4",
				Timestamp: util.GetTimestamp().Format(time.RFC3339),
			},
		},
		Pagination: struct {
			Cursor string "json:\"cursor\""
		}{
			Cursor: "1234",
		},
	}
	ok, err := json.Marshal(okModel)
	a.Nil(err)

	var emptyModel = &models.DropsEntitlementsResponse{
		Data: []models.DropsEntitlementsData{},
	}

	empty, err := json.Marshal(emptyModel)
	a.Nil(err)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/data" {
			w.WriteHeader(http.StatusOK)
			cursor := r.URL.Query().Get("after")
			if cursor != "" {
				w.Write(empty)
				return
			}
			w.Write(ok)
			_, err := ioutil.ReadAll(r.Body)
			a.Nil(err)
		}
		if r.URL.Path == "/empty" {
			w.WriteHeader(http.StatusOK)
			w.Write(empty)
			_, err := ioutil.ReadAll(r.Body)
			a.Nil(err)
		}
		if r.URL.Path == "/error" {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(ok))
			_, err := ioutil.ReadAll(r.Body)
			a.Nil(err)
		}
	}))
	filename := ".testing-csv.csv"

	viper.Set("BASE_URL", ts.URL+"/data")
	ExportEntitlements(filename, "", "")

	viper.Set("BASE_URL", ts.URL+"/empty")
	ExportEntitlements(filename, "1", "")

	viper.Set("BASE_URL", ts.URL+"/error")
	ExportEntitlements(filename, "", "2")
}
