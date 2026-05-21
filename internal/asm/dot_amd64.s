#include "textflag.h"

TEXT ·DotProductAVX2(SB), NOSPLIT, $0-56

    MOVQ a_base+0(FP), SI
    MOVQ b_base+24(FP), DI
    MOVQ a_len+8(FP), CX

    VXORPS Y0, Y0, Y0
    XORQ AX, AX

loop:

    CMPQ AX, CX
    JGE done

    VMOVUPS (SI)(AX*4), Y1
    VMOVUPS (DI)(AX*4), Y2

    VFMADD231PS Y1, Y2, Y0

    ADDQ $8, AX
    JMP loop

done:

    VEXTRACTF128 $1, Y0, X1
    VADDPS X1, X0, X0

    VHADDPS X0, X0, X0
    VHADDPS X0, X0, X0

    MOVSS X0, ret+48(FP)

    VZEROUPPER
    RET
