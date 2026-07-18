package amx

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/pawnkit/goamx/vm"
)

func TestRuntimePublicAPI(t *testing.T) {
	runtime, err := LoadBytes("public.amx", publicTestAMX(vm.OP_CONST_PRI, 42, vm.OP_HALT, 0))
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := runtime.Close(); err != nil {
			t.Fatal(err)
		}
	}()
	publics, err := runtime.Publics()
	if err != nil || len(publics) != 1 || publics[0].Name != "test_public" {
		t.Fatalf("Publics() = %#v, %v", publics, err)
	}
	instructions, err := runtime.Decode()
	if err != nil || len(instructions) != 2 || instructions[0].Opcode != OP_CONST_PRI {
		t.Fatalf("Decode() = %#v, %v", instructions, err)
	}
	value, err := runtime.ExecPublic(publics[0].Index)
	if err != nil || value != 42 {
		t.Fatalf("ExecPublic() = %d, %v; want 42", value, err)
	}
	if pubvars, err := runtime.PubVars(); err != nil || len(pubvars) != 0 {
		t.Fatalf("PubVars() = %#v, %v", pubvars, err)
	}
	if tags, err := runtime.Tags(); err != nil || len(tags) != 0 {
		t.Fatalf("Tags() = %#v, %v", tags, err)
	}
}

func TestLoadFileAndInvalidImage(t *testing.T) {
	path := filepath.Join(t.TempDir(), "public.amx")
	if err := os.WriteFile(path, publicTestAMX(vm.OP_HALT, 0), 0o644); err != nil {
		t.Fatal(err)
	}
	runtime, err := LoadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := runtime.Close(); err != nil {
		t.Fatal(err)
	}
	if _, err := LoadBytes("bad.amx", []byte("bad")); err == nil {
		t.Fatal("LoadBytes() accepted an invalid image")
	}
}

func TestRuntimeMemoryCloneAndReset(t *testing.T) {
	runtime, err := LoadBytes("memory.amx", publicTestAMX(vm.OP_HALT, 0))
	if err != nil {
		t.Fatal(err)
	}
	addr, err := runtime.Allot(2)
	if err != nil {
		t.Fatal(err)
	}
	if err := runtime.WriteCell(addr, 42); err != nil {
		t.Fatal(err)
	}
	if err := runtime.WriteBytes(addr+4, []byte{1, 2, 3, 4}); err != nil {
		t.Fatal(err)
	}
	clone := runtime.Clone()
	if got, err := clone.ReadCell(addr); err != nil || got != 42 {
		t.Fatalf("clone ReadCell() = %d, %v", got, err)
	}
	if err := runtime.ResetMemory(); err != nil {
		t.Fatal(err)
	}
	if got, err := runtime.ReadCell(addr); err != nil || got != 0 {
		t.Fatalf("reset ReadCell() = %d, %v", got, err)
	}
	if got, err := clone.ReadCell(addr); err != nil || got != 42 {
		t.Fatalf("clone changed after source reset: %d, %v", got, err)
	}
	if err := clone.Release(addr); err != nil {
		t.Fatal(err)
	}
}

func TestRuntimeRegistersNative(t *testing.T) {
	runtime, err := LoadBytes("native.amx", publicNativeAMX())
	if err != nil {
		t.Fatal(err)
	}
	native, ok := runtime.FindNative("increment")
	if !ok || native.Registered {
		t.Fatalf("FindNative() = %#v, %v", native, ok)
	}
	if err := runtime.RegisterNative("increment", func(ctx NativeContext, params []Cell) (Cell, error) {
		if len(params) != 1 || params[0] != 41 {
			t.Fatalf("native params = %#v", params)
		}
		return params[0] + 1, nil
	}); err != nil {
		t.Fatal(err)
	}
	value, err := runtime.ExecPublic(0)
	if err != nil || value != 42 {
		t.Fatalf("ExecPublic() = %d, %v", value, err)
	}
	if native, _ := runtime.FindNative("increment"); !native.Registered {
		t.Fatal("native not marked registered")
	}
	err = runtime.RegisterNative("missing", func(NativeContext, []Cell) (Cell, error) {
		return 0, nil
	})
	if !errors.Is(err, ErrNativeNotDeclared) {
		t.Fatalf("missing registration error = %v", err)
	}
}

