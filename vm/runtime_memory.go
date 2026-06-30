package vm

import "fmt"

func (vm *VM) Allot(cells int) (Cell, error) {
	if cells < 0 {
		return 0, fmt.Errorf("%w: negative allocation %d", ErrInvalidMemoryAccess, cells)
	}
	bytes := Cell(cells * cellBytes)
	if cells != 0 && int(bytes)/cellBytes != cells {
		return 0, fmt.Errorf("%w: allocation overflow", ErrInvalidMemoryAccess)
	}
	if vm.stk-vm.hea-bytes < stackMargin {
		return 0, fmt.Errorf("%w: heap/stack collision", ErrInvalidMemoryAccess)
	}
	addr := vm.hea
	vm.hea += bytes
	return addr, nil
}

func (vm *VM) Release(addr Cell) error {
	base := Cell(vm.header.HEA - vm.header.DAT)
	if addr < base || addr > vm.hea || addr%cellBytes != 0 {
		return fmt.Errorf("%w: invalid heap release address %d", ErrInvalidMemoryAccess, addr)
	}
	vm.hea = addr
	return nil
}

func (vm *VM) ReadString(addr Cell) (string, error) {
	return readString(vm.memory, addr)
}

func (vm *VM) WriteString(addr Cell, value string) error {
	return writeString(vm.memory, addr, value)
}

func (vm *VM) WriteStringN(addr Cell, value string, maxCells int, packed bool) error {
	return writeStringN(vm.memory, addr, value, maxCells, packed)
}

func (vm *VM) ReadCell(addr Cell) (Cell, error) {
	return readCell(vm.memory, addr)
}

func (vm *VM) WriteCell(addr Cell, value Cell) error {
	return writeCell(vm.memory, addr, value)
}

func (vm *VM) ReadBytes(addr Cell, size int) ([]byte, error) {
	start, end, err := checkedByteRange(vm.memory, addr, Cell(size))
	if err != nil {
		return nil, err
	}
	return append([]byte(nil), vm.memory[start:end]...), nil
}

func (vm *VM) WriteBytes(addr Cell, value []byte) error {
	start, end, err := checkedByteRange(vm.memory, addr, Cell(len(value)))
	if err != nil {
		return err
	}
	copy(vm.memory[start:end], value)
	return nil
}
