package launchd

import (
	"fmt"
	"net"
	"os"
	"syscall"
)

// Files returns the *os.File for a given socket name
func Files(name string) ([]*os.File, error) {
	fds, err := libxpc_launch_activate_socket(name)
	if err != nil {
		return nil, err
	}
	files := make([]*os.File, len(fds))
	for idx, fd := range fds {
		files[idx] = os.NewFile(uintptr(fd), "")
	}
	return files, nil
}

// Sockets returns the net.Listener for each socket name
func Sockets(name string) ([]net.Listener, error) {
	files, err := Files(name)
	if err != nil {
		return nil, err
	}
	listeners := make([]net.Listener, len(files))
	for idx, file := range files {
		listener, err := net.FileListener(file)
		if err != nil {
			return nil, fmt.Errorf("net.FileListener for %d failed: %w", file.Fd(), err)
		}
		file.Close()
		listeners[idx] = listener
	}
	return listeners, nil
}

// Activates a single net.Listener with the given socket name
// If anything other than a single file descriptor is available syscall.EINVAL is returned
func Activate(name string) (net.Listener, error) {
	listeners, err := Sockets(name)
	if err != nil {
		return nil, err
	} else if len(listeners) != 1 {
		return nil, syscall.EINVAL
	}
	return listeners[0], nil
}
