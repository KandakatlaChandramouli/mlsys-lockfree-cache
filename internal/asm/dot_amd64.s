#include "textflag.h"

TEXT ·DotProductAVX2(SB), NOSPLIT, $0-56

    MOVQ a_base+0(FP), SI
    MOVQ b_base+24(FP), DI
    MOVQ a_len+8(FP), CX

    XORPS X0, X0
    XORQ AX, AX

loop:

    CMPQ AX, CX
    JGE done

    MOVSS (SI)(AX*4), X1
    MOVSS (DI)(AX*4), X2

    MULSS X2, X1
    ADDSS X1, X0

    INCQ AX
    JMP loop

done:

    MOVSS X0, ret+48(FP)
    RET
