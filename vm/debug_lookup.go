package vm

func (info DebugInfo) FileAt(address uint32) (DebugFile, bool) {
	var found DebugFile
	var ok bool
	for _, entry := range info.Files {
		if entry.Address > address {
			break
		}
		found, ok = entry, true
	}
	return found, ok
}

func (info DebugInfo) LineAt(address uint32) (DebugLine, bool) {
	var found DebugLine
	var ok bool
	for _, entry := range info.Lines {
		if entry.Address > address {
			break
		}
		found, ok = entry, true
	}
	return found, ok
}

func (info DebugInfo) FunctionAt(address uint32) (DebugSymbol, bool) {
	for _, symbol := range info.Symbols {
		if symbol.containsFunction(address) {
			return symbol, true
		}
	}
	return DebugSymbol{}, false
}

func (symbol DebugSymbol) containsFunction(address uint32) bool {
	return symbol.Ident == debugFunction &&
		symbol.CodeStart <= address &&
		address < symbol.CodeEnd &&
		isUserSymbol(symbol.Name)
}

func isUserSymbol(name string) bool {
	return name == "" || name[0] != '@'
}
