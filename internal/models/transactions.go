// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package models

type TransactionWebsubEvent struct {
	ID              string             `json:"id"`
	Timestamp       string             `json:"timestamp"`
	BroadcasterID   string             `json:"broadcaster_id"`
	BroadcasterName string             `json:"broadcaster_name"`
	UserID          string             `json:"user_id"`
	UserName        string             `json:"user_name"`
	ProductType     string             `json:"product_type"`
	Product         TransactionProduct `json:"product_data"`
}

type TransactionProduct struct {
	Sku           string          `json:"sku"`
	Cost          TransactionCost `json:"cost"`
	DisplayName   string          `json:"displayName"`
	InDevelopment bool            `json:"inDevelopment"`
	Broadcast     bool            `json:"broadcast"`
	Domain        string          `json:"domain"`
	Expiration    string          `json:"expiration"`
}

type TransactionCost struct {
	Amount int64  `json:"amount"`
	Type   string `json:"type"`
}

type TransactionWebSubResponse struct {
	Data []TransactionWebsubEvent `json:"data"`
}
