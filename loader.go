package amx

import (
	"github.com/pawnkit/goamx/vm"
)

func LoadFile(path string) (*Runtime, error) {
	machine, err := vm.LoadFile(path)
	if err != nil {
		return nil, err
	}
	return &Runtime{vm: machine, userData: map[int64]any{}}, nil
}

func LoadBytes(name string, data []byte) (*Runtime, error) {
	machine, err := vm.LoadBytes(name, data)
	if err != nil {
		return nil, err
	}
	return &Runtime{vm: machine, userData: map[int64]any{}}, nil
}
