//go:build darwin
// +build darwin

#include "textflag.h"

// The trampolines are ABIInternal as they are address-taken in
// Go code.

TEXT libc_free_trampoline<>(SB),NOSPLIT,$0-0
    JMP	libc_free(SB)

GLOBL	·libc_free_trampoline_addr(SB), RODATA, $8
DATA	·libc_free_trampoline_addr(SB)/8, $libc_free_trampoline<>(SB)
