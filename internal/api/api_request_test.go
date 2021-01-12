// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package api

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/twitchdev/twitch-cli/internal/util"
)

func TestApiRequest(t *testing.T) {
	a := util.SetupTestEnv(t)

	var ok = "{\"status\":\"ok\"}"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(ok))

		a.Equal("1234", r.Header.Get("Client-ID"), "Client ID invalid.")
		a.Equal("Bearer 4567", r.Header.Get("Authorization"), "Token invalid.")

		_, err := ioutil.ReadAll(r.Body)
		a.Nil(err)
	}))

	defer ts.Close()

	params := *&apiRequestParameters{
		ClientID: "1234",
		Token:    "4567",
	}

	resp, err := apiRequest(http.MethodGet, ts.URL, nil, params)
	a.Nil(err)

	a.Equal(http.StatusOK, resp.StatusCode)
	a.Equal(string(resp.Body), ok)

	resp, err = apiRequest(http.MethodGet, "potato", nil, params)
	a.NotNil(err)
}
