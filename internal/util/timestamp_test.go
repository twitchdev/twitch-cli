// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetTimestamp(t *testing.T) {
	a := assert.New(t)

	ts := GetTimestamp()
	a.NotNil(ts)
}
