// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package events

var triggerSupported = map[string]bool{
	"subscribe":         	true,
	"unsubscribe":       	true,
	"gift":              	true,
	"cheer":             	true,
	"transaction":       	true,
	"follow":            	true,
	"add-redemption":    	true,
	"update-redemption": 	true,
	"add-reward":        	true,
	"update-reward":     	true,
	"remove-reward":     	true,
	"add-moderator":     	true,
	"remove-moderator":  	true,
	"ban":					true,
	"unban": 			 	true,
	"hype-train-begin": 	true,
	"hype-train-progress":  true,
	"hype-train-end":       true,
}

var transportSupported = map[string]bool{
	"websub":     true,
	"eventsub":   true,
	"websockets": false,
}
