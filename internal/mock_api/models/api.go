// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package models

type APIPagination struct {
	Cursor *string `json:"cursor"`
}

type APIResponse struct {
	Total      *int           `json:"total,omitempty"`
	Data       interface{}    `json:"data"`
	Pagination *APIPagination `json:"pagination,omitempty"`
}
