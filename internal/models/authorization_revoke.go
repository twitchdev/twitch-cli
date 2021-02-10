// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package models

type AuthorizationRevokeEvent struct {
	UserID    string `json:"user_id"`
	UserLogin string `json:"user_login"`
	UserName  string `json:"user_name"`
	ClientID  string `json:"client_id"`
}

type AuthorizationRevokeEventSubResponse struct {
	Subscription EventsubSubscription     `json:"subscription"`
	Event        AuthorizationRevokeEvent `json:"event"`
}
