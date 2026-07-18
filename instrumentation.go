package amx

import "github.com/pawnkit/goamx/vm"

type EventKind string

const (
	EventInstruction EventKind = "instruction"
	EventPublicEnter EventKind = "public-enter"
	EventPublicExit  EventKind = "public-exit"
	EventNativeEnter EventKind = "native-enter"
	EventNativeExit  EventKind = "native-exit"
	EventException   EventKind = "exception"
)

type InstrumentationEvent struct {
	Kind        EventKind
	Name        string
	Instruction Instruction
	State       State
	File        string
	Line        int
	Result      Cell
	Err         error
}

type InstrumentationHook func(InstrumentationEvent) error

func (r *Runtime) SetInstrumentationHook(hook InstrumentationHook) {
	r.hook = hook
	if hook == nil {
		r.vm.SetDebugHook(nil)
		return
	}
	r.vm.SetDebugHook(func(ins vm.Instruction, state vm.State) error {
		file, line, _, _ := r.vm.DebugLocation(vm.Cell(ins.Offset))
		return hook(InstrumentationEvent{Kind: EventInstruction, Instruction: publicInstruction(ins), State: publicState(state), File: file, Line: line})
	})
}

func (r *Runtime) emit(event InstrumentationEvent) error {
	if r.hook == nil {
		return nil
	}
	return r.hook(event)
}
