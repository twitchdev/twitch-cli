// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package models

type UserUpdateEventSubResponse struct {
    Subscription EventsubSubscription  `json:"subscription"`
    Event        StreamUpEventSubEvent `json:"event"`
}

type UserUpdateEventSubEvent struct {
    UserID               string `json:"user_id"`
    UserLogin            string `json:"user_login"`
    UserName             string `json:"user_name"`
    Description          string `json:"description"`
}
