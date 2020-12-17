// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package trigger

var triggerTypeMap = map[string]map[string]string{
	"eventsub": {
		"cheer":       "channels.cheer",
		"subscribe":   "channels.subscribe",
		"unsubscribe": "channels.unsubscribe",
		"gift":        "channels.subscribe",
		"follow":      "users.follow",
		"transaction": "",
	},
	"websub": {
		"cheer":       "",
		"subscribe":   "subscribe",
		"unsubscribe": "subscribe",
		"gift":        "subscribe",
		"follow":      "follow",
		"transaction": "transaction",
	},
	"websockets": {
		"cheer":       "channels.cheer",
		"subscribe":   "channels.subscribe",
		"unsubscribe": "channels.unsubscribe",
		"gift":        "channels.subscribe",
		"follow":      "users.follow",
		"transaction": "",
	},
}

var triggerSupported = map[string]bool{
	"subscribe":   true,
	"unsubscribe": true,
	"gift":        true,
	"cheer":       true,
	"transaction": true,
	"follow":      false,
}

var transportSupported = map[string]bool{
	"websub":     true,
	"eventsub":   true,
	"websockets": false,
}
