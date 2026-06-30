package vm

import (
	"encoding/binary"
	"fmt"
	"unicode/utf8"
)

const (
	cellBytes   = 4
	unpackedMax = 0x00ffffff
)

func initialMemory(image []byte, h header) ([]byte, error) {
	if h.DAT > h.Size || h.Size > uint32(len(image)) {
		return nil, fmt.Errorf("%w: invalid data segment bounds", ErrInvalidAMX)
	}
	dataImage := image[h.DAT:h.Size]
	if h.STP < h.DAT {
		return nil, fmt.Errorf("%w: stack top precedes data segment", ErrInvalidAMX)
	}
	size := int(h.STP - h.DAT)
	if size < len(dataImage) {
		size = len(dataImage)
	}
	memory := make([]byte, size)
	copy(memory, dataImage)
	return memory, nil
}

func readCell(memory []byte, addr Cell) (Cell, error) {
	off, err := checkedCellOffset(memory, addr)
	if err != nil {
		return 0, err
	}
	return Cell(int32(binary.LittleEndian.Uint32(memory[off : off+cellBytes]))), nil
}

func writeCell(memory []byte, addr, value Cell) error {
	off, err := checkedCellOffset(memory, addr)
	if err != nil {
		return err
	}
	binary.LittleEndian.PutUint32(memory[off:off+cellBytes], uint32(int32(value)))
	return nil
}

func readSized(memory []byte, addr, size Cell) (Cell, error) {
	start, end, err := checkedByteRange(memory, addr, size)
	if err != nil {
		return 0, err
	}
	switch size {
	case 1:
		return Cell(memory[start]), nil
	case 2:
		return Cell(binary.LittleEndian.Uint16(memory[start:end])), nil
	case 4:
		return Cell(int32(binary.LittleEndian.Uint32(memory[start:end]))), nil
	default:
		return 0, fmt.Errorf("%w: unsupported byte access size %d", ErrInvalidMemoryAccess, size)
	}
}

func writeSized(memory []byte, addr, size, value Cell) error {
	start, end, err := checkedByteRange(memory, addr, size)
	if err != nil {
		return err
	}
	switch size {
	case 1:
		memory[start] = byte(value)
	case 2:
		binary.LittleEndian.PutUint16(memory[start:end], uint16(value))
	case 4:
		binary.LittleEndian.PutUint32(memory[start:end], uint32(int32(value)))
	default:
		return fmt.Errorf("%w: unsupported byte access size %d", ErrInvalidMemoryAccess, size)
	}
	return nil
}

func readString(memory []byte, addr Cell) (string, error) {
	first, err := readCell(memory, addr)
	if err != nil {
		return "", err
	}
	if uint32(first) > unpackedMax {
		return readPackedString(memory, addr)
	}
	return readUnpackedString(memory, addr)
}

func readUnpackedString(memory []byte, addr Cell) (string, error) {
	var out []rune
	for {
		cell, err := readCell(memory, addr)
		if err != nil {
			return "", err
		}
		if cell == 0 {
			return string(out), nil
		}
		if cell < 0 || cell > utf8.MaxRune {
			return "", fmt.Errorf("%w: invalid string rune %d", ErrInvalidMemoryAccess, cell)
		}
		out = append(out, rune(cell))
		addr += cellBytes
	}
}

func readPackedString(memory []byte, addr Cell) (string, error) {
	var out []byte
	for {
		cell, err := readCell(memory, addr)
		if err != nil {
			return "", err
		}
		raw := uint32(cell)
		for shift := 24; shift >= 0; shift -= 8 {
			ch := byte(raw >> uint(shift))
			if ch == 0 {
				return string(out), nil
			}
			out = append(out, ch)
		}
		addr += cellBytes
	}
}

func writeString(memory []byte, addr Cell, value string) error {
	for _, r := range value {
		if err := writeCell(memory, addr, Cell(r)); err != nil {
			return err
		}
		addr += cellBytes
	}
	return writeCell(memory, addr, 0)
}

func writeStringN(memory []byte, addr Cell, value string, maxCells int, packed bool) error {
	if maxCells <= 0 {
		return fmt.Errorf("%w: string capacity must be positive", ErrInvalidMemoryAccess)
	}
	if !packed {
		runes := []rune(value)
		if len(runes) >= maxCells {
			runes = runes[:maxCells-1]
		}
		for _, r := range runes {
			if err := writeCell(memory, addr, Cell(r)); err != nil {
				return err
			}
			addr += cellBytes
		}
		return writeCell(memory, addr, 0)
	}
	valueBytes := []byte(value)
	maxBytes := maxCells*cellBytes - 1
	if len(valueBytes) > maxBytes {
		valueBytes = valueBytes[:maxBytes]
	}
	for cellIndex := 0; cellIndex < maxCells; cellIndex++ {
		var raw uint32
		for byteIndex := 0; byteIndex < cellBytes; byteIndex++ {
			index := cellIndex*cellBytes + byteIndex
			if index < len(valueBytes) {
				raw |= uint32(valueBytes[index]) << uint(24-byteIndex*8)
			}
		}
		if err := writeCell(memory, addr+Cell(cellIndex*cellBytes), Cell(int32(raw))); err != nil {
			return err
		}
		if (cellIndex+1)*cellBytes >= len(valueBytes)+1 {
			return nil
		}
	}
	return nil
}

func checkedCellOffset(memory []byte, addr Cell) (int, error) {
	if addr < 0 {
		return 0, fmt.Errorf("%w: negative address %d", ErrInvalidMemoryAccess, addr)
	}
	if addr%cellBytes != 0 {
		return 0, fmt.Errorf("%w: unaligned cell address %d", ErrInvalidMemoryAccess, addr)
	}
	off := int(addr)
	if off > len(memory)-cellBytes {
		return 0, fmt.Errorf("%w: address %d outside data memory", ErrInvalidMemoryAccess, addr)
	}
	return off, nil
}

func checkedByteRange(memory []byte, addr, size Cell) (int, int, error) {
	if addr < 0 {
		return 0, 0, fmt.Errorf("%w: negative address %d", ErrInvalidMemoryAccess, addr)
	}
	if size < 0 {
		return 0, 0, fmt.Errorf("%w: negative range size %d", ErrInvalidMemoryAccess, size)
	}
	start := int(addr)
	end := start + int(size)
	if start < 0 || end < start || end > len(memory) {
		return 0, 0, fmt.Errorf("%w: range %d..%d outside data memory", ErrInvalidMemoryAccess, addr, addr+size)
	}
	return start, end, nil
}
