package vm

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestExecPublicArithmetic(t *testing.T) {
	code := appendInstr(nil, OP_CONST_PRI, 2)
	code = appendInstr(code, OP_CONST_ALT, 3)
	code = appendInstr(code, OP_ADD)
	code = appendInstr(code, OP_SMUL_C, 4)
	code = appendInstr(code, OP_HALT, 0)

	vm, err := LoadBytes("arith.amx", amxWithCode(t, code))
	if err != nil {
		t.Fatal(err)
	}
	got, err := vm.ExecPublic(0)
	if err != nil {
		t.Fatal(err)
	}
	if got != 20 {
		t.Fatalf("ExecPublic() = %d, want 20", got)
	}
}

func TestExecPublicAcceptsAlignmentAndDebugOpcodes(t *testing.T) {
	code := appendInstr(nil, OP_CONST_PRI, 42)
	code = appendInstr(code, OP_ALIGN_PRI, 1)
	code = appendInstr(code, OP_ALIGN_ALT, 1)
	code = appendInstr(code, OP_FILE, 0)
	code = appendInstr(code, OP_LINE, 12)
	code = appendInstr(code, OP_SYMBOL, 0)
	code = appendInstr(code, OP_SRANGE, 0)
	code = appendInstr(code, OP_SYMTAG, 0)
	code = appendInstr(code, OP_HALT, 0)

	vm, err := LoadBytes("debug.amx", amxWithCode(t, code))
	if err != nil {
		t.Fatal(err)
	}
	got, err := vm.ExecPublic(0)
	if err != nil {
		t.Fatal(err)
	}
	if got != 42 {
		t.Fatalf("ExecPublic() = %d, want 42", got)
	}
}

func TestExecPublicBranch(t *testing.T) {
	code := appendInstr(nil, OP_CONST_PRI, 0)
	code = appendInstr(code, OP_JZER, 32)
	code = appendInstr(code, OP_CONST_PRI, 1)
	code = appendInstr(code, OP_HALT, 0)
	code = appendInstr(code, OP_CONST_PRI, 9)
	code = appendInstr(code, OP_HALT, 0)

	vm, err := LoadBytes("branch.amx", amxWithCode(t, code))
	if err != nil {
		t.Fatal(err)
	}
	got, err := vm.ExecPublic(0)
	if err != nil {
		t.Fatal(err)
	}
	if got != 9 {
		t.Fatalf("ExecPublic() = %d, want 9", got)
	}
}

func TestExecPublicMemoryOps(t *testing.T) {
	code := appendInstr(nil, OP_CONST, 0, 40)
	code = appendInstr(code, OP_INC, 0)
	code = appendInstr(code, OP_LOAD_PRI, 0)
	code = appendInstr(code, OP_ADD_C, 1)
	code = appendInstr(code, OP_STOR_PRI, 4)
	code = appendInstr(code, OP_LOAD_PRI, 4)
	code = appendInstr(code, OP_HALT, 0)

	vm, err := LoadBytes("memory.amx", amxWithCode(t, code))
	if err != nil {
		t.Fatal(err)
	}
	got, err := vm.ExecPublic(0)
	if err != nil {
		t.Fatal(err)
	}
	if got != 42 {
		t.Fatalf("ExecPublic() = %d, want 42", got)
	}
}

func TestExecPublicCallAndReturn(t *testing.T) {
	code := appendInstr(nil, OP_CALL, 24)
	code = appendInstr(code, OP_ADD_C, 1)
	code = appendInstr(code, OP_HALT, 0)
	code = appendInstr(code, OP_PROC)
	code = appendInstr(code, OP_CONST_PRI, 41)
	code = appendInstr(code, OP_RET)

	vm, err := LoadBytes("call.amx", amxWithCode(t, code))
	if err != nil {
		t.Fatal(err)
	}
	got, err := vm.ExecPublic(0)
	if err != nil {
		t.Fatal(err)
	}
	if got != 42 {
		t.Fatalf("ExecPublic() = %d, want 42", got)
	}
}

func TestExecPublicStackFrameLocals(t *testing.T) {
	code := appendInstr(nil, OP_PROC)
	code = appendInstr(code, OP_STACK, -4)
	code = appendInstr(code, OP_CONST_PRI, 40)
	code = appendInstr(code, OP_STOR_S_PRI, -4)
	code = appendInstr(code, OP_INC_S, -4)
	code = appendInstr(code, OP_LOAD_S_PRI, -4)
	code = appendInstr(code, OP_ADD_C, 1)
	code = appendInstr(code, OP_STACK, 4)
	code = appendInstr(code, OP_RETN)

	vm, err := LoadBytes("locals.amx", amxWithCode(t, code))
	if err != nil {
		t.Fatal(err)
	}
	got, err := vm.ExecPublic(0)
	if err != nil {
		t.Fatal(err)
	}
	if got != 42 {
		t.Fatalf("ExecPublic() = %d, want 42", got)
	}
}

