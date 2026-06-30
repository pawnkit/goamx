package amx

import (
	"errors"

	"github.com/pawnkit/goamx/vm"
)

var ErrNativeNotDeclared = errors.New("native is not declared by the AMX image")

func (r *Runtime) RegisterNative(name string, fn NativeFunc) error {
	if fn == nil {
		return errors.New("native callback is nil")
	}
	if _, ok := r.FindNative(name); !ok {
		return ErrNativeNotDeclared
	}
	return r.vm.RegisterNative(name, func(ctx vm.NativeContext, params []vm.Cell) (vm.Cell, error) {
		machine, ok := ctx.(*vm.VM)
		if !ok {
			return 0, errors.New("native context is not a Go AMX runtime")
		}
		value, err := fn(runtimeContext{vm: machine}, publicCells(params))
		return vm.Cell(value), err
	})
}

func (r *Runtime) RegisterNatives(natives map[string]NativeFunc) error {
	for name, fn := range natives {
		if err := r.RegisterNative(name, fn); err != nil {
			return err
		}
	}
	return nil
}

type runtimeContext struct{ vm *vm.VM }

func (ctx runtimeContext) ReadString(addr Cell) (string, error) {
	return ctx.vm.ReadString(vm.Cell(addr))
}

func (ctx runtimeContext) WriteString(addr Cell, value string) error {
	return ctx.vm.WriteString(vm.Cell(addr), value)
}

func (ctx runtimeContext) ReadCell(addr Cell) (Cell, error) {
	value, err := ctx.vm.ReadCell(vm.Cell(addr))
	return Cell(value), err
}

func (ctx runtimeContext) WriteCell(addr Cell, value Cell) error {
	return ctx.vm.WriteCell(vm.Cell(addr), vm.Cell(value))
}

func (ctx runtimeContext) ReadBytes(addr Cell, size int) ([]byte, error) {
	return ctx.vm.ReadBytes(vm.Cell(addr), size)
}

func (ctx runtimeContext) WriteBytes(addr Cell, value []byte) error {
	return ctx.vm.WriteBytes(vm.Cell(addr), value)
}

func (ctx runtimeContext) CallPublic(name string, args ...Cell) (Cell, error) {
	value, err := ctx.vm.CallPublic(name, vmCells(args)...)
	return Cell(value), translateError(err)
}
