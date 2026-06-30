package vm

import "fmt"

func (vm *VM) execControlInstruction(ins Instruction, next *int) (handled, halt bool, ret Cell, err error) {
	param := func(i int) Cell { return ins.Params[i] }

	switch ins.Opcode {
	case OP_JUMP:
		*next = int(param(0))
	case OP_JZER:
		if vm.pri == 0 {
			*next = int(param(0))
		}
	case OP_JNZ:
		if vm.pri != 0 {
			*next = int(param(0))
		}
	case OP_JEQ:
		if vm.pri == vm.alt {
			*next = int(param(0))
		}
	case OP_JNEQ:
		if vm.pri != vm.alt {
			*next = int(param(0))
		}
	case OP_JLESS:
		if uint32(vm.pri) < uint32(vm.alt) {
			*next = int(param(0))
		}
	case OP_JLEQ:
		if uint32(vm.pri) <= uint32(vm.alt) {
			*next = int(param(0))
		}
	case OP_JGRTR:
		if uint32(vm.pri) > uint32(vm.alt) {
			*next = int(param(0))
		}
	case OP_JGEQ:
		if uint32(vm.pri) >= uint32(vm.alt) {
			*next = int(param(0))
		}
	case OP_JSLESS:
		if vm.pri < vm.alt {
			*next = int(param(0))
		}
	case OP_JSLEQ:
		if vm.pri <= vm.alt {
			*next = int(param(0))
		}
	case OP_JSGRTR:
		if vm.pri > vm.alt {
			*next = int(param(0))
		}
	case OP_JSGEQ:
		if vm.pri >= vm.alt {
			*next = int(param(0))
		}
	case OP_JREL:
		*next += int(param(0))
	case OP_CALL:
		if err := vm.push(Cell(*next)); err != nil {
			return true, false, 0, err
		}
		*next = int(param(0))
	case OP_CALL_PRI:
		if err := vm.push(Cell(*next)); err != nil {
			return true, false, 0, err
		}
		*next = int(vm.pri)
	case OP_JUMP_PRI:
		*next = int(vm.pri)
	case OP_BOUNDS:
		if uint32(vm.pri) > uint32(param(0)) {
			return true, false, 0, RuntimeError{Code: RuntimeErrorBounds, Message: "bounds check failed", CIP: Cell(ins.Offset)}
		}
	case OP_SWITCH:
		target, err := vm.switchTarget(param(0))
		if err != nil {
			return true, false, 0, err
		}
		*next = target
	case OP_CASETBL:
		return true, false, 0, fmt.Errorf("%w: casetbl at executable offset %d", ErrInvalidInstruction, ins.Offset)
	case OP_SYSREQ_C:
		ret, err := vm.callNative(int(param(0)), -1)
		if err != nil {
			return true, false, 0, err
		}
		vm.pri = ret
	case OP_SYSREQ_PRI:
		ret, err := vm.callNative(int(vm.pri), -1)
		if err != nil {
			return true, false, 0, err
		}
		vm.pri = ret
	case OP_SYSREQ_N:
		ret, err := vm.callNative(int(param(0)), int(param(1)))
		if err != nil {
			return true, false, 0, err
		}
		vm.pri = ret
	case OP_FILE, OP_LINE, OP_SYMBOL, OP_SRANGE, OP_SYMTAG, OP_NOP, OP_BREAK:
		// Debug metadata is consumed by tooling and has no runtime effect.
	case OP_RET, OP_RETN:
		oldFrame, err := vm.pop()
		if err != nil {
			return true, true, vm.pri, nil
		}
		returnAddr, err := vm.pop()
		if err != nil {
			return true, true, vm.pri, nil
		}
		vm.frm = oldFrame
		if returnAddr == 0 {
			return true, true, vm.pri, nil
		}
		if ins.Opcode == OP_RETN {
			paramBytes, err := vm.ReadCell(vm.stk)
			if err != nil {
				return true, false, 0, err
			}
			vm.stk += paramBytes + cellBytes
			if err := vm.checkStack(); err != nil {
				return true, false, 0, err
			}
		}
		*next = int(returnAddr)
	case OP_HALT:
		if len(ins.Params) > 0 && ins.Params[0] != 0 {
			if ins.Params[0] == 12 {
				return true, false, 0, RuntimeError{Code: RuntimeErrorSleep, Message: "execution suspended", CIP: Cell(ins.Offset)}
			}
			return true, false, 0, RuntimeError{
				Code:    RuntimeErrorHalt,
				Message: fmt.Sprintf("halted with code %d", ins.Params[0]),
				CIP:     Cell(ins.Offset),
			}
		}
		return true, true, vm.pri, nil
	default:
		return false, false, 0, nil
	}
	return true, false, 0, nil
}