func TestExecPublicBitwiseShiftAndDivision(t *testing.T) {
	code := appendInstr(nil, OP_CONST_PRI, 7)
	code = appendInstr(code, OP_CONST_ALT, 3)
	code = appendInstr(code, OP_AND)
	code = appendInstr(code, OP_SHL_C_PRI, 4)
	code = appendInstr(code, OP_CONST_ALT, 2)
	code = appendInstr(code, OP_SDIV)
	code = appendInstr(code, OP_HALT, 0)

	vm, err := LoadBytes("ops.amx", amxWithCode(t, code))
	if err != nil {
		t.Fatal(err)
	}
	got, err := vm.ExecPublic(0)
	if err != nil {
		t.Fatal(err)
	}
	if got != 24 {
		t.Fatalf("ExecPublic() = %d, want 24", got)
	}
}

func TestExecPublicSignedBranch(t *testing.T) {
	code := appendInstr(nil, OP_CONST_PRI, -1)
	code = appendInstr(code, OP_CONST_ALT, 1)
	code = appendInstr(code, OP_JSLESS, 40)
	code = appendInstr(code, OP_CONST_PRI, 1)
	code = appendInstr(code, OP_HALT, 0)
	code = appendInstr(code, OP_CONST_PRI, 42)
	code = appendInstr(code, OP_HALT, 0)

	vm, err := LoadBytes("branch.amx", amxWithCode(t, code))
	if err != nil {
		t.Fatal(err)
	}
	got, err := vm.ExecPublic(0)
	if err != nil {
		t.Fatal(err)
	}
	if got != 42 {
		t.Fatalf("ExecPublic() = %d, want 42", got)
	}
}

func TestExecPublicArrayIndexing(t *testing.T) {
	code := appendInstr(nil, OP_CONST, 0, 11)
	code = appendInstr(code, OP_CONST, 4, 22)
	code = appendInstr(code, OP_CONST_ALT, 0)
	code = appendInstr(code, OP_CONST_PRI, 1)
	code = appendInstr(code, OP_IDXADDR)
	code = appendInstr(code, OP_LOAD_I)
	code = appendInstr(code, OP_HALT, 0)

	vm, err := LoadBytes("array.amx", amxWithCode(t, code))
	if err != nil {
		t.Fatal(err)
	}
	got, err := vm.ExecPublic(0)
	if err != nil {
		t.Fatal(err)
	}
	if got != 22 {
		t.Fatalf("ExecPublic() = %d, want 22", got)
	}
}

func TestExecPublicLoadBoth(t *testing.T) {
	code := appendInstr(nil, OP_CONST, 0, 10)
	code = appendInstr(code, OP_CONST, 4, 32)
	code = appendInstr(code, OP_LOAD_BOTH, 0, 4)
	code = appendInstr(code, OP_ADD)
	code = appendInstr(code, OP_HALT, 0)

	vm, err := LoadBytes("loadboth.amx", amxWithCode(t, code))
	if err != nil {
		t.Fatal(err)
	}
	got, err := vm.ExecPublic(0)
	if err != nil {
		t.Fatal(err)
	}
	if got != 42 {
		t.Fatalf("ExecPublic() = %d, want 42", got)
	}
}

func TestExecPublicMacroPushCallsNative(t *testing.T) {
	code := appendInstr(nil, OP_PUSH2_C, 5, 4)
	code = appendInstr(code, OP_SYSREQ_N, 0, 8)
	code = appendInstr(code, OP_HALT, 0)

	vm, err := LoadBytes("macro-native.amx", amxWithCodeAndNatives(t, code, []string{"add"}))
	if err != nil {
		t.Fatal(err)
	}
	if err := vm.RegisterNative("add", func(ctx NativeContext, params []Cell) (Cell, error) {
		if diff := cmp.Diff([]Cell{4, 5}, params); diff != "" {
			t.Fatalf("native params mismatch (-want +got):\n%s", diff)
		}
		return params[0] + params[1], nil
	}); err != nil {
		t.Fatal(err)
	}
	got, err := vm.ExecPublic(0)
	if err != nil {
		t.Fatal(err)
	}
	if got != 9 {
		t.Fatalf("ExecPublic() = %d, want 9", got)
	}
}

