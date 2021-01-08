// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package util

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRandomUserId(t *testing.T) {
	a := assert.New(t)

	userID := RandomUserID()
	a.NotEqual(0, len(userID), "RandomUserID() returned string with a length of 0")
}

func TestRandomGUID(t *testing.T) {
	a := assert.New(t)
	r, _ := regexp.Compile("^[{]?[0-9a-fA-F]{8}-([0-9a-fA-F]{4}-){3}[0-9a-fA-F]{12}[}]?$")
	guid := RandomGUID()

	a.NotEqual(0, len(guid), "RandomGUID() returned string with a length of 0")
	a.Equal(true, r.MatchString(guid), "RandomGUID() returned a string with value %v, which does not meet the GUID pattern", guid)
}
