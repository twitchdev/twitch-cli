// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package events

import (
	"testing"

	"github.com/twitchdev/twitch-cli/test_setup"
)

func TestValidTransports(t *testing.T) {
	a := test_setup.SetupTestEnv(t)

	t1 := ValidTransports()
	a.NotEmpty(t1)
}
