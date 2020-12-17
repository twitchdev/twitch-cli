// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package util

var version = "source"

// SetVersion sets the version number for use later in the version command and for request headers.
func SetVersion(v string) {
	version = v
}

// GetVersion retrieves the version number for use later in the version command and for request headers.
func GetVersion() string {
	return version
}
