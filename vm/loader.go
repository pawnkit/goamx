package vm

import (
	"os"
)

func LoadFile(path string) (*VM, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return LoadBytes(path, data)
}

func LoadBytes(name string, data []byte) (*VM, error) {
	h, err := parseHeader(data)
	if err != nil {
		return nil, err
	}
	sourceFlags := h.Flags
	if h.Flags&amxFlagCompact != 0 {
		data, h, err = expandCompactImage(data, h)
		if err != nil {
			return nil, err
		}
	}
	memory, err := initialMemory(data, h)
	if err != nil {
		return nil, err
	}
	publics, err := parseFuncTable(data, h.Publics, h.Natives, h.DefSize)
	if err != nil {
		return nil, err
	}
	natives, err := parseFuncTable(data, h.Natives, h.Libraries, h.DefSize)
	if err != nil {
		return nil, err
	}
	libraries, err := parseFuncTable(data, h.Libraries, h.PubVars, h.DefSize)
	if err != nil {
		return nil, err
	}
	pubvars, err := parseFuncTable(data, h.PubVars, h.Tags, h.DefSize)
	if err != nil {
		return nil, err
	}
	tagEnd := h.NameTable
	if h.DefSize != 8 {
		tagEnd = h.COD
	}
	tags, err := parseFuncTable(data, h.Tags, tagEnd, h.DefSize)
	if err != nil {
		return nil, err
	}
	debug, err := parseDebugInfo(data, h)
	if err != nil {
		return nil, err
	}
	vm := &VM{
		name:        name,
		image:       append([]byte(nil), data...),
		initial:     append([]byte(nil), memory...),
		memory:      memory,
		header:      h,
		sourceFlags: sourceFlags,
		publics:     publics,
		natives:     natives,
		libraries:   libraries,
		pubvars:     pubvars,
		tags:        tags,
		debug:       debug,
		maxSteps:    maxExecSteps,
		registered:  map[string]NativeFunc{},
	}
	if err := vm.Reset(); err != nil {
		return nil, err
	}
	return vm, nil
}
