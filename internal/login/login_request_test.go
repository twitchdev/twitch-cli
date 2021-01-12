// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package login

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/twitchdev/twitch-cli/internal/util"
)

func TestLoginRequest(t *testing.T) {
	a := util.SetupTestEnv(t)

	var ok = "{\"status\":\"ok\"}"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(ok))

		_, err := ioutil.ReadAll(r.Body)
		a.Nil(err)

	}))

	defer ts.Close()

	resp, err := loginRequest(http.MethodGet, ts.URL, nil)
	a.Nil(err)
	a.Equal(http.StatusOK, resp.StatusCode, "Expected status %v, got %v")
	a.Equal(ok, string(resp.Body), "Expected %v, got %v", ok, resp.Body)
}
