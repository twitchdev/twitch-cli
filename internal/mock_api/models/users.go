// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package models

type UserAPIData struct {
	ID              string `json:"id"`
	Login           string `json:"login"`
	DisplayName     string `json:"display_name"`
	Type            string `json:"type"`
	BroadcasterType string `json:"broadcaster_type"`
	Description     string `json:"description"`
	ProfileImageURL string `json:"profile_image_url"`
	OfflineImageURL string `json:"offline_image_url"`
	ViewCount       int64  `json:"view_count"`
	Email           string `json:"email,omitempty"`
	CreatedAt       string `json:"created_at"`
}

type UserFollowsAPIData struct {
	FromID     string `json:"from_id"`
	FromLogin  string `json:"from_login"`
	FromName   string `json:"from_name"`
	ToID       string `json:"to_id"`
	ToLogin    string `json:"to_login"`
	ToName     string `json:"to_name"`
	FollowedAt string `json:"followed_at"`
}

type UserBlocksAPIData struct {
	ID          string `json:"user_id"`
	Login       string `json:"user_login"`
	DisplayName string `json:"display_name"`
}
