package vm

import "fmt"

func (vm *VM) RegisterNative(name string, fn NativeFunc) error {
	vm.registered[name] = fn
	return nil
}

func (vm *VM) NativeRegistered(name string) bool { _, ok := vm.registered[name]; return ok }

func (vm *VM) ExecPublic(index int, args ...Cell) (Cell, error) {
	if index < 0 || index >= len(vm.publics) {
		return 0, ErrPublicIndexOutOfRange
	}
	return vm.execAt(vm.publics[index].Address, args...)
}

func (vm *VM) ExecMain(args ...Cell) (Cell, error) {
	if vm.header.CIP == ^uint32(0) {
		return 0, fmt.Errorf("%w: main", ErrPublicIndexOutOfRange)
	}
	return vm.execAt(vm.header.CIP, args...)
}

func (vm *VM) Continue() (Cell, error) {
	if !vm.suspended {
		return 0, ErrNotSleeping
	}
	vm.suspended = false
	return vm.execLoop(vm.resumeCIP)
}

func (vm *VM) Suspended() bool { return vm.suspended }

func (vm *VM) CallPublic(name string, args ...Cell) (Cell, error) {
	index := -1
	for i, public := range vm.publics {
		if public.Name == name {
			index = i
			break
		}
	}
	if index < 0 {
		return 0, fmt.Errorf("%w: %s", ErrPublicIndexOutOfRange, name)
	}
	pri, alt, hea, stk, stp, frm := vm.pri, vm.alt, vm.hea, vm.stk, vm.stp, vm.frm
	dynamicStart := int(hea)
	dynamicEnd := int(stp)
	if dynamicStart < 0 || dynamicEnd < dynamicStart || dynamicEnd > len(vm.memory) {
		return 0, fmt.Errorf("%w: callback memory bounds %d..%d", ErrInvalidMemoryAccess, dynamicStart, dynamicEnd)
	}
	dynamicMemory := append([]byte(nil), vm.memory[dynamicStart:dynamicEnd]...)
	defer func() {
		copy(vm.memory[dynamicStart:dynamicEnd], dynamicMemory)
		vm.pri, vm.alt, vm.hea, vm.stk, vm.stp, vm.frm = pri, alt, hea, stk, stp, frm
	}()
	return vm.execAt(vm.publics[index].Address, args...)
}
