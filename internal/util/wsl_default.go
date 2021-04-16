// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
// +build !linux

package util

// non-linux platforms cannot be WSL
func IsWsl() bool {
	return false
}
