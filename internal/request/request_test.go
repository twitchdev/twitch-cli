// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package request

import (
	"testing"

	"github.com/twitchdev/twitch-cli/test_setup"
)

func TestNewRequest(t *testing.T) {
	a := test_setup.SetupTestEnv(t)

	r, err := NewRequest("GET", "https://api.twitch.tv/helix/users", nil)
	a.Nil(err)
	a.Contains(r.Header.Get("User-Agent"), "twitch-cli/")
}
