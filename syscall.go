//go:build darwin
// +build darwin

package launchd

import (
	"syscall"
	_ "unsafe"
)

// Implemented in the runtime package (runtime/sys_darwin.go)
func syscall_syscall(fn, a1, a2, a3 uintptr) (r1, r2 uintptr, err syscall.Errno)

//go:linkname syscall_syscall syscall.syscall
