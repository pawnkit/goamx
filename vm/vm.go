package vm

type VM struct {
	name        string
	image       []byte
	initial     []byte
	memory      []byte
	header      header
	sourceFlags uint16

	publics   []funcStub
	natives   []funcStub
	pubvars   []funcStub
	tags      []funcStub
	libraries []funcStub
	debug     DebugInfo

	registered map[string]NativeFunc
	debugHook  func(Instruction, State) error
	maxSteps   int

	pri Cell
	alt Cell
	hea Cell
	stk Cell
	stp Cell
	frm Cell

	suspended bool
	resumeCIP int
}

type State struct {
	PRI, ALT           Cell
	HEA, STK, STP, FRM Cell
	CIP                int
}

func (vm *VM) SetDebugHook(hook func(Instruction, State) error) { vm.debugHook = hook }
func (vm *VM) DebugInfo() DebugInfo                             { return vm.debug }
func (vm *VM) State() State {
	return State{PRI: vm.pri, ALT: vm.alt, HEA: vm.hea, STK: vm.stk, STP: vm.stp, FRM: vm.frm, CIP: vm.resumeCIP}
}

func (vm *VM) SetInstructionLimit(limit int) {
	if limit <= 0 {
		vm.maxSteps = maxExecSteps
		return
	}
	vm.maxSteps = limit
}

const stackMargin = 16 * cellBytes
