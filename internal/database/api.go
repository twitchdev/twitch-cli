// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package database

import (
	"log"
	"strings"
)

type DBResposne struct {
	Cursor string      `json:"cursor"`
	Total  int         `json:"total"`
	Limit  int         `json:"-"`
	Data   interface{} `json:"data"`
}

func (c CLIDatabase) IsFirstRun() bool {
	var userCount = 0

	err := c.DB.Get(&userCount, "select count(*) from users")
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return true
		}
		log.Print(err)
	}

	return userCount == 0
}
