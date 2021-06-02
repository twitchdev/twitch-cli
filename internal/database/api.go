// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package database

type DBResposne struct {
	Cursor string      `json:"cursor"`
	Total  int         `json:"total"`
	Limit  int         `json:"-"`
	Data   interface{} `json:"data"`
}
