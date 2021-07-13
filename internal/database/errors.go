// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package database

import (
	"reflect"

	"github.com/mattn/go-sqlite3"
)

// Wrapper for the various SQLite errors to allow for easier error checking
func DatabaseErrorIs(err error, sqliteError error) bool {
	if reflect.TypeOf(err).String() == "sqlite3.Error" {
		tempErr := err.(sqlite3.Error)
		switch v := sqliteError.(type) {
		case sqlite3.ErrNo:
			return tempErr.Code == v
		case sqlite3.ErrNoExtended:
			return tempErr.ExtendedCode == v
		}
	}

	return false
}
