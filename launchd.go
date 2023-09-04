package launchd

import (
	"fmt"
	"net"
	"os"
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

// File returns a single *os.File for a given socket name
func File(name string) (*os.File, error) {
	files, err := Files(name)
	if err != nil {
		return nil, err
	} else if len(files) != 1 {
		return nil, fmt.Errorf("expected 1 file, got %d", len(files))
	}
	return files[0], nil
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
		listeners[idx] = listener
	}
	return listeners, nil
}

// Socket returns a single net.Listener for a given socket name
func Socket(name string) (net.Listener, error) {
	listeners, err := Sockets(name)
	if err != nil {
		return nil, err
	} else if len(listeners) != 1 {
		return nil, fmt.Errorf("expected 1 listener, got %d", len(listeners))
	}
	return listeners[0], nil
}
