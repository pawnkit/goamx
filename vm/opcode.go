package vm

type Opcode int32

const (
	OP_NONE Opcode = iota
	OP_LOAD_PRI
	OP_LOAD_ALT
	OP_LOAD_S_PRI
	OP_LOAD_S_ALT
	OP_LREF_PRI
	OP_LREF_ALT
	OP_LREF_S_PRI
	OP_LREF_S_ALT
	OP_LOAD_I
	OP_LODB_I
	OP_CONST_PRI
	OP_CONST_ALT
	OP_ADDR_PRI
	OP_ADDR_ALT
	OP_STOR_PRI
	OP_STOR_ALT
	OP_STOR_S_PRI
	OP_STOR_S_ALT
	OP_SREF_PRI
	OP_SREF_ALT
	OP_SREF_S_PRI
	OP_SREF_S_ALT
	OP_STOR_I
	OP_STRB_I
	OP_LIDX
	OP_LIDX_B
	OP_IDXADDR
	OP_IDXADDR_B
	OP_ALIGN_PRI
	OP_ALIGN_ALT
	OP_LCTRL
	OP_SCTRL
	OP_MOVE_PRI
	OP_MOVE_ALT
	OP_XCHG
	OP_PUSH_PRI
	OP_PUSH_ALT
	OP_PUSH_R
	OP_PUSH_C
	OP_PUSH
	OP_PUSH_S
	OP_POP_PRI
	OP_POP_ALT
	OP_STACK
	OP_HEAP
	OP_PROC
	OP_RET
	OP_RETN
	OP_CALL
	OP_CALL_PRI
	OP_JUMP
	OP_JREL
	OP_JZER
	OP_JNZ
	OP_JEQ
	OP_JNEQ
	OP_JLESS
	OP_JLEQ
	OP_JGRTR
	OP_JGEQ
	OP_JSLESS
	OP_JSLEQ
	OP_JSGRTR
	OP_JSGEQ
	OP_SHL
	OP_SHR
	OP_SSHR
	OP_SHL_C_PRI
	OP_SHL_C_ALT
	OP_SHR_C_PRI
	OP_SHR_C_ALT
	OP_SMUL
	OP_SDIV
	OP_SDIV_ALT
	OP_UMUL
	OP_UDIV
	OP_UDIV_ALT
	OP_ADD
	OP_SUB
	OP_SUB_ALT
	OP_AND
	OP_OR
	OP_XOR
	OP_NOT
	OP_NEG
	OP_INVERT
	OP_ADD_C
	OP_SMUL_C
	OP_ZERO_PRI
	OP_ZERO_ALT
	OP_ZERO
	OP_ZERO_S
	OP_SIGN_PRI
	OP_SIGN_ALT
	OP_EQ
	OP_NEQ
	OP_LESS
	OP_LEQ
	OP_GRTR
	OP_GEQ
	OP_SLESS
	OP_SLEQ
	OP_SGRTR
	OP_SGEQ
	OP_EQ_C_PRI
	OP_EQ_C_ALT
	OP_INC_PRI
	OP_INC_ALT
	OP_INC
	OP_INC_S
	OP_INC_I
	OP_DEC_PRI
	OP_DEC_ALT
	OP_DEC
	OP_DEC_S
	OP_DEC_I
	OP_MOVS
	OP_CMPS
	OP_FILL
	OP_HALT
	OP_BOUNDS
	OP_SYSREQ_PRI
	OP_SYSREQ_C
	OP_FILE
	OP_LINE
	OP_SYMBOL
	OP_SRANGE
	OP_JUMP_PRI
	OP_SWITCH
	OP_CASETBL
	OP_SWAP_PRI
	OP_SWAP_ALT
	OP_PUSH_ADR
	OP_NOP
	OP_SYSREQ_N
	OP_SYMTAG
	OP_BREAK
	OP_PUSH2_C
	OP_PUSH2
	OP_PUSH2_S
	OP_PUSH2_ADR
	OP_PUSH3_C
	OP_PUSH3
	OP_PUSH3_S
	OP_PUSH3_ADR
	OP_PUSH4_C
	OP_PUSH4
	OP_PUSH4_S
	OP_PUSH4_ADR
	OP_PUSH5_C
	OP_PUSH5
	OP_PUSH5_S
	OP_PUSH5_ADR
	OP_LOAD_BOTH
	OP_LOAD_S_BOTH
	OP_CONST
	OP_CONST_S
	_OP_SYSREQ_D
	_OP_SYSREQ_ND
	OP_NUM_OPCODES
)

type OpcodeInfo struct {
	Name       string
	ParamCount int
	CaseTable  bool
}

var opcodeInfo = map[Opcode]OpcodeInfo{}

