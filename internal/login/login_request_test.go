// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package login

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLoginRequest(t *testing.T) {
	var ok = "{\"status\":\"ok\"}"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(ok))

		_, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Error(err)
		}

	}))

	defer ts.Close()

	resp, err := loginRequest(http.MethodGet, ts.URL, nil)
	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status %v, got %v", http.StatusOK, resp.StatusCode)
	}

	if string(resp.Body) != ok {
		t.Errorf("Expected %v, got %v", ok, resp.Body)
	}

}
