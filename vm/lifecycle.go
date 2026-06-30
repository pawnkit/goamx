package vm

import "fmt"

func (vm *VM) Reset() error {
	vm.pri = 0
	vm.alt = 0
	vm.resetExecutionFrame()
	vm.suspended = false
	vm.resumeCIP = 0
	return nil
}

func (vm *VM) resetExecutionFrame() {
	vm.hea = Cell(vm.header.HEA - vm.header.DAT)
	vm.stk = Cell(vm.header.STP - vm.header.DAT)
	vm.stp = Cell(vm.header.STP - vm.header.DAT)
	vm.frm = 0
}

func (vm *VM) ResetMemory() error {
	copy(vm.memory, vm.initial)
	clear(vm.memory[len(vm.initial):])
	return vm.Reset()
}

func (vm *VM) SnapshotMemory() []byte {
	return append([]byte(nil), vm.memory...)
}

func (vm *VM) RestoreMemory(snapshot []byte) error {
	if len(snapshot) != len(vm.memory) {
		return fmt.Errorf("%w: snapshot size %d, want %d", ErrInvalidMemoryAccess, len(snapshot), len(vm.memory))
	}
	copy(vm.memory, snapshot)
	return vm.Reset()
}

func (vm *VM) Clone() *VM {
	clone := *vm
	clone.image = append([]byte(nil), vm.image...)
	clone.initial = append([]byte(nil), vm.initial...)
	clone.memory = append([]byte(nil), vm.memory...)
	clone.publics = append([]funcStub(nil), vm.publics...)
	clone.natives = append([]funcStub(nil), vm.natives...)
	clone.libraries = append([]funcStub(nil), vm.libraries...)
	clone.pubvars = append([]funcStub(nil), vm.pubvars...)
	clone.tags = append([]funcStub(nil), vm.tags...)
	clone.debug.Files = append([]DebugFile(nil), vm.debug.Files...)
	clone.debug.Lines = append([]DebugLine(nil), vm.debug.Lines...)
	clone.debug.Symbols = append([]DebugSymbol(nil), vm.debug.Symbols...)
	for i := range clone.debug.Symbols {
		clone.debug.Symbols[i].Dimensions = append([]DebugDimension(nil), vm.debug.Symbols[i].Dimensions...)
	}
	clone.debug.Tags = append([]DebugTag(nil), vm.debug.Tags...)
	clone.debug.Automata = append([]DebugAutomaton(nil), vm.debug.Automata...)
	clone.debug.States = append([]DebugState(nil), vm.debug.States...)
	clone.registered = make(map[string]NativeFunc, len(vm.registered))
	for name, fn := range vm.registered {
		clone.registered[name] = fn
	}
	return &clone
}

func (vm *VM) Close() error {
	return nil
}
