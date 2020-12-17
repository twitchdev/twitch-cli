// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package models

type CheerEventSubEvent struct {
	UserID              string  `json:"user_id"`
	UserName            string  `json:"user_name"`
	BroadcasterUserID   string  `json:"broadcast_user_id"`
	BroadcasterUserName string  `json:"broadcast_user_name"`
	IsAnonymous         bool    `json:"is_anonymous"`
	Message             string  `json:"message"`
	Bits                float64 `json:"bits"`
}
