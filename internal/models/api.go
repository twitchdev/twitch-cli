// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package models

type APIResponse struct {
	Data       any                       `json:"data,omitempty"`
	Pagination *APIPagination            `json:"pagination,omitempty"`
	Error      string                    `json:"error,omitempty"`
	Status     int                       `json:"status,omitempty"`
	Message    string                    `json:"message,omitempty"`
	Template   string                    `json:"template,omitempty"`
	Total      *int                      `json:"total,omitempty"`
	DateRange  *BitsLeaderboardDateRange `json:"date_range,omitempty"`
}

type APIPagination struct {
	Cursor string `json:"cursor,omitempty"`
}

type BitsLeaderboardDateRange struct {
	StartedAt string `json:"started_at"`
	EndedAt   string `json:"ended_at"`
}

type ExtensionAPIResponse struct { // extensions/live
	Data       interface{} `json:"data,omitempty"`
	Pagination *string     `json:"pagination,omitempty"`
	Error      string      `json:"error,omitempty"`
	Status     int         `json:"status,omitempty"`
	Message    string      `json:"message,omitempty"`
}
