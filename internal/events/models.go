// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package trigger

var triggerTypeMap = map[string]map[string]string{
	"eventsub": {
		"cheer":             "channel.cheer",
		"subscribe":         "channel.subscribe",
		"unsubscribe":       "channel.unsubscribe",
		"gift":              "channel.subscribe",
		"follow":            "channel.follow",
		"transaction":       "",
		"add-redemption":    "channel.channel_points_custom_reward_redemption.add",
		"update-redemption": "channel.channel_points_custom_reward_redemption.update",
		"add-reward":        "channel.channel_points_custom_reward.add",
		"update-reward":     "channel.channel_points_custom_reward.update",
		"remove-reward":     "channel.channel_points_custom_reward.remove",
	},
	"websub": {
		"cheer":             "",
		"subscribe":         "subscribe",
		"unsubscribe":       "subscribe",
		"gift":              "subscribe",
		"follow":            "follow",
		"transaction":       "transaction",
		"add-redemption":    "",
		"update-redemption": "",
		"add-reward":        "",
		"update-reward":     "",
		"remove-reward":     "",
	},
	"websockets": {
		"cheer":             "channel.cheer",
		"subscribe":         "channel.subscribe",
		"unsubscribe":       "channel.unsubscribe",
		"gift":              "channel.subscribe",
		"follow":            "channel.follow",
		"transaction":       "",
		"add-redemption":    "channel.channel_points_custom_reward_redemption.add",
		"update-redemption": "channel.channel_points_custom_reward_redemption.update",
		"add-reward":        "channel.channel_points_custom_reward.add",
		"update-reward":     "channel.channel_points_custom_reward.update",
		"remove-reward":     "channel.channel_points_custom_reward.remove",
	},
}

var triggerSupported = map[string]bool{
	"subscribe":         true,
	"unsubscribe":       true,
	"gift":              true,
	"cheer":             true,
	"transaction":       true,
	"follow":            true,
	"add-redemption":    true,
	"update-redemption": true,
	"add-reward":        true,
	"update-reward":     true,
	"remove-reward":     true,
}

var transportSupported = map[string]bool{
	"websub":     true,
	"eventsub":   true,
	"websockets": false,
}
