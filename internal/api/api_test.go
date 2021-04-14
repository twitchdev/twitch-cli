// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/spf13/viper"
	"github.com/twitchdev/twitch-cli/internal/util"
)

var params = apiRequestParameters{
	ClientID: "1111",
	Token:    "4567",
}

func TestNewRequest(t *testing.T) {
	a := util.SetupTestEnv(t)
	viper.Set("clientid", "1111")
	viper.Set("clientsecret", "2222")
	viper.Set("accesstoken", "4567")
	viper.Set("refreshtoken", "123")
	viper.Set("tokenexpiration", "0")

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		a.Equal(params.ClientID, r.Header.Get("Client-ID"), "ClientID mismatch")
		a.Equal("Bearer "+params.Token, r.Header.Get("Authorization"), "Token mismatch")

		if r.URL.Path == "/error" {
			w.WriteHeader(http.StatusInternalServerError)
		} else if r.URL.Path == "/nocontent" {
			w.WriteHeader(http.StatusNoContent)
			return
		} else if r.URL.Path == "/cursor" {
			if strings.Contains(r.URL.RawQuery, "after") {
				w.Write([]byte("{}"))
				return
			}
			w.Write([]byte("{\"data\":[],\"pagination\":{\"cursor\":\"test\"}}"))
			return
		}
		w.Write([]byte("{}"))
	}))
	defer ts.Close()

	viper.Set("BASE_URL", ts.URL)
	viper.Set("clientid", "1111")
	viper.Set("clientsecret", "2222")
	viper.Set("accesstoken", "4567")
	viper.Set("refreshtoken", "123")

	// tests for normal get requests
	NewRequest("GET", "", []string{"test=1", "test=2"}, nil, true, false)
	NewRequest("GET", "", []string{"test=1", "test=2"}, nil, false, true)

	// testing cursors autopagination
	NewRequest("GET", "/cursor", []string{"test=1", "test=2"}, nil, false, true)

	// testing 204 no-content apis
	NewRequest("POST", "/nocontent", []string{"test=1", "test=2"}, nil, false, false)

	// testing 500 errors
	NewRequest("GET", "/error", []string{"test=1", "test=2"}, nil, false, true)
}

func TestValidOptions(t *testing.T) {
	a := util.SetupTestEnv(t)

	get := ValidOptions("GET")
	a.NotEmpty(get)

	potato := ValidOptions("potato")
	a.Empty(potato)
}

func TestGetClientInformation(t *testing.T) {
	a := util.SetupTestEnv(t)

	viper.Set("clientid", "1111")
	viper.Set("clientsecret", "2222")
	viper.Set("accesstoken", "4567")
	viper.Set("refreshtoken", "123")

	// check in the future
	viper.Set("tokenexpiration", util.GetTimestamp().Add(10*time.Minute).Format(time.RFC3339Nano))
	clientInfo, err := GetClientInformation()
	a.Nil(err)
	a.Equal(clientInfo.Token, "4567")

	// non-expiring tokens
	viper.Set("tokenexpiration", "0")
	clientInfo, err = GetClientInformation()
	a.Nil(err)
	a.Equal(clientInfo.Token, "4567")

	// expired, but will fail since it's not valid :)
	viper.Set("tokenexpiration", "1")
	clientInfo, err = GetClientInformation()
	a.NotNil(err)
}
