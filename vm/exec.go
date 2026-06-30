package vm

import (
	"errors"
	"fmt"
)

const maxExecSteps = 1_000_000

func (vm *VM) execAt(address uint32, args ...Cell) (Cell, error) {
	if err := vm.Reset(); err != nil {
		return 0, err
	}
	for i := len(args) - 1; i >= 0; i-- {
		if err := vm.push(args[i]); err != nil {
			return 0, err
		}
	}
	if err := vm.push(Cell(len(args) * cellBytes)); err != nil {
		return 0, err
	}
	if err := vm.push(0); err != nil {
		return 0, err
	}
	code := vm.code()
	if int(address) >= len(code) {
		return 0, fmt.Errorf("%w: entry address %d outside code", ErrInvalidInstruction, address)
	}
	return vm.execLoop(int(address))
}

func (vm *VM) execLoop(cip int) (Cell, error) {
	code := vm.code()
	for steps := 0; steps < vm.maxSteps; steps++ {
		ins, err := decodeOne(code, cip)
		if err != nil {
			return 0, err
		}
		if vm.debugHook != nil {
			state := State{PRI: vm.pri, ALT: vm.alt, HEA: vm.hea, STK: vm.stk, STP: vm.stp, FRM: vm.frm, CIP: cip}
			if err := vm.debugHook(ins, state); err != nil {
				return 0, err
			}
		}
		next := cip + ins.Size
		halt, ret, err := vm.execInstruction(ins, &next)
		if err != nil {
			var runtimeErr RuntimeError
			if errors.As(err, &runtimeErr) && runtimeErr.Code == RuntimeErrorSleep {
				vm.suspended = true
				vm.resumeCIP = next
				return vm.pri, err
			}
			vm.resetExecutionFrame()
			return 0, err
		}
		if halt {
			vm.resetExecutionFrame()
			return ret, nil
		}
		if next < 0 || next > len(code) {
			return 0, fmt.Errorf("%w: instruction pointer %d outside code", ErrInvalidInstruction, next)
		}
		cip = next
	}
	vm.resetExecutionFrame()
	return 0, fmt.Errorf("%w: execution step limit exceeded", ErrUnsupportedExecution)
}

