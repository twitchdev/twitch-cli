// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package trigger

var triggerTypeMap = map[string]map[string]string{
	"eventsub": {
		"cheer":       "channel.cheer",
		"subscribe":   "channel.subscribe",
		"unsubscribe": "channel.unsubscribe",
		"gift":        "channel.subscribe",
		"follow":      "channel.follow",
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
		"cheer":       "channel.cheer",
		"subscribe":   "channel.subscribe",
		"unsubscribe": "channel.unsubscribe",
		"gift":        "channel.subscribe",
		"follow":      "channel.follow",
		"transaction": "",
	},
}

var triggerSupported = map[string]bool{
	"subscribe":   true,
	"unsubscribe": true,
	"gift":        true,
	"cheer":       true,
	"transaction": true,
	"follow":      true,
}

var transportSupported = map[string]bool{
	"websub":     true,
	"eventsub":   true,
	"websockets": false,
}
