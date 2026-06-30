package vm

import (
	"errors"
	"testing"
)

func TestExecContinueAfterSleep(t *testing.T) {
	code := appendInstr(nil, OP_CONST_PRI, 9)
	code = appendInstr(code, OP_HALT, 12)
	code = appendInstr(code, OP_CONST_PRI, 42)
	code = appendInstr(code, OP_HALT, 0)
	loaded, err := LoadBytes("sleep.amx", amxWithCode(t, code))
	if err != nil {
		t.Fatal(err)
	}
	vm := loaded
	if value, err := vm.ExecPublic(0); err == nil {
		t.Fatal("ExecPublic() did not report sleep")
	} else {
		if value != 9 {
			t.Fatalf("sleep value = %d, want current PRI 9", value)
		}
		var runtimeErr RuntimeError
		if !errors.As(err, &runtimeErr) || runtimeErr.Code != RuntimeErrorSleep {
			t.Fatalf("ExecPublic() error = %v, want sleep", err)
		}
	}
	if !vm.Suspended() {
		t.Fatal("VM is not suspended")
	}
	value, err := vm.Continue()
	if err != nil {
		t.Fatal(err)
	}
	if value != 42 {
		t.Fatalf("Continue() = %d, want 42", value)
	}
	if vm.Suspended() {
		t.Fatal("VM remained suspended")
	}
	if vm.stk != vm.stp || vm.hea != Cell(vm.header.HEA-vm.header.DAT) {
		t.Fatalf("execution frame was not restored: hea=%d stk=%d stp=%d", vm.hea, vm.stk, vm.stp)
	}
	if _, err := vm.Continue(); !errors.Is(err, ErrNotSleeping) {
		t.Fatalf("second Continue() error = %v, want ErrNotSleeping", err)
	}
}
