package amx

import (
	"github.com/pawnkit/goamx/vm"
)

func (r *Runtime) ReadCell(addr Cell) (Cell, error) {
	value, err := r.vm.ReadCell(vm.Cell(addr))
	return Cell(value), err
}

func (r *Runtime) WriteCell(addr Cell, value Cell) error {
	return r.vm.WriteCell(vm.Cell(addr), vm.Cell(value))
}

func (r *Runtime) ReadString(addr Cell) (string, error) {
	return r.vm.ReadString(vm.Cell(addr))
}

func (r *Runtime) WriteString(addr Cell, value string) error {
	return r.vm.WriteString(vm.Cell(addr), value)
}

func (r *Runtime) WriteStringN(addr Cell, value string, maxCells int, packed bool) error {
	return r.vm.WriteStringN(vm.Cell(addr), value, maxCells, packed)
}

func (r *Runtime) ReadBytes(addr Cell, size int) ([]byte, error) {
	return r.vm.ReadBytes(vm.Cell(addr), size)
}

func (r *Runtime) WriteBytes(addr Cell, value []byte) error {
	return r.vm.WriteBytes(vm.Cell(addr), value)
}

func (r *Runtime) Allot(cells int) (Cell, error) {
	value, err := r.vm.Allot(cells)
	return Cell(value), err
}

func (r *Runtime) Release(addr Cell) error { return r.vm.Release(vm.Cell(addr)) }
func (r *Runtime) Reset() error            { return r.vm.Reset() }
func (r *Runtime) ResetMemory() error      { return r.vm.ResetMemory() }

func (r *Runtime) AllotCells(values []Cell) (Cell, error) {
	addr, err := r.Allot(len(values))
	if err != nil {
		return 0, err
	}
	for i, value := range values {
		if err := r.WriteCell(addr+Cell(i*CellBytes), value); err != nil {
			_ = r.Release(addr)
			return 0, err
		}
	}
	return addr, nil
}

func (r *Runtime) AllotString(value string, packed bool) (Cell, error) {
	cells := len([]rune(value)) + 1
	if packed {
		cells = (len([]byte(value)) + 4) / 4
	}
	addr, err := r.Allot(cells)
	if err != nil {
		return 0, err
	}
	if err := r.WriteStringN(addr, value, cells, packed); err != nil {
		_ = r.Release(addr)
		return 0, err
	}
	return addr, nil
}

func (r *Runtime) Clone() *Runtime {
	clone := &Runtime{vm: r.vm.Clone(), userData: map[int64]any{}}
	for tag, value := range r.userData {
		clone.userData[tag] = value
	}
	return clone
}
