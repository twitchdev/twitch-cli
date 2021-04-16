// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
// +build linux

package util

import "syscall"

// Syscall wraps syscalls used in the application for unit testing purposes
type Syscall struct {
	Uname func(buf *syscall.Utsname) (err error)
}

var DefaultSyscall = Syscall{
	Uname: syscall.Uname,
}
