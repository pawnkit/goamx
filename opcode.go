package amx

import "github.com/pawnkit/goamx/vm"

// Opcode is a 32-bit AMX instruction opcode.
type Opcode = vm.Opcode

const (
	OP_NONE        Opcode = vm.OP_NONE
	OP_LOAD_PRI    Opcode = vm.OP_LOAD_PRI
	OP_LOAD_ALT    Opcode = vm.OP_LOAD_ALT
	OP_LOAD_S_PRI  Opcode = vm.OP_LOAD_S_PRI
	OP_LOAD_S_ALT  Opcode = vm.OP_LOAD_S_ALT
	OP_LREF_PRI    Opcode = vm.OP_LREF_PRI
	OP_LREF_ALT    Opcode = vm.OP_LREF_ALT
	OP_LREF_S_PRI  Opcode = vm.OP_LREF_S_PRI
	OP_LREF_S_ALT  Opcode = vm.OP_LREF_S_ALT
	OP_LOAD_I      Opcode = vm.OP_LOAD_I
	OP_LODB_I      Opcode = vm.OP_LODB_I
	OP_CONST_PRI   Opcode = vm.OP_CONST_PRI
	OP_CONST_ALT   Opcode = vm.OP_CONST_ALT
	OP_ADDR_PRI    Opcode = vm.OP_ADDR_PRI
	OP_ADDR_ALT    Opcode = vm.OP_ADDR_ALT
	OP_STOR_PRI    Opcode = vm.OP_STOR_PRI
	OP_STOR_ALT    Opcode = vm.OP_STOR_ALT
	OP_STOR_S_PRI  Opcode = vm.OP_STOR_S_PRI
	OP_STOR_S_ALT  Opcode = vm.OP_STOR_S_ALT
	OP_SREF_PRI    Opcode = vm.OP_SREF_PRI
	OP_SREF_ALT    Opcode = vm.OP_SREF_ALT
	OP_SREF_S_PRI  Opcode = vm.OP_SREF_S_PRI
	OP_SREF_S_ALT  Opcode = vm.OP_SREF_S_ALT
	OP_STOR_I      Opcode = vm.OP_STOR_I
	OP_STRB_I      Opcode = vm.OP_STRB_I
	OP_LIDX        Opcode = vm.OP_LIDX
	OP_LIDX_B      Opcode = vm.OP_LIDX_B
	OP_IDXADDR     Opcode = vm.OP_IDXADDR
	OP_IDXADDR_B   Opcode = vm.OP_IDXADDR_B
	OP_ALIGN_PRI   Opcode = vm.OP_ALIGN_PRI
	OP_ALIGN_ALT   Opcode = vm.OP_ALIGN_ALT
	OP_LCTRL       Opcode = vm.OP_LCTRL
	OP_SCTRL       Opcode = vm.OP_SCTRL
	OP_MOVE_PRI    Opcode = vm.OP_MOVE_PRI
	OP_MOVE_ALT    Opcode = vm.OP_MOVE_ALT
	OP_XCHG        Opcode = vm.OP_XCHG
	OP_PUSH_PRI    Opcode = vm.OP_PUSH_PRI
	OP_PUSH_ALT    Opcode = vm.OP_PUSH_ALT
	OP_PUSH_R      Opcode = vm.OP_PUSH_R
	OP_PUSH_C      Opcode = vm.OP_PUSH_C
	OP_PUSH        Opcode = vm.OP_PUSH
	OP_PUSH_S      Opcode = vm.OP_PUSH_S
	OP_POP_PRI     Opcode = vm.OP_POP_PRI
	OP_POP_ALT     Opcode = vm.OP_POP_ALT
	OP_STACK       Opcode = vm.OP_STACK
	OP_HEAP        Opcode = vm.OP_HEAP
	OP_PROC        Opcode = vm.OP_PROC
	OP_RET         Opcode = vm.OP_RET
	OP_RETN        Opcode = vm.OP_RETN
	OP_CALL        Opcode = vm.OP_CALL
	OP_CALL_PRI    Opcode = vm.OP_CALL_PRI
	OP_JUMP        Opcode = vm.OP_JUMP
	OP_JREL        Opcode = vm.OP_JREL
	OP_JZER        Opcode = vm.OP_JZER
	OP_JNZ         Opcode = vm.OP_JNZ
	OP_JEQ         Opcode = vm.OP_JEQ
	OP_JNEQ        Opcode = vm.OP_JNEQ
	OP_JLESS       Opcode = vm.OP_JLESS
	OP_JLEQ        Opcode = vm.OP_JLEQ
	OP_JGRTR       Opcode = vm.OP_JGRTR
	OP_JGEQ        Opcode = vm.OP_JGEQ
	OP_JSLESS      Opcode = vm.OP_JSLESS
	OP_JSLEQ       Opcode = vm.OP_JSLEQ
	OP_JSGRTR      Opcode = vm.OP_JSGRTR
	OP_JSGEQ       Opcode = vm.OP_JSGEQ
	OP_SHL         Opcode = vm.OP_SHL
	OP_SHR         Opcode = vm.OP_SHR
	OP_SSHR        Opcode = vm.OP_SSHR
	OP_SHL_C_PRI   Opcode = vm.OP_SHL_C_PRI
	OP_SHL_C_ALT   Opcode = vm.OP_SHL_C_ALT
	OP_SHR_C_PRI   Opcode = vm.OP_SHR_C_PRI
	OP_SHR_C_ALT   Opcode = vm.OP_SHR_C_ALT
	OP_SMUL        Opcode = vm.OP_SMUL
	OP_SDIV        Opcode = vm.OP_SDIV
	OP_SDIV_ALT    Opcode = vm.OP_SDIV_ALT
	OP_UMUL        Opcode = vm.OP_UMUL
	OP_UDIV        Opcode = vm.OP_UDIV
	OP_UDIV_ALT    Opcode = vm.OP_UDIV_ALT
	OP_ADD         Opcode = vm.OP_ADD
	OP_SUB         Opcode = vm.OP_SUB
	OP_SUB_ALT     Opcode = vm.OP_SUB_ALT
	OP_AND         Opcode = vm.OP_AND
	OP_OR          Opcode = vm.OP_OR
	OP_XOR         Opcode = vm.OP_XOR
	OP_NOT         Opcode = vm.OP_NOT
	OP_NEG         Opcode = vm.OP_NEG
	OP_INVERT      Opcode = vm.OP_INVERT
	OP_ADD_C       Opcode = vm.OP_ADD_C
	OP_SMUL_C      Opcode = vm.OP_SMUL_C
	OP_ZERO_PRI    Opcode = vm.OP_ZERO_PRI
	OP_ZERO_ALT    Opcode = vm.OP_ZERO_ALT
	OP_ZERO        Opcode = vm.OP_ZERO
	OP_ZERO_S      Opcode = vm.OP_ZERO_S
	OP_SIGN_PRI    Opcode = vm.OP_SIGN_PRI
	OP_SIGN_ALT    Opcode = vm.OP_SIGN_ALT
	OP_EQ          Opcode = vm.OP_EQ
	OP_NEQ         Opcode = vm.OP_NEQ
	OP_LESS        Opcode = vm.OP_LESS
	OP_LEQ         Opcode = vm.OP_LEQ
	OP_GRTR        Opcode = vm.OP_GRTR
	OP_GEQ         Opcode = vm.OP_GEQ
	OP_SLESS       Opcode = vm.OP_SLESS
	OP_SLEQ        Opcode = vm.OP_SLEQ
	OP_SGRTR       Opcode = vm.OP_SGRTR
	OP_SGEQ        Opcode = vm.OP_SGEQ
	OP_EQ_C_PRI    Opcode = vm.OP_EQ_C_PRI
	OP_EQ_C_ALT    Opcode = vm.OP_EQ_C_ALT
	OP_INC_PRI     Opcode = vm.OP_INC_PRI
	OP_INC_ALT     Opcode = vm.OP_INC_ALT
	OP_INC         Opcode = vm.OP_INC
	OP_INC_S       Opcode = vm.OP_INC_S
	OP_INC_I       Opcode = vm.OP_INC_I
	OP_DEC_PRI     Opcode = vm.OP_DEC_PRI
	OP_DEC_ALT     Opcode = vm.OP_DEC_ALT
	OP_DEC         Opcode = vm.OP_DEC
	OP_DEC_S       Opcode = vm.OP_DEC_S
	OP_DEC_I       Opcode = vm.OP_DEC_I
	OP_MOVS        Opcode = vm.OP_MOVS
	OP_CMPS        Opcode = vm.OP_CMPS
	OP_FILL        Opcode = vm.OP_FILL
	OP_HALT        Opcode = vm.OP_HALT
	OP_BOUNDS      Opcode = vm.OP_BOUNDS
	OP_SYSREQ_PRI  Opcode = vm.OP_SYSREQ_PRI
	OP_SYSREQ_C    Opcode = vm.OP_SYSREQ_C
	OP_FILE        Opcode = vm.OP_FILE
	OP_LINE        Opcode = vm.OP_LINE
	OP_SYMBOL      Opcode = vm.OP_SYMBOL
	OP_SRANGE      Opcode = vm.OP_SRANGE
	OP_JUMP_PRI    Opcode = vm.OP_JUMP_PRI
	OP_SWITCH      Opcode = vm.OP_SWITCH
	OP_CASETBL     Opcode = vm.OP_CASETBL
	OP_SWAP_PRI    Opcode = vm.OP_SWAP_PRI
	OP_SWAP_ALT    Opcode = vm.OP_SWAP_ALT
	OP_PUSH_ADR    Opcode = vm.OP_PUSH_ADR
	OP_NOP         Opcode = vm.OP_NOP
	OP_SYSREQ_N    Opcode = vm.OP_SYSREQ_N
	OP_SYMTAG      Opcode = vm.OP_SYMTAG
	OP_BREAK       Opcode = vm.OP_BREAK
	OP_PUSH2_C     Opcode = vm.OP_PUSH2_C
	OP_PUSH2       Opcode = vm.OP_PUSH2
	OP_PUSH2_S     Opcode = vm.OP_PUSH2_S
	OP_PUSH2_ADR   Opcode = vm.OP_PUSH2_ADR
	OP_PUSH3_C     Opcode = vm.OP_PUSH3_C
	OP_PUSH3       Opcode = vm.OP_PUSH3
	OP_PUSH3_S     Opcode = vm.OP_PUSH3_S
	OP_PUSH3_ADR   Opcode = vm.OP_PUSH3_ADR
	OP_PUSH4_C     Opcode = vm.OP_PUSH4_C
	OP_PUSH4       Opcode = vm.OP_PUSH4
	OP_PUSH4_S     Opcode = vm.OP_PUSH4_S
	OP_PUSH4_ADR   Opcode = vm.OP_PUSH4_ADR
	OP_PUSH5_C     Opcode = vm.OP_PUSH5_C
	OP_PUSH5       Opcode = vm.OP_PUSH5
	OP_PUSH5_S     Opcode = vm.OP_PUSH5_S
	OP_PUSH5_ADR   Opcode = vm.OP_PUSH5_ADR
	OP_LOAD_BOTH   Opcode = vm.OP_LOAD_BOTH
	OP_LOAD_S_BOTH Opcode = vm.OP_LOAD_S_BOTH
	OP_CONST       Opcode = vm.OP_CONST
	OP_CONST_S     Opcode = vm.OP_CONST_S
	OP_NUM_OPCODES Opcode = vm.OP_NUM_OPCODES
)

// OpcodeInfo describes an instruction's encoded name and operands.
type OpcodeInfo struct {
	Name       string
	ParamCount int
	CaseTable  bool
}

func OpcodeMetadata(op Opcode) (OpcodeInfo, bool) {
	info, ok := op.Info()
	return OpcodeInfo{Name: info.Name, ParamCount: info.ParamCount, CaseTable: info.CaseTable}, ok
}
