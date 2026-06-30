package vm

func (vm *VM) DebugLocation(cip Cell) (file string, line int, function string, ok bool) {
	address := uint32(cip)
	debugFile, fileOK := vm.debug.FileAt(address)
	debugLine, lineOK := vm.debug.LineAt(address)
	debugFunction, functionOK := vm.debug.FunctionAt(address)
	if fileOK {
		file = debugFile.Name
	}
	if lineOK {
		line = int(debugLine.Line)
	}
	if functionOK {
		function = debugFunction.Name
	}
	return file, line, function, fileOK || lineOK || functionOK
}

func (vm *VM) CoverageLocations() []CoverageLocation {
	seen := map[CoverageLocation]bool{}
	var locations []CoverageLocation
	for _, entry := range vm.debug.Lines {
		file, line, function, ok := vm.DebugLocation(Cell(entry.Address))
		location := CoverageLocation{File: file, Line: line, Function: function}
		if ok && file != "" && line > 0 && !seen[location] {
			seen[location] = true
			locations = append(locations, location)
		}
	}
	return locations
}

func (vm *VM) SetInstructionObserver(observer func(cip Cell)) {
	if observer == nil {
		vm.debugHook = nil
		return
	}
	vm.debugHook = func(instruction Instruction, _ State) error {
		observer(Cell(instruction.Offset))
		return nil
	}
}