func TestRuntimePackedStringsDebugHookAndUserData(t *testing.T) {
	runtime, err := LoadBytes("host.amx", publicTestAMX(vm.OP_CONST_PRI, 7, vm.OP_HALT, 0))
	if err != nil {
		t.Fatal(err)
	}
	addr, err := runtime.AllotString("packed value", true)
	if err != nil {
		t.Fatal(err)
	}
	if got, err := runtime.ReadString(addr); err != nil || got != "packed value" {
		t.Fatalf("ReadString() = %q, %v", got, err)
	}
	runtime.SetUserData(10, "host state")
	if got, ok := runtime.UserData(10); !ok || got != "host state" {
		t.Fatalf("UserData() = %#v, %v", got, ok)
	}
	var offsets []int32
	runtime.SetDebugHook(func(event DebugEvent) error {
		offsets = append(offsets, event.Instruction.Offset)
		if event.State.CIP != int(event.Instruction.Offset) {
			t.Fatalf("debug CIP = %d", event.State.CIP)
		}
		return nil
	})
	if _, err := runtime.ExecPublic(0); err != nil {
		t.Fatal(err)
	}
	if len(offsets) != 2 {
		t.Fatalf("debug hook calls = %d, want 2", len(offsets))
	}
	clone := runtime.Clone()
	if got, ok := clone.UserData(10); !ok || got != "host state" {
		t.Fatalf("clone UserData() = %#v, %v", got, ok)
	}
}

func TestRuntimeParsesDebugMetadata(t *testing.T) {
	runtime, err := LoadBytes("debug.amx", debugTestAMX())
	if err != nil {
		t.Fatal(err)
	}
	debug := runtime.DebugInfo()
	if file, ok := debug.FileAt(0); !ok || file.Name != "test.pwn" {
		t.Fatalf("FileAt() = %#v, %v", file, ok)
	}
	if line, ok := debug.LineAt(0); !ok || line.Line != 12 {
		t.Fatalf("LineAt() = %#v, %v", line, ok)
	}
	if fn, ok := debug.FunctionAt(1); !ok || fn.Name != "test_public" || len(fn.Dimensions) != 1 {
		t.Fatalf("FunctionAt() = %#v, %v", fn, ok)
	}
	if len(debug.Tags) != 1 || len(debug.Automata) != 1 || len(debug.States) != 1 {
		t.Fatalf("DebugInfo() = %#v", debug)
	}
}

func TestRuntimeInstructionLimitAndFloatCells(t *testing.T) {
	cell := CellFromFloat32(1.5)
	if got := cell.Float32(); got != 1.5 {
		t.Fatalf("Float32() = %v", got)
	}
	runtime, err := LoadBytes("loop.amx", publicTestAMX(vm.OP_JUMP, 0))
	if err != nil {
		t.Fatal(err)
	}
	runtime.SetInstructionLimit(3)
	if _, err := runtime.ExecPublic(0); !errors.Is(err, ErrUnsupportedExecution) {
		t.Fatalf("ExecPublic() error = %v, want instruction limit", err)
	}
	runtime.SetInstructionLimit(0)
}

func TestRuntimePauseAndContinue(t *testing.T) {
	runtime, err := LoadBytes("pause.amx", publicTestAMX(vm.OP_CONST_PRI, 42, vm.OP_HALT, 0))
	if err != nil {
		t.Fatal(err)
	}
	paused := false
	runtime.SetDebugHook(func(event DebugEvent) error {
		if !paused {
			paused = true
			return ErrExecutionPaused
		}
		return nil
	})
	if _, err := runtime.ExecPublic(0); !errors.Is(err, ErrExecutionPaused) {
		t.Fatalf("pause error = %v", err)
	}
	if !runtime.Suspended() || runtime.State().CIP != 0 {
		t.Fatalf("state = %+v", runtime.State())
	}
	value, err := runtime.Continue()
	if err != nil || value != 42 {
		t.Fatalf("Continue() = %d, %v", value, err)
	}
}
