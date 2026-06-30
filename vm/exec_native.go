package vm

import "fmt"

func (vm *VM) callNative(index, paramBytes int) (Cell, error) {
	if index < 0 || index >= len(vm.natives) {
		return 0, fmt.Errorf("%w: native index %d out of range", ErrInvalidInstruction, index)
	}
	name := vm.natives[index].Name
	fn, ok := vm.registered[name]
	if !ok {
		return 0, fmt.Errorf("%w: native %s is not registered", ErrUnsupportedExecution, name)
	}
	args, err := vm.nativeArgs(paramBytes)
	if err != nil {
		return 0, err
	}
	return fn(vm, args)
}

func (vm *VM) nativeArgs(paramBytes int) ([]Cell, error) {
	if paramBytes >= 0 {
		return vm.nativeArgsFromSysreqN(paramBytes)
	}
	return vm.nativeArgsFromStackHeader()
}

func (vm *VM) nativeArgsFromSysreqN(paramBytes int) ([]Cell, error) {
	if paramBytes%cellBytes != 0 {
		return nil, fmt.Errorf("%w: native parameter byte count %d is not cell-aligned", ErrInvalidInstruction, paramBytes)
	}
	count := paramBytes / cellBytes
	args := make([]Cell, 0, count)
	for i := 0; i < count; i++ {
		value, err := vm.ReadCell(vm.stk + Cell(i*cellBytes))
		if err != nil {
			return nil, err
		}
		args = append(args, value)
	}
	vm.stk += Cell(paramBytes)
	if err := vm.checkStack(); err != nil {
		return nil, err
	}
	return args, nil
}

func (vm *VM) nativeArgsFromStackHeader() ([]Cell, error) {
	paramCountCell, err := vm.ReadCell(vm.stk)
	if err != nil {
		return nil, err
	}
	if paramCountCell < 0 || paramCountCell%cellBytes != 0 {
		return nil, fmt.Errorf("%w: native parameter byte count %d is invalid", ErrInvalidInstruction, paramCountCell)
	}
	count := int(paramCountCell) / cellBytes
	args := make([]Cell, 0, count)
	for i := 0; i < count; i++ {
		value, err := vm.ReadCell(vm.stk + Cell((i+1)*cellBytes))
		if err != nil {
			return nil, err
		}
		args = append(args, value)
	}
	return args, nil
}
