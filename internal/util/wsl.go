// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
// +build linux

package util

import (
	"strings"
	"syscall"
)

func IsWsl() bool {
	return isWsl(DefaultSyscall)
}

// check for Windows Subsystem for Linux
func isWsl(sc Syscall) bool {
	// the common factor between WSL distros is the Microsoft-specific kernel version, so we check for that
	// SUSE, WSLv1: 4.4.0-19041-Microsoft
	// Ubuntu, WSLv2: 4.19.128-microsoft-standard
	const wslIdentifier = "microsoft"
	var uname syscall.Utsname
	if err := sc.Uname(&uname); err == nil {
		var kernel []byte
		for _, b := range uname.Release {
			if b == 0 {
				break
			}
			kernel = append(kernel, byte(b))
		}
		return strings.Contains(strings.ToLower(string(kernel)), wslIdentifier)
	}
	return false
}
