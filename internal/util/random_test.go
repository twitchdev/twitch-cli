// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package util

import (
	"regexp"
	"testing"
)

func TestRandomUserId(t *testing.T) {
	userID := RandomUserID()

	if len(userID) == 0 {
		t.Errorf("RandomUserID() returned string with a length of 0")
	}
}

func TestRandomGUID(t *testing.T) {
	guid := RandomGUID()
	if len(guid) == 0 {
		t.Errorf("RandomGUID() returned string with a length of 0")
	}

	r, _ := regexp.Compile("^[{]?[0-9a-fA-F]{8}-([0-9a-fA-F]{4}-){3}[0-9a-fA-F]{12}[}]?$")

	if r.MatchString(guid) != true {
		t.Errorf("RandomGUID() returned a string with value %v, which does not meet the GUID pattern", guid)
	}
}
