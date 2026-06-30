package vm

import "fmt"

func (vm *VM) push(value Cell) error {
	vm.stk -= cellBytes
	if err := vm.checkStack(); err != nil {
		return err
	}
	return vm.WriteCell(vm.stk, value)
}

func (vm *VM) pushParams(params []Cell) error {
	for _, param := range params {
		if err := vm.push(param); err != nil {
			return err
		}
	}
	return nil
}

func (vm *VM) pop() (Cell, error) {
	if vm.stk >= vm.stp {
		return 0, fmt.Errorf("%w: stack underflow", ErrInvalidMemoryAccess)
	}
	value, err := vm.ReadCell(vm.stk)
	if err != nil {
		return 0, err
	}
	vm.stk += cellBytes
	return value, nil
}

func (vm *VM) checkStack() error {
	if vm.stk < 0 || vm.stk > vm.stp {
		return fmt.Errorf("%w: stack pointer %d outside stack bounds", ErrInvalidMemoryAccess, vm.stk)
	}
	if vm.stk-vm.hea < stackMargin {
		return fmt.Errorf("%w: heap/stack collision", ErrInvalidMemoryAccess)
	}
	return nil
}

func (vm *VM) checkHeap() error {
	if vm.hea < 0 || vm.hea > vm.stp {
		return fmt.Errorf("%w: heap pointer %d outside memory bounds", ErrInvalidMemoryAccess, vm.hea)
	}
	if vm.stk-vm.hea < stackMargin {
		return fmt.Errorf("%w: heap/stack collision", ErrInvalidMemoryAccess)
	}
	return nil
}
