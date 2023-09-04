# go-launchd [![Go Reference](https://pkg.go.dev/badge/github.com/bored-engineer/go-launchd.svg)](https://pkg.go.dev/github.com/bored-engineer/go-launchd) [![Test Workflow](https://github.com/bored-engineer/go-launchd/actions/workflows/test.yml/badge.svg)](https://github.com/bored-engineer/go-launchd/actions/workflows/test.yml)
[Golang](https://go.dev/) support for macOS (Darwin) launchd socket activation ([launch_activate_socket](https://developer.apple.com/documentation/xpc/1505523-launch_activate_socket)) without [cgo](https://pkg.go.dev/cmd/cgo).

## How it works
Using a similar technique as the [crypto/x509 package](https://go-review.googlesource.com/c/go/+/232397) and [golang.org/x/sys/unix package](https://pkg.go.dev/golang.org/x/sys/unix) it is possible to import the [launch_activate_socket](https://developer.apple.com/documentation/xpc/1505523-launch_activate_socket) function from `libxpc.dylib` at runtime via the `go:cgo_import_dynamic` directive:
```go
// launch_activate_socket is defined in libxpc.dylib
var libxpc_launch_activate_socket_trampoline_addr uintptr

//go:cgo_import_dynamic libxpc_launch_activate_socket launch_activate_socket "/usr/lib/system/libxpc.dylib"
```
This must be combined with a golang assembly file ([libxpc.s](./libxpc.s)) to populate the "trampoline" variable with a pointer to the relevant function:
```
TEXT libxpc_launch_activate_socket_trampoline<>(SB),NOSPLIT,$0-0
    JMP	libxpc_launch_activate_socket(SB)

GLOBL	·libxpc_launch_activate_socket_trampoline_addr(SB), RODATA, $8
DATA	·libxpc_launch_activate_socket_trampoline_addr(SB)/8, $libxpc_launch_activate_socket_trampoline<>(SB)
```
Finally by linking the internal `syscall.syscall` function we are able to invoke the [launch_activate_socket](https://developer.apple.com/documentation/xpc/1505523-launch_activate_socket) function:
```go
// Implemented in the runtime package (runtime/sys_darwin.go)
func syscall_syscall(fn, a1, a2, a3 uintptr) (r1, r2 uintptr, err syscall.Errno)

//go:linkname syscall_syscall syscall.syscall
```
The remaining logic is primarily implemented in [libxpc.go](./libxpc.go) to convert from C structures into Golang equivalents.

## Usage
See also the [example/](./example/) directory for the associated launchd job definition (plist):
```go
package main

import launchd "github.com/bored-engineer/go-launchd"

func main() {
	l, err := launchd.Socket("Listener")
	if err != nil {
		log.Fatalf("launchd.Socket failed: %s", err)
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Printf("(net.Listener).Accept failed: %s", err)
			continue
		}
		go func(conn net.Conn) {
			defer conn.Close()
			io.Copy(conn, conn)
		}(conn)
	}
}

```
