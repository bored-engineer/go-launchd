//go:build darwin
// +build darwin

package launchd

import (
	"runtime"
	"syscall"
	"unsafe"
)

// Supported is true if launchd socket activation is supported on this platform
const Supported = true

// As a sanity check, refuse any returned array larger than this with syscall.EINVAL
const maxFDs = 1 << 20

// launch_activate_socket is defined in libxpc.dylib
var libxpc_launch_activate_socket_trampoline_addr uintptr

//go:cgo_import_dynamic libxpc_launch_activate_socket launch_activate_socket "/usr/lib/system/libxpc.dylib"

// invokes launch_activate_socket
func libxpc_launch_activate_socket(name string) ([]int, error) {
	c_name_ptr, err := syscall.BytePtrFromString(name)
	if err != nil {
		return nil, err
	}
	var c_fds_ptr *uintptr // *C.int
	defer func() {
		if c_fds_ptr != nil {
			free(unsafe.Pointer(c_fds_ptr))
		}
	}()
	var c_cnt uint // C.size_t
	res, _, _ := syscall_syscall(
		libxpc_launch_activate_socket_trampoline_addr,
		uintptr(unsafe.Pointer(c_name_ptr)),
		uintptr(unsafe.Pointer(&c_fds_ptr)),
		uintptr(unsafe.Pointer(&c_cnt)),
	)
	runtime.KeepAlive(c_name_ptr)
	if res != 0 {
		return nil, syscall.Errno(res)
	} else if c_cnt > maxFDs {
		return nil, syscall.EINVAL
	}
	c_fds := (*[maxFDs]int)(unsafe.Pointer(c_fds_ptr))
	fds := make([]int, c_cnt)
	copy(fds, (*c_fds)[0:c_cnt])
	return fds, nil
}
