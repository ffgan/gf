//go:build amd64

#include "textflag.h"

// User assembly is the stack ABI. The compiler inserts the register
// wrapper, so the leaf and the four results live at fixed FP offsets.
TEXT ·cpuidRegs(SB), NOSPLIT, $0-24
	MOVL eax+0(FP), AX
	MOVL ecx+4(FP), CX
	CPUID
	MOVL AX, a+8(FP)
	MOVL BX, b+12(FP)
	MOVL CX, c+16(FP)
	MOVL DX, d+20(FP)
	RET
