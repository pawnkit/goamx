package vm

func (vm *VM) execArithmeticInstruction(ins Instruction) (handled bool, err error) {
	param := func(i int) Cell { return ins.Params[i] }

	switch ins.Opcode {
	case OP_ADD:
		vm.pri += vm.alt
	case OP_ADD_C:
		vm.pri += param(0)
	case OP_SUB:
		vm.pri -= vm.alt
	case OP_SUB_ALT:
		vm.pri = vm.alt - vm.pri
	case OP_AND:
		vm.pri &= vm.alt
	case OP_OR:
		vm.pri |= vm.alt
	case OP_XOR:
		vm.pri ^= vm.alt
	case OP_NOT:
		vm.pri = boolCell(vm.pri == 0)
	case OP_NEG:
		vm.pri = -vm.pri
	case OP_INVERT:
		vm.pri = ^vm.pri
	case OP_SHL:
		vm.pri <<= vm.alt
	case OP_SHR:
		vm.pri = Cell(uint32(vm.pri) >> uint(vm.alt))
	case OP_SSHR:
		vm.pri >>= vm.alt
	case OP_SHL_C_PRI:
		vm.pri <<= param(0)
	case OP_SHL_C_ALT:
		vm.alt <<= param(0)
	case OP_SHR_C_PRI:
		vm.pri = Cell(uint32(vm.pri) >> uint(param(0)))
	case OP_SHR_C_ALT:
		vm.alt = Cell(uint32(vm.alt) >> uint(param(0)))
	case OP_SMUL:
		vm.pri *= vm.alt
	case OP_SMUL_C:
		vm.pri *= param(0)
	case OP_UMUL:
		vm.pri = Cell(uint32(vm.pri) * uint32(vm.alt))
	case OP_SDIV:
		if vm.alt == 0 {
			return true, divideByZero(ins)
		}
		dividend := vm.pri
		divisor := vm.alt
		vm.pri = dividend / divisor
		vm.alt = dividend % divisor
	case OP_UDIV:
		if vm.alt == 0 {
			return true, divideByZero(ins)
		}
		dividend := uint32(vm.pri)
		divisor := uint32(vm.alt)
		vm.pri = Cell(dividend / divisor)
		vm.alt = Cell(dividend % divisor)
	case OP_SDIV_ALT:
		if vm.pri == 0 {
			return true, divideByZero(ins)
		}
		dividend := vm.alt
		divisor := vm.pri
		vm.pri = dividend / divisor
		vm.alt = dividend % divisor
	case OP_UDIV_ALT:
		if vm.pri == 0 {
			return true, divideByZero(ins)
		}
		dividend := uint32(vm.alt)
		divisor := uint32(vm.pri)
		vm.pri = Cell(dividend / divisor)
		vm.alt = Cell(dividend % divisor)
	case OP_SIGN_PRI:
		if vm.pri&0xff >= 0x80 {
			vm.pri |= Cell(-256)
		}
	case OP_SIGN_ALT:
		if vm.alt&0xff >= 0x80 {
			vm.alt |= Cell(-256)
		}
	case OP_EQ:
		vm.pri = boolCell(vm.pri == vm.alt)
	case OP_NEQ:
		vm.pri = boolCell(vm.pri != vm.alt)
	case OP_LESS:
		vm.pri = boolCell(uint32(vm.pri) < uint32(vm.alt))
	case OP_LEQ:
		vm.pri = boolCell(uint32(vm.pri) <= uint32(vm.alt))
	case OP_GRTR:
		vm.pri = boolCell(uint32(vm.pri) > uint32(vm.alt))
	case OP_GEQ:
		vm.pri = boolCell(uint32(vm.pri) >= uint32(vm.alt))
	case OP_SLESS:
		vm.pri = boolCell(vm.pri < vm.alt)
	case OP_SLEQ:
		vm.pri = boolCell(vm.pri <= vm.alt)
	case OP_SGRTR:
		vm.pri = boolCell(vm.pri > vm.alt)
	case OP_SGEQ:
		vm.pri = boolCell(vm.pri >= vm.alt)
	case OP_EQ_C_PRI:
		vm.pri = boolCell(vm.pri == param(0))
	case OP_EQ_C_ALT:
		vm.pri = boolCell(vm.alt == param(0))
	default:
		return false, nil
	}
	return true, nil
}

func divideByZero(ins Instruction) RuntimeError {
	return RuntimeError{Code: RuntimeErrorDivideByZero, Message: "divide by zero", CIP: Cell(ins.Offset)}
}