func TestExecPublicBlockMemoryOps(t *testing.T) {
	code := appendInstr(nil, OP_CONST, 0, 42)
	code = appendInstr(code, OP_CONST_PRI, 7)
	code = appendInstr(code, OP_CONST_ALT, 8)
	code = appendInstr(code, OP_FILL, 8)
	code = appendInstr(code, OP_CONST_PRI, 0)
	code = appendInstr(code, OP_CONST_ALT, 16)
	code = appendInstr(code, OP_MOVS, 4)
	code = appendInstr(code, OP_CONST_PRI, 0)
	code = appendInstr(code, OP_CONST_ALT, 16)
	code = appendInstr(code, OP_CMPS, 4)
	code = appendInstr(code, OP_EQ_C_PRI, 0)
	code = appendInstr(code, OP_HALT, 0)

	vm, err := LoadBytes("blocks.amx", amxWithCode(t, code))
	if err != nil {
		t.Fatal(err)
	}
	got, err := vm.ExecPublic(0)
	if err != nil {
		t.Fatal(err)
	}
	if got != 1 {
		t.Fatalf("ExecPublic() = %d, want 1", got)
	}
	cell, err := vm.ReadCell(8)
	if err != nil {
		t.Fatal(err)
	}
	if cell != 7 {
		t.Fatalf("filled cell = %d, want 7", cell)
	}
}

func TestExecPublicIndirectRefsAndByteAccess(t *testing.T) {
	code := appendInstr(nil, OP_CONST, 0, 8)
	code = appendInstr(code, OP_CONST_PRI, 42)
	code = appendInstr(code, OP_SREF_PRI, 0)
	code = appendInstr(code, OP_LREF_ALT, 0)
	code = appendInstr(code, OP_CONST_PRI, 0x1234)
	code = appendInstr(code, OP_CONST_ALT, 12)
	code = appendInstr(code, OP_STRB_I, 2)
	code = appendInstr(code, OP_CONST_PRI, 12)
	code = appendInstr(code, OP_LODB_I, 2)
	code = appendInstr(code, OP_LREF_ALT, 0)
	code = appendInstr(code, OP_ADD)
	code = appendInstr(code, OP_HALT, 0)

	vm, err := LoadBytes("refs.amx", amxWithCode(t, code))
	if err != nil {
		t.Fatal(err)
	}
	got, err := vm.ExecPublic(0)
	if err != nil {
		t.Fatal(err)
	}
	if got != 0x1234+42 {
		t.Fatalf("ExecPublic() = %d, want %d", got, 0x1234+42)
	}
}

func TestExecPublicHeapRepeatedPushSwapAndSysreqPRI(t *testing.T) {
	code := appendInstr(nil, OP_HEAP, 8)
	code = appendInstr(code, OP_CONST_PRI, 3)
	code = appendInstr(code, OP_PUSH_R, 2)
	code = appendInstr(code, OP_SWAP_PRI)
	code = appendInstr(code, OP_CONST_ALT, 3)
	code = appendInstr(code, OP_SWAP_ALT)
	code = appendInstr(code, OP_PUSH_C, 8)
	code = appendInstr(code, OP_CONST_PRI, 0)
	code = appendInstr(code, OP_SYSREQ_PRI)
	code = appendInstr(code, OP_HALT, 0)

	vm, err := LoadBytes("heap-sysreq-pri.amx", amxWithCodeAndNatives(t, code, []string{"sum"}))
	if err != nil {
		t.Fatal(err)
	}
	if err := vm.RegisterNative("sum", func(ctx NativeContext, params []Cell) (Cell, error) {
		if diff := cmp.Diff([]Cell{3, 3}, params); diff != "" {
			t.Fatalf("native params mismatch (-want +got):\n%s", diff)
		}
		return params[0] + params[1], nil
	}); err != nil {
		t.Fatal(err)
	}
	got, err := vm.ExecPublic(0)
	if err != nil {
		t.Fatal(err)
	}
	if got != 6 {
		t.Fatalf("ExecPublic() = %d, want 6", got)
	}
}

