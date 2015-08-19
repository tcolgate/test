#include "textflag.h"

// func Sum(t *testing.B)
TEXT ·Sum(SB),NOSPLIT,$0
  MOVQ    (TLS),CX
  CMPQ    SP,16(CX)
  JHI     ,22
  CALL    runtime·morestack_noctxt(SB)
  JMP     ,0
  SUBQ    $8,SP
  FUNCDATA        $0,gclocals·519efd86263089ddb84df3cfe7fd2992+0(SB)
  FUNCDATA        $1,gclocals·3280bececceccd33cb74587feedb1f9f+0(SB)
  MOVQ    "".t+16(FP),BX
  MOVQ    BX,(SP)
  PCDATA  $0,$0
  CALL    testing.(*B)·ResetTimer(SB)
  MOVQ    "".t+16(FP),CX
  MOVQ    $0,AX
  NOP
  MOVQ    112(CX),BX
  CMPQ    BX,AX
  JLE     $0,68
  INCQ    AX
  NOP
  MOVQ    112(CX),BX
  CMPQ    BX,AX
  JGT     $0,56
  ADDQ    $8,SP
  RET

