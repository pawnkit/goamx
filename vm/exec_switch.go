package vm

import "fmt"

func (vm *VM) switchTarget(tableOffset Cell) (int, error) {
	code := vm.code()
	ins, err := decodeOne(code, int(tableOffset))
	if err != nil {
		return 0, err
	}
	if ins.Opcode != OP_CASETBL {
		return 0, fmt.Errorf("%w: switch target %d is not casetbl", ErrInvalidInstruction, tableOffset)
	}
	if len(ins.Params) < 2 {
		return 0, fmt.Errorf("%w: empty casetbl at offset %d", ErrInvalidInstruction, tableOffset)
	}
	target := int(ins.Params[1])
	for i := 2; i+1 < len(ins.Params); i += 2 {
		if ins.Params[i] == vm.pri {
			target = int(ins.Params[i+1])
			break
		}
	}
	return target, nil
}
