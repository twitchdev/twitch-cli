// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package events

var triggerSupported = map[string]bool{
	"add-moderator":       true,
	"add-redemption":      true,
	"add-reward":          true,
	"ban":                 true,
	"cheer":               true,
	"drop":                true,
	"follow":              true,
	"gift":                true,
	"goal-begin":          true,
	"goal-end":            true,
	"goal-progress":       true,
	"grant":               true,
	"hype-train-begin":    true,
	"hype-train-end":      true,
	"hype-train-progress": true,
	"poll-begin":          true,
	"poll-progress":       true,
	"poll-end":            true,
	"prediction-begin":    true,
	"prediction-progress": true,
	"prediction-lock":     true,
	"prediction-end":      true,
	"raid":                true,
	"remove-moderator":    true,
	"remove-reward":       true,
	"revoke":              true,
	"stream-change":       true,
	"streamdown":          true,
	"streamup":            true,
	"subscribe":           true,
	"transaction":         true,
	"unban":               true,
	"unsubscribe":         true,
	"update-redemption":   true,
	"update-reward":       true,
	"user-update":         true,
}

var transportSupported = map[string]bool{
	"websub":     false,
	"eventsub":   true,
	"websockets": false,
}
