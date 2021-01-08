// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package util

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetApplicationDir(t *testing.T) {
	a := assert.New(t)

	dir, err := GetApplicationDir()
	a.Nil(err, "GetApplicationDir() failed with error  %v", err)
	a.Equal(true, strings.HasSuffix(dir, ".twitch-cli"), "GetApplicationDir() expected to end with %v, got %v", ".twitch-cli", dir)
}

func TestGetConfigPath(t *testing.T) {
	a := assert.New(t)

	config, err := GetConfigPath()
	a.Nil(err, "GetConfigPath() failed with error  %v", err)
	a.Equal(true, strings.HasSuffix(config, ".twitch-cli.env"), "GetConfigPath() expected to end with %v, got %v", ".twitch-cli", config)
}
