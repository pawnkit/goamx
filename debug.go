package amx

const (
	SymbolVariable  uint8 = 1
	SymbolReference uint8 = 2
	SymbolArray     uint8 = 3
	SymbolRefArray  uint8 = 4
	SymbolFunction  uint8 = 9
)

type DebugInfo struct {
	Files    []DebugFile
	Lines    []DebugLine
	Symbols  []DebugSymbol
	Tags     []DebugTag
	Automata []DebugAutomaton
	States   []DebugState
}

type DebugFile struct {
	Address uint32
	Name    string
}

type DebugLine struct {
	Address uint32
	Line    int32
}

type DebugDimension struct {
	Tag  int16
	Size uint32
}

type DebugSymbol struct {
	Address            uint32
	Tag                int16
	CodeStart, CodeEnd uint32
	Ident, Class       uint8
	Name               string
	Dimensions         []DebugDimension
}

type DebugTag struct {
	ID   int16
	Name string
}

type DebugAutomaton struct {
	ID      int16
	Address uint32
	Name    string
}

type DebugState struct {
	ID, Automaton int16
	Name          string
}

func (r *Runtime) DebugInfo() DebugInfo {
	source := r.vm.DebugInfo()
	info := DebugInfo{}

	for _, entry := range source.Files {
		info.Files = append(info.Files, DebugFile{Address: entry.Address, Name: entry.Name})
	}

	for _, entry := range source.Lines {
		info.Lines = append(info.Lines, DebugLine{Address: entry.Address, Line: entry.Line})
	}

	for _, entry := range source.Symbols {
		symbol := DebugSymbol{
			Address:   entry.Address,
			Tag:       entry.Tag,
			CodeStart: entry.CodeStart,
			CodeEnd:   entry.CodeEnd,
			Ident:     entry.Ident,
			Class:     entry.Class,
			Name:      entry.Name,
		}

		for _, dimension := range entry.Dimensions {
			symbol.Dimensions = append(symbol.Dimensions, DebugDimension{Tag: dimension.Tag, Size: dimension.Size})
		}

		info.Symbols = append(info.Symbols, symbol)
	}

	for _, entry := range source.Tags {
		info.Tags = append(info.Tags, DebugTag{ID: entry.ID, Name: entry.Name})
	}

	for _, entry := range source.Automata {
		info.Automata = append(info.Automata, DebugAutomaton{ID: entry.ID, Address: entry.Address, Name: entry.Name})
	}

	for _, entry := range source.States {
		info.States = append(info.States, DebugState{ID: entry.ID, Automaton: entry.Automaton, Name: entry.Name})
	}

	return info
}

func (info DebugInfo) FileAt(address uint32) (DebugFile, bool) {
	var found DebugFile
	ok := false
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
	ok := false
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
	return symbol.Ident == SymbolFunction &&
		symbol.CodeStart <= address &&
		address < symbol.CodeEnd &&
		isUserSymbol(symbol.Name)
}

func isUserSymbol(name string) bool {
	return name == "" || name[0] != '@'
}
