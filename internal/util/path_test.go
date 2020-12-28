// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package util

import (
	"strings"
	"testing"
)

func TestGetApplicationDir(t *testing.T) {
	dir, err := GetApplicationDir()
	if err != nil {
		t.Errorf("GetApplicationDir() failed with error  %v", err)
	}

	if strings.HasSuffix(dir, ".twitch-cli") != true {
		t.Errorf("GetApplicationDir() expected to end with %v, got %v", ".twitch-cli", dir)
	}
}
