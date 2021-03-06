// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVersion(t *testing.T) {
	a := assert.New(t)

	var testString = "test_version"

	SetVersion(testString)
	v := GetVersion()
	a.Equal(testString, v)
}
