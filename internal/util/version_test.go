// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package util

import (
	"testing"
)

func TestVersion(t *testing.T) {
	var testString = "test_version"

	SetVersion(testString)

	v := GetVersion()

	if v != testString {
		t.Errorf("Version failed, set version to %v, received %v", testString, v)
	}
}
