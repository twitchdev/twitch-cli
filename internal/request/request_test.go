// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package request

import (
	"testing"

	"github.com/twitchdev/twitch-cli/internal/util"
)

func TestNewRequest(t *testing.T) {
	a := util.SetupTestEnv(t)

	r, err := NewRequest("GET", "https://api.twitch.tv/helix/users", nil)
	a.Nil(err)
	a.Contains(r.Header.Get("User-Agent"), "twitch-cli/")
}