func (vm *VM) execInstruction(ins Instruction, next *int) (bool, Cell, error) {
	param := func(i int) Cell { return ins.Params[i] }
	if handled, halt, ret, err := vm.execControlInstruction(ins, next); handled {
		return halt, ret, err
	}
	if handled, err := vm.execArithmeticInstruction(ins); handled {
		return false, 0, err
	}

	switch ins.Opcode {
	case OP_LOAD_PRI:
		value, err := vm.ReadCell(param(0))
		if err != nil {
			return false, 0, err
		}
		vm.pri = value
	case OP_LOAD_ALT:
		value, err := vm.ReadCell(param(0))
		if err != nil {
			return false, 0, err
		}
		vm.alt = value
	case OP_LOAD_S_PRI:
		value, err := vm.ReadCell(vm.frm + param(0))
		if err != nil {
			return false, 0, err
		}
		vm.pri = value
	case OP_LOAD_S_ALT:
		value, err := vm.ReadCell(vm.frm + param(0))
		if err != nil {
			return false, 0, err
		}
		vm.alt = value
	case OP_LOAD_BOTH:
		pri, err := vm.ReadCell(param(0))
		if err != nil {
			return false, 0, err
		}
		alt, err := vm.ReadCell(param(1))
		if err != nil {
			return false, 0, err
		}
		vm.pri = pri
		vm.alt = alt
	case OP_LOAD_S_BOTH:
		pri, err := vm.ReadCell(vm.frm + param(0))
		if err != nil {
			return false, 0, err
		}
		alt, err := vm.ReadCell(vm.frm + param(1))
		if err != nil {
			return false, 0, err
		}
		vm.pri = pri
		vm.alt = alt
	case OP_LOAD_I:
		value, err := vm.ReadCell(vm.pri)
		if err != nil {
			return false, 0, err
		}
		vm.pri = value
	case OP_LODB_I:
		value, err := readSized(vm.memory, vm.pri, param(0))
		if err != nil {
			return false, 0, err
		}
		vm.pri = value
	case OP_LREF_PRI:
		addr, err := vm.ReadCell(param(0))
		if err != nil {
			return false, 0, err
		}
		value, err := vm.ReadCell(addr)
		if err != nil {
			return false, 0, err
		}
		vm.pri = value
	case OP_LREF_ALT:
		addr, err := vm.ReadCell(param(0))
		if err != nil {
			return false, 0, err
		}
		value, err := vm.ReadCell(addr)
		if err != nil {
			return false, 0, err
		}
		vm.alt = value
	case OP_LREF_S_PRI:
		addr, err := vm.ReadCell(vm.frm + param(0))
		if err != nil {
			return false, 0, err
		}
		value, err := vm.ReadCell(addr)
		if err != nil {
			return false, 0, err
		}
		vm.pri = value
	case OP_LREF_S_ALT:
		addr, err := vm.ReadCell(vm.frm + param(0))
		if err != nil {
			return false, 0, err
		}
		value, err := vm.ReadCell(addr)
		if err != nil {
			return false, 0, err
		}
		vm.alt = value
	case OP_LIDX:
		value, err := vm.ReadCell(vm.alt + vm.pri*cellBytes)
		if err != nil {
			return false, 0, err
		}
		vm.pri = value
	case OP_LIDX_B:
		value, err := vm.ReadCell(vm.alt + (vm.pri << param(0)))
		if err != nil {
			return false, 0, err
		}
		vm.pri = value
	case OP_IDXADDR:
		vm.pri = vm.alt + vm.pri*cellBytes
	case OP_IDXADDR_B:
		vm.pri = vm.alt + (vm.pri << param(0))
	case OP_ALIGN_PRI, OP_ALIGN_ALT:
		// AMX addresses only need adjustment when the host and image byte order
		// differ. The loader reads the little-endian AMX image into native cell
		// values, so alignment is intentionally a no-op here.
	case OP_CONST_PRI:
		vm.pri = param(0)
	case OP_CONST_ALT:
		vm.alt = param(0)
	case OP_ADDR_PRI:
		vm.pri = vm.frm + param(0)
	case OP_ADDR_ALT:
		vm.alt = vm.frm + param(0)
	case OP_CONST:
		if err := vm.WriteCell(param(0), param(1)); err != nil {
			return false, 0, err
		}
	case OP_CONST_S:
		if err := vm.WriteCell(vm.frm+param(0), param(1)); err != nil {
			return false, 0, err
		}
	case OP_STOR_PRI:
		if err := vm.WriteCell(param(0), vm.pri); err != nil {
			return false, 0, err
		}
	case OP_STOR_ALT:
		if err := vm.WriteCell(param(0), vm.alt); err != nil {
			return false, 0, err
		}
	case OP_STOR_S_PRI:
		if err := vm.WriteCell(vm.frm+param(0), vm.pri); err != nil {
			return false, 0, err
		}
	case OP_STOR_S_ALT:
		if err := vm.WriteCell(vm.frm+param(0), vm.alt); err != nil {
			return false, 0, err
		}
	case OP_STOR_I:
		if err := vm.WriteCell(vm.alt, vm.pri); err != nil {
			return false, 0, err
		}
	case OP_STRB_I:
		if err := writeSized(vm.memory, vm.alt, param(0), vm.pri); err != nil {
			return false, 0, err
		}
	case OP_SREF_PRI:
		addr, err := vm.ReadCell(param(0))
		if err != nil {
			return false, 0, err
		}
		if err := vm.WriteCell(addr, vm.pri); err != nil {
			return false, 0, err
		}
	case OP_SREF_ALT:
		addr, err := vm.ReadCell(param(0))
		if err != nil {
			return false, 0, err
		}
		if err := vm.WriteCell(addr, vm.alt); err != nil {
			return false, 0, err
		}
	case OP_SREF_S_PRI:
		addr, err := vm.ReadCell(vm.frm + param(0))
		if err != nil {
			return false, 0, err
		}
		if err := vm.WriteCell(addr, vm.pri); err != nil {
			return false, 0, err
		}
	case OP_SREF_S_ALT:
		addr, err := vm.ReadCell(vm.frm + param(0))
		if err != nil {
			return false, 0, err
		}
		if err := vm.WriteCell(addr, vm.alt); err != nil {
			return false, 0, err
		}
	case OP_ZERO_PRI:
		vm.pri = 0
	case OP_ZERO_ALT:
		vm.alt = 0
	case OP_ZERO:
		if err := vm.WriteCell(param(0), 0); err != nil {
			return false, 0, err
		}
	case OP_ZERO_S:
		if err := vm.WriteCell(vm.frm+param(0), 0); err != nil {
			return false, 0, err
		}
	case OP_INC_PRI:
		vm.pri++
	case OP_INC_ALT:
		vm.alt++
	case OP_INC:
		value, err := vm.ReadCell(param(0))
		if err != nil {
			return false, 0, err
		}
		if err := vm.WriteCell(param(0), value+1); err != nil {
			return false, 0, err
		}
	case OP_INC_S:
		addr := vm.frm + param(0)
		value, err := vm.ReadCell(addr)
		if err != nil {
			return false, 0, err
		}
		if err := vm.WriteCell(addr, value+1); err != nil {
			return false, 0, err
		}
	case OP_DEC_PRI:
		vm.pri--
	case OP_DEC_ALT:
		vm.alt--
	case OP_DEC:
		value, err := vm.ReadCell(param(0))
		if err != nil {
			return false, 0, err
		}
		if err := vm.WriteCell(param(0), value-1); err != nil {
			return false, 0, err
		}
	case OP_DEC_S:
		addr := vm.frm + param(0)
		value, err := vm.ReadCell(addr)
		if err != nil {
			return false, 0, err
		}
		if err := vm.WriteCell(addr, value-1); err != nil {
			return false, 0, err
		}
	case OP_INC_I:
		value, err := vm.ReadCell(vm.pri)
		if err != nil {
			return false, 0, err
		}
		if err := vm.WriteCell(vm.pri, value+1); err != nil {
			return false, 0, err
		}
	case OP_DEC_I:
		value, err := vm.ReadCell(vm.pri)
		if err != nil {
			return false, 0, err
		}
		if err := vm.WriteCell(vm.pri, value-1); err != nil {
			return false, 0, err
		}
	case OP_MOVS:
		srcStart, srcEnd, err := checkedByteRange(vm.memory, vm.pri, param(0))
		if err != nil {
			return false, 0, err
		}
		dstStart, _, err := checkedByteRange(vm.memory, vm.alt, param(0))
		if err != nil {
			return false, 0, err
		}
		copy(vm.memory[dstStart:dstStart+(srcEnd-srcStart)], vm.memory[srcStart:srcEnd])
	case OP_CMPS:
		leftStart, leftEnd, err := checkedByteRange(vm.memory, vm.alt, param(0))
		if err != nil {
			return false, 0, err
		}
		rightStart, rightEnd, err := checkedByteRange(vm.memory, vm.pri, param(0))
		if err != nil {
			return false, 0, err
		}
		vm.pri = compareBytes(vm.memory[leftStart:leftEnd], vm.memory[rightStart:rightEnd])
	case OP_FILL:
		if param(0)%cellBytes != 0 {
			return false, 0, fmt.Errorf("%w: fill byte count %d is not cell-aligned", ErrInvalidInstruction, param(0))
		}
		if _, _, err := checkedByteRange(vm.memory, vm.alt, param(0)); err != nil {
			return false, 0, err
		}
		for offset := Cell(0); offset < param(0); offset += cellBytes {
			if err := vm.WriteCell(vm.alt+offset, vm.pri); err != nil {
				return false, 0, err
			}
		}
	case OP_PUSH_PRI:
		if err := vm.push(vm.pri); err != nil {
			return false, 0, err
		}
	case OP_PUSH_ALT:
		if err := vm.push(vm.alt); err != nil {
			return false, 0, err
		}
	case OP_PUSH_C:
		if err := vm.push(param(0)); err != nil {
			return false, 0, err
		}
	case OP_PUSH_R:
		for i := Cell(0); i < param(0); i++ {
			if err := vm.push(vm.pri); err != nil {
				return false, 0, err
			}
		}
	case OP_PUSH:
		value, err := vm.ReadCell(param(0))
		if err != nil {
			return false, 0, err
		}
		if err := vm.push(value); err != nil {
			return false, 0, err
		}
	case OP_PUSH_S:
		value, err := vm.ReadCell(vm.frm + param(0))
		if err != nil {
			return false, 0, err
		}
		if err := vm.push(value); err != nil {
			return false, 0, err
		}
	case OP_PUSH_ADR:
		if err := vm.push(vm.frm + param(0)); err != nil {
			return false, 0, err
		}
	case OP_PUSH2_C, OP_PUSH3_C, OP_PUSH4_C, OP_PUSH5_C:
		if err := vm.pushParams(ins.Params); err != nil {
			return false, 0, err
		}
	case OP_PUSH2, OP_PUSH3, OP_PUSH4, OP_PUSH5:
		for _, addr := range ins.Params {
			value, err := vm.ReadCell(addr)
			if err != nil {
				return false, 0, err
			}
			if err := vm.push(value); err != nil {
				return false, 0, err
			}
		}
	case OP_PUSH2_S, OP_PUSH3_S, OP_PUSH4_S, OP_PUSH5_S:
		for _, addr := range ins.Params {
			value, err := vm.ReadCell(vm.frm + addr)
			if err != nil {
				return false, 0, err
			}
			if err := vm.push(value); err != nil {
				return false, 0, err
			}
		}
	case OP_PUSH2_ADR, OP_PUSH3_ADR, OP_PUSH4_ADR, OP_PUSH5_ADR:
		for _, addr := range ins.Params {
			if err := vm.push(vm.frm + addr); err != nil {
				return false, 0, err
			}
		}
	case OP_POP_PRI:
		value, err := vm.pop()
		if err != nil {
			return false, 0, err
		}
		vm.pri = value
	case OP_POP_ALT:
		value, err := vm.pop()
		if err != nil {
			return false, 0, err
		}
		vm.alt = value
	case OP_STACK:
		vm.alt = vm.stk
		vm.stk += param(0)
		if err := vm.checkStack(); err != nil {
			return false, 0, err
		}
	case OP_HEAP:
		vm.alt = vm.hea
		vm.hea += param(0)
		if err := vm.checkHeap(); err != nil {
			return false, 0, err
		}
	case OP_PROC:
		if err := vm.push(vm.frm); err != nil {
			return false, 0, err
		}
		vm.frm = vm.stk
	case OP_LCTRL:
		switch param(0) {
		case 0:
			vm.pri = Cell(vm.header.COD)
		case 1:
			vm.pri = Cell(vm.header.DAT)
		case 2:
			vm.pri = vm.hea
		case 3:
			vm.pri = vm.stp
		case 4:
			vm.pri = vm.stk
		case 5:
			vm.pri = vm.frm
		case 6:
			vm.pri = Cell(*next)
		}
	case OP_SCTRL:
		switch param(0) {
		case 2:
			vm.hea = vm.pri
		case 4:
			vm.stk = vm.pri
			if err := vm.checkStack(); err != nil {
				return false, 0, err
			}
		case 5:
			vm.frm = vm.pri
		case 6:
			*next = int(vm.pri)
		}
	case OP_MOVE_PRI:
		vm.pri = vm.alt
	case OP_MOVE_ALT:
		vm.alt = vm.pri
	case OP_XCHG:
		vm.pri, vm.alt = vm.alt, vm.pri
	case OP_SWAP_PRI:
		value, err := vm.ReadCell(vm.stk)
		if err != nil {
			return false, 0, err
		}
		if err := vm.WriteCell(vm.stk, vm.pri); err != nil {
			return false, 0, err
		}
		vm.pri = value
	case OP_SWAP_ALT:
		value, err := vm.ReadCell(vm.stk)
		if err != nil {
			return false, 0, err
		}
		if err := vm.WriteCell(vm.stk, vm.alt); err != nil {
			return false, 0, err
		}
		vm.alt = value
	default:
		info, _ := ins.Opcode.Info()
		return false, 0, fmt.Errorf("%w: %s at offset %d", ErrUnsupportedExecution, info.Name, ins.Offset)
	}
	return false, 0, nil
}
