package amx

import (
	"errors"

	"github.com/pawnkit/goamx/vm"
)

func (r *Runtime) ExecPublic(index int, args ...Cell) (Cell, error) {
	ret, err := r.vm.ExecPublic(index, vmCells(args)...)
	if err != nil {
		return Cell(ret), translateError(err)
	}
	return Cell(ret), nil
}

func (r *Runtime) ExecMain(args ...Cell) (Cell, error) {
	ret, err := r.vm.ExecMain(vmCells(args)...)
	if err != nil {
		return Cell(ret), translateError(err)
	}
	return Cell(ret), nil
}

func (r *Runtime) Continue() (Cell, error) {
	ret, err := r.vm.Continue()
	if err != nil {
		return Cell(ret), translateError(err)
	}
	return Cell(ret), nil
}

func (r *Runtime) Suspended() bool { return r.vm.Suspended() }

func (r *Runtime) SetDebugHook(hook DebugHook) {
	if hook == nil {
		r.vm.SetDebugHook(nil)
		return
	}
	r.vm.SetDebugHook(func(ins vm.Instruction, state vm.State) error {
		return hook(DebugEvent{
			Instruction: publicInstruction(ins),
			State:       publicState(state),
		})
	})
}

func vmCells(values []Cell) []vm.Cell {
	out := make([]vm.Cell, len(values))
	for i, value := range values {
		out[i] = vm.Cell(value)
	}
	return out
}

func publicCells(values []vm.Cell) []Cell {
	out := make([]Cell, len(values))
	for i, value := range values {
		out[i] = Cell(value)
	}
	return out
}

func publicInstruction(ins vm.Instruction) Instruction {
	return Instruction{
		Offset: ins.Offset,
		Opcode: ins.Opcode,
		Params: publicCells(ins.Params),
		Size:   ins.Size,
	}
}

func publicState(state vm.State) State {
	return State{
		PRI: Cell(state.PRI),
		ALT: Cell(state.ALT),
		HEA: Cell(state.HEA),
		STK: Cell(state.STK),
		STP: Cell(state.STP),
		FRM: Cell(state.FRM),
		CIP: state.CIP,
	}
}

func (r *Runtime) SetUserData(tag int64, value any) {
	if r.userData == nil {
		r.userData = map[int64]any{}
	}
	r.userData[tag] = value
}

func (r *Runtime) UserData(tag int64) (any, bool) {
	value, ok := r.userData[tag]
	return value, ok
}

func (r *Runtime) Close() error {
	return r.vm.Close()
}

func translateError(err error) error {
	var runtimeErr vm.RuntimeError
	if errors.As(err, &runtimeErr) {
		return RuntimeError{
			Code:    RuntimeErrorCode(runtimeErr.Code),
			Message: runtimeErr.Message,
			CIP:     Cell(runtimeErr.CIP),
		}
	}
	return err
}
