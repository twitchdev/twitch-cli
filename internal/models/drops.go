// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package models

type DropsEntitlementsData struct {
	ID        string `json:"id"`
	BenefitID string `json:"benefit_id"`
	Timestamp string `json:"timestamp"`
	UserID    string `json:"user_id"`
	GameID    string `json:"game_id"`
}

type DropsEntitlementsResponse struct {
	Pagination struct {
		Cursor string `json:"cursor"`
	} `json:"pagination"`
	Data []DropsEntitlementsData `json:"data"`
}
