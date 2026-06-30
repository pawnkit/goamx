package vm

import "testing"

func TestOpcodeMetadataCoversKnownOpcodes(t *testing.T) {
	for op := OP_NONE + 1; op < OP_NUM_OPCODES; op++ {
		info, ok := op.Info()
		if !ok {
			t.Fatalf("missing opcode metadata for %d", op)
		}
		if info.Name == "" {
			t.Fatalf("opcode %d has empty name", op)
		}
	}
}

func TestOpcodeParameterCounts(t *testing.T) {
	tests := []struct {
		op    Opcode
		count int
	}{
		{OP_CONST_PRI, 1},
		{OP_ADD, 0},
		{OP_SYSREQ_N, 2},
		{OP_PUSH5_C, 5},
		{OP_LOAD_BOTH, 2},
	}
	for _, tt := range tests {
		info, ok := tt.op.Info()
		if !ok {
			t.Fatalf("missing opcode metadata for %d", tt.op)
		}
		if info.ParamCount != tt.count {
			t.Fatalf("%s ParamCount = %d, want %d", info.Name, info.ParamCount, tt.count)
		}
	}
}
