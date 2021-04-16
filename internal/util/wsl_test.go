// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
// +build linux

package util

import (
	"errors"
	"syscall"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsWsl(t *testing.T) {
	a := assert.New(t)

	var (
		// syscall.Utsname.Release value on various systems

		// Ubuntu 20.04 on WSL2 on Windows 10 x64 20H2
		ubuntu20Wsl2 = [65]int8{52, 46, 49, 57, 46, 49, 50, 56, 45, 109, 105, 99, 114, 111, 115, 111, 102, 116, 45, 115, 116, 97, 110, 100, 97, 114, 100, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

		// Arch Linux on baremetal on 2021-04-02
		archReal = [65]int8{53, 46, 49, 49, 46, 49, 49, 45, 97, 114, 99, 104, 49, 45, 49, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	)

	result := isWsl(Syscall{
		Uname: func(buf *syscall.Utsname) (err error) {
			buf.Release = ubuntu20Wsl2
			return nil
		},
	})
	a.True(result)

	result = isWsl(Syscall{
		Uname: func(buf *syscall.Utsname) (err error) {
			buf.Release = archReal
			return nil
		},
	})
	a.False(result)

	result = isWsl(Syscall{
		Uname: func(buf *syscall.Utsname) (err error) {
			return errors.New("mocked error")
		},
	})
	a.False(result)
}
