package util

import "syscall"

// Syscall wraps syscalls used in the application for unit testing purposes
type Syscall struct {
	Uname func(buf *syscall.Utsname) (err error)
}

var DefaultSyscall = Syscall{
	Uname: syscall.Uname,
}
