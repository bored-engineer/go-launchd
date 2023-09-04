//go:build darwin
// +build darwin

#include "textflag.h"

// The trampolines are ABIInternal as they are address-taken in
// Go code.

TEXT libxpc_launch_activate_socket_trampoline<>(SB),NOSPLIT,$0-0
    JMP	libxpc_launch_activate_socket(SB)

GLOBL	·libxpc_launch_activate_socket_trampoline_addr(SB), RODATA, $8
DATA	·libxpc_launch_activate_socket_trampoline_addr(SB)/8, $libxpc_launch_activate_socket_trampoline<>(SB)
