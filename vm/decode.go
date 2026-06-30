package vm

import (
	"encoding/binary"
	"fmt"
)

type Instruction struct {
	Offset int32
	Opcode Opcode
	Params []Cell
	Size   int
}

func Decode(code []byte) ([]Instruction, error) {
	var out []Instruction
	for off := 0; off < len(code); {
		ins, err := decodeOne(code, off)
		if err != nil {
			return nil, err
		}
		out = append(out, ins)
		off += ins.Size
	}
	return out, nil
}

func (vm *VM) Decode() ([]Instruction, error) {
	return Decode(vm.code())
}

func (vm *VM) code() []byte {
	if vm.header.COD >= vm.header.DAT || int(vm.header.DAT) > len(vm.image) {
		return nil
	}
	return vm.image[vm.header.COD:vm.header.DAT]
}

func decodeOne(code []byte, off int) (Instruction, error) {
	if off%cellBytes != 0 {
		return Instruction{}, fmt.Errorf("%w: unaligned offset %d", ErrInvalidInstruction, off)
	}
	if off > len(code)-cellBytes {
		return Instruction{}, fmt.Errorf("%w: truncated opcode at offset %d", ErrInvalidInstruction, off)
	}
	op := Opcode(int32(binary.LittleEndian.Uint32(code[off : off+cellBytes])))
	info, ok := op.Info()
	if !ok {
		return Instruction{}, fmt.Errorf("%w: unknown opcode %d at offset %d", ErrInvalidInstruction, op, off)
	}
	paramCount := info.ParamCount
	if info.CaseTable {
		count, err := readParam(code, off, 0)
		if err != nil {
			return Instruction{}, err
		}
		if count < 0 {
			return Instruction{}, fmt.Errorf("%w: negative casetbl count %d at offset %d", ErrInvalidInstruction, count, off)
		}
		paramCount = int(count) * 2
	}
	size := cellBytes * (1 + paramCount)
	if off > len(code)-size {
		return Instruction{}, fmt.Errorf("%w: truncated %s at offset %d", ErrInvalidInstruction, info.Name, off)
	}
	params := make([]Cell, 0, paramCount)
	for i := 0; i < paramCount; i++ {
		param, err := readParam(code, off, i)
		if err != nil {
			return Instruction{}, err
		}
		params = append(params, param)
	}
	return Instruction{Offset: int32(off), Opcode: op, Params: params, Size: size}, nil
}

func readParam(code []byte, off, index int) (Cell, error) {
	pos := off + cellBytes*(index+1)
	if pos > len(code)-cellBytes {
		return 0, fmt.Errorf("%w: truncated parameter at offset %d", ErrInvalidInstruction, off)
	}
	return Cell(int32(binary.LittleEndian.Uint32(code[pos : pos+cellBytes]))), nil
}
