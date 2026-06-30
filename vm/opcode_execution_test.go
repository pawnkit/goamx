package vm

import (
	"errors"
	"testing"
)

func TestEveryActiveOpcodeHasExecutionDispatch(t *testing.T) {
	for opcode := OP_NONE + 1; opcode < OP_NUM_OPCODES; opcode++ {
		if opcode == _OP_SYSREQ_D || opcode == _OP_SYSREQ_ND {
			continue
		}
		t.Run(opcodeName(opcode), func(t *testing.T) {
			info, ok := opcode.Info()
			if !ok {
				t.Fatalf("missing metadata for opcode %d", opcode)
			}
			params := make([]Cell, info.ParamCount)
			if opcode == OP_CASETBL {
				params = []Cell{0, 0}
			}
			vm := newOpcodeTestVM()
			_ = vm.WriteCell(0, 0)
			_ = vm.WriteCell(128, 0)
			_ = vm.WriteCell(512, 0)
			next := 4
			_, _, err := vm.execInstruction(Instruction{Opcode: opcode, Params: params}, &next)
			if errors.Is(err, ErrUnsupportedExecution) {
				t.Fatalf("active opcode %s reached unsupported execution: %v", info.Name, err)
			}
		})
	}
}

func newOpcodeTestVM() *VM {
	return &VM{
		memory:  make([]byte, 1024),
		hea:     64,
		stk:     512,
		stp:     1024,
		frm:     128,
		natives: []funcStub{{Name: "native"}},
		registered: map[string]NativeFunc{
			"native": func(NativeContext, []Cell) (Cell, error) {
				return 0, nil
			},
		},
	}
}

func opcodeName(opcode Opcode) string {
	if info, ok := opcode.Info(); ok {
		return info.Name
	}
	return "unknown"
}
