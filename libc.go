//go:build darwin
// +build darwin

package launchd

import "unsafe"

// free is defined in libxpc.dylib
var libc_free_trampoline_addr uintptr

//go:cgo_import_dynamic libc_free free "/usr/lib/libSystem.B.dylib"

//go:nosplit
func free(ptr unsafe.Pointer) {
	syscall_syscall(libc_free_trampoline_addr, uintptr(ptr), 0, 0)
}
