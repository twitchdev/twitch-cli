// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package models

const TransportEventSub = "eventsub"
const TransportWebSub = "websub"
const TransportWebsockets = "websockets"

var TransportSupported = map[string]bool{
	"websub":     true,
	"eventsub":   true,
	"websockets": false,
}
