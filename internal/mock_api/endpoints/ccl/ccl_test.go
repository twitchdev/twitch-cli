// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package ccl

import (
	"net/http"
	"testing"

	"github.com/twitchdev/twitch-cli/test_setup"
	"github.com/twitchdev/twitch-cli/test_setup/test_server"
)

func TestContentClassificationLabels(t *testing.T) {
	a := test_setup.SetupTestEnv(t)
	ts := test_server.SetupTestServer(ContentClassificationLabels{})

	// get
	req, _ := http.NewRequest(http.MethodGet, ts.URL+ContentClassificationLabels{}.Path(), nil)
	q := req.URL.Query()
	req.URL.RawQuery = q.Encode()
	resp, err := http.DefaultClient.Do(req)
	a.Nil(err)
	a.Equal(200, resp.StatusCode)
}
