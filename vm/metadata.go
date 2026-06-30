package vm

func (vm *VM) Publics() ([]Public, error) {
	out := make([]Public, 0, len(vm.publics))
	for i, pub := range vm.publics {
		out = append(out, Public{Index: i, Name: pub.Name})
	}
	return out, nil
}

func (vm *VM) Natives() ([]string, error) {
	out := make([]string, 0, len(vm.natives))
	for _, n := range vm.natives {
		out = append(out, n.Name)
	}
	return out, nil
}

func (vm *VM) PubVars() ([]Public, error) {
	out := make([]Public, 0, len(vm.pubvars))
	for i, pubvar := range vm.pubvars {
		out = append(out, Public{Index: i, Name: pubvar.Name})
	}
	return out, nil
}

func (vm *VM) Tags() ([]Public, error) {
	out := make([]Public, 0, len(vm.tags))
	for i, tag := range vm.tags {
		out = append(out, Public{Index: i, Name: tag.Name})
	}
	return out, nil
}

func (vm *VM) Libraries() ([]Public, error) {
	out := make([]Public, 0, len(vm.libraries))
	for i, library := range vm.libraries {
		out = append(out, Public{Index: i, Name: library.Name})
	}
	return out, nil
}

func (vm *VM) PubVarAddress(index int) Cell {
	if index < 0 || index >= len(vm.pubvars) {
		return 0
	}
	return Cell(vm.pubvars[index].Address)
}

func (vm *VM) TagID(index int) Cell {
	if index < 0 || index >= len(vm.tags) {
		return 0
	}
	return Cell(vm.tags[index].Address)
}
