// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package database

type DBResposne struct {
	Cursor string
	Total  string
	Limit  int
	Data   interface{}
}

type DBPagination struct {
	Limit  int
	Cursor string
}
