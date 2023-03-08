// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package models_mock_websocket

import "github.com/twitchdev/twitch-cli/internal/models"

type ReconnectWebSocketEvent struct {
}

type ReconnectWebSocketResponse struct {
	Subscription models.EventsubSubscription `json:"ban"`
	Event        ReconnectWebSocketEvent     `json:"event"`
}