func init() {
	names := []string{
		"none", "load.pri", "load.alt", "load.s.pri", "load.s.alt", "lref.pri", "lref.alt", "lref.s.pri", "lref.s.alt",
		"load.i", "lodb.i", "const.pri", "const.alt", "addr.pri", "addr.alt", "stor.pri", "stor.alt", "stor.s.pri",
		"stor.s.alt", "sref.pri", "sref.alt", "sref.s.pri", "sref.s.alt", "stor.i", "strb.i", "lidx", "lidx.b",
		"idxaddr", "idxaddr.b", "align.pri", "align.alt", "lctrl", "sctrl", "move.pri", "move.alt", "xchg",
		"push.pri", "push.alt", "push.r", "push.c", "push", "push.s", "pop.pri", "pop.alt", "stack", "heap",
		"proc", "ret", "retn", "call", "call.pri", "jump", "jrel", "jzer", "jnz", "jeq", "jneq", "jless", "jleq",
		"jgrtr", "jgeq", "jsless", "jsleq", "jsgrtr", "jsgeq", "shl", "shr", "sshr", "shl.c.pri", "shl.c.alt",
		"shr.c.pri", "shr.c.alt", "smul", "sdiv", "sdiv.alt", "umul", "udiv", "udiv.alt", "add", "sub", "sub.alt",
		"and", "or", "xor", "not", "neg", "invert", "add.c", "smul.c", "zero.pri", "zero.alt", "zero", "zero.s",
		"sign.pri", "sign.alt", "eq", "neq", "less", "leq", "grtr", "geq", "sless", "sleq", "sgrtr", "sgeq",
		"eq.c.pri", "eq.c.alt", "inc.pri", "inc.alt", "inc", "inc.s", "inc.i", "dec.pri", "dec.alt", "dec",
		"dec.s", "dec.i", "movs", "cmps", "fill", "halt", "bounds", "sysreq.pri", "sysreq.c", "file", "line",
		"symbol", "srange", "jump.pri", "switch", "casetbl", "swap.pri", "swap.alt", "push.adr", "nop", "sysreq.n",
		"symtag", "break", "push2.c", "push2", "push2.s", "push2.adr", "push3.c", "push3", "push3.s", "push3.adr",
		"push4.c", "push4", "push4.s", "push4.adr", "push5.c", "push5", "push5.s", "push5.adr", "load.both",
		"load.s.both", "const", "const.s", "sysreq.d", "sysreq.nd",
	}
	for op, name := range names {
		opcodeInfo[Opcode(op)] = OpcodeInfo{Name: name}
	}
	for _, op := range []Opcode{
		OP_LOAD_PRI, OP_LOAD_ALT, OP_LOAD_S_PRI, OP_LOAD_S_ALT, OP_LREF_PRI, OP_LREF_ALT, OP_LREF_S_PRI, OP_LREF_S_ALT,
		OP_LODB_I, OP_CONST_PRI, OP_CONST_ALT, OP_ADDR_PRI, OP_ADDR_ALT, OP_STOR_PRI, OP_STOR_ALT, OP_STOR_S_PRI,
		OP_STOR_S_ALT, OP_SREF_PRI, OP_SREF_ALT, OP_SREF_S_PRI, OP_SREF_S_ALT, OP_STRB_I, OP_LIDX_B, OP_IDXADDR_B,
		OP_ALIGN_PRI, OP_ALIGN_ALT, OP_LCTRL, OP_SCTRL, OP_PUSH_R, OP_PUSH_C, OP_PUSH, OP_PUSH_S, OP_STACK, OP_HEAP,
		OP_CALL, OP_JUMP, OP_JREL, OP_JZER, OP_JNZ, OP_JEQ, OP_JNEQ, OP_JLESS, OP_JLEQ, OP_JGRTR, OP_JGEQ, OP_JSLESS,
		OP_JSLEQ, OP_JSGRTR, OP_JSGEQ, OP_SHL_C_PRI, OP_SHL_C_ALT, OP_SHR_C_PRI, OP_SHR_C_ALT, OP_ADD_C, OP_SMUL_C,
		OP_ZERO, OP_ZERO_S, OP_EQ_C_PRI, OP_EQ_C_ALT, OP_INC, OP_INC_S, OP_DEC, OP_DEC_S, OP_MOVS, OP_CMPS, OP_FILL,
		OP_HALT, OP_BOUNDS, OP_SYSREQ_C, OP_FILE, OP_LINE, OP_SYMBOL, OP_SRANGE, OP_SWITCH, OP_PUSH_ADR, OP_SYMTAG,
	} {
		setParamCount(op, 1)
	}
	for _, op := range []Opcode{
		OP_PUSH2_C, OP_PUSH2, OP_PUSH2_S, OP_PUSH2_ADR,
		OP_LOAD_BOTH, OP_LOAD_S_BOTH, OP_CONST, OP_CONST_S, OP_SYSREQ_N,
	} {
		setParamCount(op, 2)
	}
	for _, op := range []Opcode{OP_PUSH3_C, OP_PUSH3, OP_PUSH3_S, OP_PUSH3_ADR} {
		setParamCount(op, 3)
	}
	for _, op := range []Opcode{OP_PUSH4_C, OP_PUSH4, OP_PUSH4_S, OP_PUSH4_ADR} {
		setParamCount(op, 4)
	}
	for _, op := range []Opcode{OP_PUSH5_C, OP_PUSH5, OP_PUSH5_S, OP_PUSH5_ADR} {
		setParamCount(op, 5)
	}
	info := opcodeInfo[OP_CASETBL]
	info.CaseTable = true
	opcodeInfo[OP_CASETBL] = info
}

func setParamCount(op Opcode, count int) {
	info := opcodeInfo[op]
	info.ParamCount = count
	opcodeInfo[op] = info
}

func (op Opcode) Info() (OpcodeInfo, bool) {
	info, ok := opcodeInfo[op]
	return info, ok && op > OP_NONE && op < OP_NUM_OPCODES
}