func TestExecPublicUnsignedDivisionAndSign(t *testing.T) {
	code := appendInstr(nil, OP_CONST_PRI, -1)
	code = appendInstr(code, OP_CONST_ALT, 2)
	code = appendInstr(code, OP_UDIV)
	code = appendInstr(code, OP_CONST_ALT, 0x80)
	code = appendInstr(code, OP_SIGN_ALT)
	code = appendInstr(code, OP_ADD)
	code = appendInstr(code, OP_HALT, 0)

	vm, err := LoadBytes("unsigned-sign.amx", amxWithCode(t, code))
	if err != nil {
		t.Fatal(err)
	}
	got, err := vm.ExecPublic(0)
	if err != nil {
		t.Fatal(err)
	}
	want := Cell(int32(uint32(0xffffffff)/2)) + Cell(-128)
	if got != want {
		t.Fatalf("ExecPublic() = %d, want %d", got, want)
	}
}

func TestExecPublicSwitch(t *testing.T) {
	code := appendInstr(nil, OP_CONST_PRI, 7)
	code = appendInstr(code, OP_SWITCH, 52)
	code = appendInstr(code, OP_CONST_PRI, 1)
	code = appendInstr(code, OP_HALT, 0)
	code = appendInstr(code, OP_CONST_PRI, 42)
	code = appendInstr(code, OP_HALT, 0)
	code = appendInstr(code, OP_NOP)
	code = appendInstr(code, OP_CASETBL, 2, 16, 7, 32)

	vm, err := LoadBytes("switch.amx", amxWithCode(t, code))
	if err != nil {
		t.Fatal(err)
	}
	got, err := vm.ExecPublic(0)
	if err != nil {
		t.Fatal(err)
	}
	if got != 42 {
		t.Fatalf("ExecPublic() = %d, want 42", got)
	}
}

func TestExecPublicBoundsError(t *testing.T) {
	code := appendInstr(nil, OP_CONST_PRI, 3)
	code = appendInstr(code, OP_BOUNDS, 2)
	code = appendInstr(code, OP_HALT, 0)

	vm, err := LoadBytes("bounds.amx", amxWithCode(t, code))
	if err != nil {
		t.Fatal(err)
	}
	var runtimeErr RuntimeError
	if _, err := vm.ExecPublic(0); !errors.As(err, &runtimeErr) {
		t.Fatal("expected bounds error")
	} else if runtimeErr.Code != RuntimeErrorBounds {
		t.Fatalf("runtime error code = %s, want %s", runtimeErr.Code, RuntimeErrorBounds)
	}
}

func TestExecPublicDivideByZeroError(t *testing.T) {
	code := appendInstr(nil, OP_CONST_PRI, 1)
	code = appendInstr(code, OP_CONST_ALT, 0)
	code = appendInstr(code, OP_SDIV)
	code = appendInstr(code, OP_HALT, 0)

	vm, err := LoadBytes("divide.amx", amxWithCode(t, code))
	if err != nil {
		t.Fatal(err)
	}
	var runtimeErr RuntimeError
	if _, err := vm.ExecPublic(0); !errors.As(err, &runtimeErr) {
		t.Fatal("expected divide-by-zero runtime error")
	} else if runtimeErr.Code != RuntimeErrorDivideByZero {
		t.Fatalf("runtime error code = %s, want %s", runtimeErr.Code, RuntimeErrorDivideByZero)
	}
}

func TestExecPublicCallsNative(t *testing.T) {
	code := appendInstr(nil, OP_PUSH_C, 5)
	code = appendInstr(code, OP_PUSH_C, 4)
	code = appendInstr(code, OP_SYSREQ_N, 0, 8)
	code = appendInstr(code, OP_HALT, 0)

	vm, err := LoadBytes("native.amx", amxWithCodeAndNatives(t, code, []string{"add"}))
	if err != nil {
		t.Fatal(err)
	}
	if err := vm.RegisterNative("add", func(ctx NativeContext, params []Cell) (Cell, error) {
		if diff := cmp.Diff([]Cell{4, 5}, params); diff != "" {
			t.Fatalf("native params mismatch (-want +got):\n%s", diff)
		}
		return params[0] + params[1], nil
	}); err != nil {
		t.Fatal(err)
	}
	got, err := vm.ExecPublic(0)
	if err != nil {
		t.Fatal(err)
	}
	if got != 9 {
		t.Fatalf("ExecPublic() = %d, want 9", got)
	}
}

func TestExecPublicRejectsUnregisteredNative(t *testing.T) {
	code := appendInstr(nil, OP_SYSREQ_N, 0, 0)
	vm, err := LoadBytes("native.amx", amxWithCodeAndNatives(t, code, []string{"missing"}))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := vm.ExecPublic(0); err == nil {
		t.Fatal("expected unregistered native error")
	}
}
