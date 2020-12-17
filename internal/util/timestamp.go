// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package util

import "time"

// GetTimestamp returns the timestamp in UTC for use with signature creation and event firing.
func GetTimestamp() time.Time {
	return time.Now().UTC()
}
