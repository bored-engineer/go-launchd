//go:build !darwin
// +build !darwin

package launchd

import "syscall"

// Supported is true if launchd socket activation is supported on this platform
const Supported = false

// Always returns syscall.ENOSYS on unsupported platforms
func libxpc_launch_activate_socket(name string) ([]int, error) {
	return nil, syscall.ENOSYS
}
