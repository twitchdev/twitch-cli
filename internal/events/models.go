// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package events

var transportSupported = map[string]bool{
	"websub":     false,
	"eventsub":   true,
	"websockets": false,
}
