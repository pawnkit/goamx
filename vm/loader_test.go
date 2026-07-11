package vm

import (
	"encoding/binary"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestLoadBytesPublicsAndNatives(t *testing.T) {
	data := minimalAMX(t)
	vm, err := LoadBytes("test.amx", data)
	if err != nil {
		t.Fatal(err)
	}
	publics, err := vm.Publics()
	if err != nil {
		t.Fatal(err)
	}
	wantPublics := []Public{{Index: 0, Name: "test_addition"}, {Index: 1, Name: "helper"}}
	if diff := cmp.Diff(wantPublics, publics); diff != "" {
		t.Fatalf("publics mismatch (-want +got):\n%s", diff)
	}
	natives, err := vm.Natives()
	if err != nil {
		t.Fatal(err)
	}
	if diff := cmp.Diff([]string{"test_native"}, natives); diff != "" {
		t.Fatalf("natives mismatch (-want +got):\n%s", diff)
	}
}

func TestLoadBytesPubVarsAndTags(t *testing.T) {
	vm, err := LoadBytes("tables.amx", amxWithPubVarsAndTags(t))
	if err != nil {
		t.Fatal(err)
	}
	pubvars, err := vm.PubVars()
	if err != nil {
		t.Fatal(err)
	}
	if diff := cmp.Diff([]Public{{Index: 0, Name: "global_value"}}, pubvars); diff != "" {
		t.Fatalf("pubvars mismatch (-want +got):\n%s", diff)
	}
	tags, err := vm.Tags()
	if err != nil {
		t.Fatal(err)
	}
	if diff := cmp.Diff([]Public{{Index: 0, Name: "Vehicle"}}, tags); diff != "" {
		t.Fatalf("tags mismatch (-want +got):\n%s", diff)
	}
}

func TestLoadBytesRejectsBadMagic(t *testing.T) {
	data := minimalAMX(t)
	binary.LittleEndian.PutUint16(data[4:6], 0)
	if _, err := LoadBytes("bad.amx", data); err == nil {
		t.Fatal("expected error")
	}
}

func TestLoadBytesAllowsCompactFlag(t *testing.T) {
	data := minimalAMX(t)
	binary.LittleEndian.PutUint16(data[8:10], 0x04)
	if _, err := LoadBytes("compact.amx", data); err != nil {
		t.Fatal(err)
	}
}

func TestLoadBytesSupportsVersionNine(t *testing.T) {
	data := minimalAMX(t)
	data[6], data[7] = 9, 9
	if _, err := LoadBytes("version-nine.amx", data); err != nil {
		t.Fatal(err)
	}
}

func TestLoadBytesSupportsLegacyInlineNames(t *testing.T) {
	const header = 56
	publics, natives := uint32(header), uint32(header+24)
	data := make([]byte, natives)
	binary.LittleEndian.PutUint32(data[publics:publics+4], 0)
	copy(data[publics+4:publics+24], "legacy_public")
	cod := uint32(len(data))
	data = appendInstr(data, OP_HALT, 0)
	dat := uint32(len(data))
	binary.LittleEndian.PutUint32(data[0:4], dat)
	binary.LittleEndian.PutUint16(data[4:6], amxMagic32)
	data[6], data[7] = 6, 6
	binary.LittleEndian.PutUint16(data[10:12], 24)
	binary.LittleEndian.PutUint32(data[12:16], cod)
	binary.LittleEndian.PutUint32(data[16:20], dat)
	binary.LittleEndian.PutUint32(data[20:24], dat)
	binary.LittleEndian.PutUint32(data[24:28], dat+256)
	binary.LittleEndian.PutUint32(data[32:36], publics)
	for offset := 36; offset <= 52; offset += 4 {
		binary.LittleEndian.PutUint32(data[offset:offset+4], natives)
	}
	loaded, err := LoadBytes("legacy.amx", data)
	if err != nil {
		t.Fatal(err)
	}
	publicsList, err := loaded.Publics()
	if err != nil || len(publicsList) != 1 || publicsList[0].Name != "legacy_public" {
		t.Fatalf("Publics() = %#v, %v", publicsList, err)
	}
	value, err := loaded.ExecPublic(0)
	if err != nil || value != 0 {
		t.Fatalf("ExecPublic() = %d, %v", value, err)
	}
}

func TestLoadBytesRejectsUnsupportedCellMagic(t *testing.T) {
	data := minimalAMX(t)
	binary.LittleEndian.PutUint16(data[4:6], amxMagic64)
	if _, err := LoadBytes("wide.amx", data); err == nil {
		t.Fatal("expected unsupported cell size error")
	}
}

func TestVMMemoryHelpers(t *testing.T) {
	vm, err := LoadBytes("test.amx", minimalAMX(t))
	if err != nil {
		t.Fatal(err)
	}
	ctx := vm

	if err := ctx.WriteCell(0, 42); err != nil {
		t.Fatal(err)
	}
	got, err := ctx.ReadCell(0)
	if err != nil {
		t.Fatal(err)
	}
	if got != 42 {
		t.Fatalf("ReadCell() = %d, want 42", got)
	}

	if err := ctx.WriteString(4, "ok"); err != nil {
		t.Fatal(err)
	}
	text, err := ctx.ReadString(4)
	if err != nil {
		t.Fatal(err)
	}
	if text != "ok" {
		t.Fatalf("ReadString() = %q, want %q", text, "ok")
	}
}

func TestVMConvertsFileStackTopToDataRelativeAddress(t *testing.T) {
	loaded, err := LoadBytes("stack.amx", amxWithCode(t, appendInstr(nil, OP_HALT, 0)))
	if err != nil {
		t.Fatal(err)
	}
	vm := loaded
	if len(vm.memory) != 256 {
		t.Fatalf("memory size = %d, want 256", len(vm.memory))
	}
	if err := vm.Reset(); err != nil {
		t.Fatal(err)
	}
	if vm.stk != 256 || vm.stp != 256 {
		t.Fatalf("stack bounds = %d..%d, want 256..256", vm.stk, vm.stp)
	}
}

func TestVMEnforcesCanonicalStackMargin(t *testing.T) {
	vm := &VM{memory: make([]byte, 256), hea: 128, stk: 128 + stackMargin - 1, stp: 256}
	if err := vm.checkStack(); err == nil {
		t.Fatal("checkStack() accepted less than the canonical margin")
	}
	vm.stk = 128 + stackMargin
	if err := vm.checkStack(); err != nil {
		t.Fatalf("checkStack() rejected exact margin: %v", err)
	}
}

func TestVMReadsPackedStrings(t *testing.T) {
	vm, err := LoadBytes("test.amx", minimalAMX(t))
	if err != nil {
		t.Fatal(err)
	}
	ctx := vm
	if err := ctx.WriteCell(0, Cell(0x68656c6c)); err != nil {
		t.Fatal(err)
	}
	if err := ctx.WriteCell(4, Cell(0x6f000000)); err != nil {
		t.Fatal(err)
	}
	got, err := ctx.ReadString(0)
	if err != nil {
		t.Fatal(err)
	}
	if got != "hello" {
		t.Fatalf("ReadString() = %q, want hello", got)
	}
}

func TestVMMemoryRejectsInvalidCellAddress(t *testing.T) {
	vm, err := LoadBytes("test.amx", minimalAMX(t))
	if err != nil {
		t.Fatal(err)
	}
	ctx := vm
	if _, err := ctx.ReadCell(1); err == nil {
		t.Fatal("expected unaligned address error")
	}
	if err := ctx.WriteCell(1<<20, 1); err == nil {
		t.Fatal("expected out-of-bounds address error")
	}
}

func TestDecodeInstructions(t *testing.T) {
	code := appendInstr(nil, OP_CONST_PRI, 7)
	code = appendInstr(code, OP_PUSH_C, 2)
	code = appendInstr(code, OP_ADD)
	code = appendInstr(code, OP_HALT, 0)

	got, err := Decode(code)
	if err != nil {
		t.Fatal(err)
	}
	want := []Instruction{
		{Offset: 0, Opcode: OP_CONST_PRI, Params: []Cell{7}, Size: 8},
		{Offset: 8, Opcode: OP_PUSH_C, Params: []Cell{2}, Size: 8},
		{Offset: 16, Opcode: OP_ADD, Params: []Cell{}, Size: 4},
		{Offset: 20, Opcode: OP_HALT, Params: []Cell{0}, Size: 8},
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("decoded instructions mismatch (-want +got):\n%s", diff)
	}
}

func TestDecodeRejectsUnknownOpcode(t *testing.T) {
	code := make([]byte, 4)
	binary.LittleEndian.PutUint32(code, uint32(OP_NUM_OPCODES+99))
	if _, err := Decode(code); err == nil {
		t.Fatal("expected unknown opcode error")
	}
}

func TestDecodeCaseTable(t *testing.T) {
	code := appendInstr(nil, OP_CASETBL, 2, 16, 7, 24)
	got, err := Decode(code)
	if err != nil {
		t.Fatal(err)
	}
	want := []Instruction{{Offset: 0, Opcode: OP_CASETBL, Params: []Cell{2, 16, 7, 24}, Size: 20}}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("decoded casetbl mismatch (-want +got):\n%s", diff)
	}
}

func TestVMDecodeUsesCodeSegment(t *testing.T) {
	vm, err := LoadBytes("code.amx", amxWithCode(t, appendInstr(nil, OP_HALT, 0)))
	if err != nil {
		t.Fatal(err)
	}
	got, err := vm.Decode()
	if err != nil {
		t.Fatal(err)
	}
	want := []Instruction{{Offset: 0, Opcode: OP_HALT, Params: []Cell{0}, Size: 8}}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("decoded VM instructions mismatch (-want +got):\n%s", diff)
	}
}
