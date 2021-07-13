// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package models

type TransactionEventSubEvent struct {
	ID                   string                     `json:"id"`
	ExtensionClientID    string                     `json:"extension_client_id"`
	BroadcasterUserID    string                     `json:"broadcaster_user_id"`
	BroadcasterUserLogin string                     `json:"broadcaster_user_login"`
	BroadcasterUserName  string                     `json:"broadcaster_user_name"`
	UserName             string                     `json:"user_name"`
	UserLogin            string                     `json:"user_login"`
	UserID               string                     `json:"user_id"`
	Product              TransactionEventSubProduct `json:"product"`
}

type TransactionEventSubProduct struct {
	Name          string `json:"name"`
	Sku           string `json:"sku"`
	Bits          int64  `json:"bits"`
	InDevelopment bool   `json:"in_development"`
}

type TransactionEventSubResponse struct {
	Subscription EventsubSubscription     `json:"subscription"`
	Event        TransactionEventSubEvent `json:"event"`
}

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
