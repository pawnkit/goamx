package amx

import (
	"fmt"

	"github.com/pawnkit/goamx/vm"
)

type Runtime struct {
	vm       *vm.VM
	userData map[int64]any
	hook     InstrumentationHook
}

// VM returns the underlying low-level virtual machine.
//
// Most applications should prefer Runtime's higher-level methods. The VM is
// exposed for tooling that needs direct access to the bytecode-level API.
func (r *Runtime) VM() *vm.VM { return r.vm }

func (r *Runtime) Info() Info {
	info := r.vm.Info()
	return Info{
		Name:           info.Name,
		FileVersion:    info.FileVersion,
		AMXVersion:     info.AMXVersion,
		Flags:          info.Flags,
		DefinitionSize: info.DefinitionSize,
		CodeSize:       info.CodeSize,
		DataSize:       info.DataSize,
		StackHeapSize:  info.StackHeapSize,
		NameLength:     info.NameLength,
		HasMain:        info.Main,
	}
}

func (r *Runtime) MemoryInfo() MemoryInfo {
	info := r.Info()
	return MemoryInfo{CodeBytes: info.CodeSize, DataBytes: info.DataSize, StackHeapBytes: info.StackHeapSize}
}

// State returns the current registers, including the stopped instruction.
func (r *Runtime) State() State { return publicState(r.vm.State()) }

// SetInstructionLimit sets the maximum instructions for one execution call.
// Values less than one restore the default limit.
func (r *Runtime) SetInstructionLimit(limit int) { r.vm.SetInstructionLimit(limit) }

func (r *Runtime) Publics() ([]Public, error) {
	publics, err := r.vm.Publics()
	if err != nil {
		return nil, err
	}
	out := make([]Public, 0, len(publics))
	for _, p := range publics {
		out = append(out, Public{Index: p.Index, Name: p.Name})
	}
	return out, nil
}

func (r *Runtime) FindPublic(name string) (Public, bool) {
	publics, err := r.Publics()
	if err != nil {
		return Public{}, false
	}
	for _, public := range publics {
		if public.Name == name {
			return public, true
		}
	}
	return Public{}, false
}

func (r *Runtime) ExecPublicByName(name string, args ...Cell) (Cell, error) {
	public, ok := r.FindPublic(name)
	if !ok {
		return 0, fmt.Errorf("%w: %s", ErrPublicIndexOutOfRange, name)
	}
	return r.ExecPublic(public.Index, args...)
}

func (r *Runtime) Natives() ([]Native, error) {
	names, err := r.vm.Natives()
	if err != nil {
		return nil, err
	}
	vm := r.vm
	out := make([]Native, 0, len(names))
	for i, name := range names {
		out = append(out, Native{Index: i, Name: name, Registered: vm.NativeRegistered(name)})
	}
	return out, nil
}

func (r *Runtime) FindNative(name string) (Native, bool) {
	natives, err := r.Natives()
	if err != nil {
		return Native{}, false
	}
	for _, native := range natives {
		if native.Name == name {
			return native, true
		}
	}
	return Native{}, false
}

func (r *Runtime) Libraries() ([]Library, error) {
	entries, err := r.vm.Libraries()
	if err != nil {
		return nil, err
	}
	out := make([]Library, 0, len(entries))
	for _, entry := range entries {
		out = append(out, Library{Index: entry.Index, Name: entry.Name})
	}
	return out, nil
}

func (r *Runtime) Decode() ([]Instruction, error) {
	instructions, err := r.vm.Decode()
	if err != nil {
		return nil, err
	}
	out := make([]Instruction, 0, len(instructions))
	for _, ins := range instructions {
		out = append(out, publicInstruction(ins))
	}
	return out, nil
}

func (r *Runtime) PubVars() ([]PublicVar, error) {
	pubvars, err := r.vm.PubVars()
	if err != nil {
		return nil, err
	}
	out := make([]PublicVar, 0, len(pubvars))
	for _, pubvar := range pubvars {
		out = append(out, PublicVar{
			Index:   pubvar.Index,
			Name:    pubvar.Name,
			Address: Cell(r.vm.PubVarAddress(pubvar.Index)),
		})
	}
	return out, nil
}

func (r *Runtime) FindPubVar(name string) (PublicVar, bool) {
	entries, err := r.PubVars()
	if err != nil {
		return PublicVar{}, false
	}
	for _, entry := range entries {
		if entry.Name == name {
			return entry, true
		}
	}
	return PublicVar{}, false
}

func (r *Runtime) Tags() ([]Tag, error) {
	tags, err := r.vm.Tags()
	if err != nil {
		return nil, err
	}
	out := make([]Tag, 0, len(tags))
	for _, tag := range tags {
		out = append(out, Tag{Index: tag.Index, Name: tag.Name, ID: Cell(r.vm.TagID(tag.Index))})
	}
	return out, nil
}

func (r *Runtime) FindTag(id Cell) (Tag, bool) {
	entries, err := r.Tags()
	if err != nil {
		return Tag{}, false
	}
	for _, entry := range entries {
		if entry.ID == id {
			return entry, true
		}
	}
	return Tag{}, false
}

func (r *Runtime) FindTagByName(name string) (Tag, bool) {
	entries, err := r.Tags()
	if err != nil {
		return Tag{}, false
	}
	for _, entry := range entries {
		if entry.Name == name {
			return entry, true
		}
	}
	return Tag{}, false
}
