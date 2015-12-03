#include "textflag.h"

// func Ret(x uint64) uint64
TEXT ·Ret(SB),NOSPLIT,$0
MOVQ x+0(FP), AX
MOVQ AX, ret+8(FP)
RET

// func Sum(x, y uint64) uint64
TEXT ·Sum(SB),NOSPLIT,$0
MOVQ x+0(FP), AX
ADDQ y+8(FP), AX
MOVQ AX, ret+16(FP)
RET
