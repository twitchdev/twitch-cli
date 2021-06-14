// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package generate

import (
	"testing"

	"github.com/twitchdev/twitch-cli/test_setup"
)

func TestGenerateUsername(t *testing.T) {
	a := test_setup.SetupTestEnv(t)

	un := generateUsername()
	a.NotEmpty(un)
}

func TestGenerate(t *testing.T) {
	a := test_setup.SetupTestEnv(t)

	err := Generate(0)
	a.Nil(err)

	err = Generate(10)
	a.Nil(err)
}
