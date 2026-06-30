package vm

import (
	"encoding/binary"
	"fmt"
)

const (
	debugHeaderSize = 22
	debugMagic      = 0xf1ef
	debugFunction   = 9
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

func parseDebugInfo(data []byte, header header) (DebugInfo, error) {
	if header.Flags&amxFlagDebug == 0 {
		return DebugInfo{}, nil
	}

	chunk, counts, err := debugChunk(data, header)
	if err != nil {
		return DebugInfo{}, err
	}
	cursor := debugCursor{data: chunk, offset: debugHeaderSize}

	var info DebugInfo
	if info.Files, err = parseDebugFiles(&cursor, counts.files); err != nil {
		return DebugInfo{}, err
	}
	if info.Lines, err = parseDebugLines(&cursor, counts.lines); err != nil {
		return DebugInfo{}, err
	}
	if info.Symbols, err = parseDebugSymbols(&cursor, counts.symbols); err != nil {
		return DebugInfo{}, err
	}
	if info.Tags, err = parseDebugTags(&cursor, counts.tags); err != nil {
		return DebugInfo{}, err
	}
	if info.Automata, err = parseDebugAutomata(&cursor, counts.automata); err != nil {
		return DebugInfo{}, err
	}
	if info.States, err = parseDebugStates(&cursor, counts.states); err != nil {
		return DebugInfo{}, err
	}
	return info, nil
}

type debugCounts struct {
	files, lines, symbols, tags, automata, states int
}

func debugChunk(data []byte, header header) ([]byte, debugCounts, error) {
	start := int(header.Size)
	if start < 0 || start+debugHeaderSize > len(data) {
		return nil, debugCounts{}, fmt.Errorf("%w: missing debug header", ErrInvalidAMX)
	}
	chunk := data[start:]
	size := int(binary.LittleEndian.Uint32(chunk[0:4]))
	if binary.LittleEndian.Uint16(chunk[4:6]) != debugMagic || size < debugHeaderSize || size > len(chunk) {
		return nil, debugCounts{}, fmt.Errorf("%w: invalid debug chunk", ErrInvalidAMX)
	}

	values := [6]int{}
	for i := range values {
		values[i] = int(int16(binary.LittleEndian.Uint16(chunk[10+i*2 : 12+i*2])))
		if values[i] < 0 {
			return nil, debugCounts{}, fmt.Errorf("%w: negative debug table count", ErrInvalidAMX)
		}
	}
	return chunk[:size], debugCounts{
		files: values[0], lines: values[1], symbols: values[2],
		tags: values[3], automata: values[4], states: values[5],
	}, nil
}

type debugCursor struct {
	data   []byte
	offset int
}

func (cursor *debugCursor) take(size int) ([]byte, error) {
	if size < 0 || cursor.offset > len(cursor.data)-size {
		return nil, fmt.Errorf("%w: truncated debug table", ErrInvalidAMX)
	}
	value := cursor.data[cursor.offset : cursor.offset+size]
	cursor.offset += size
	return value, nil
}

func (cursor *debugCursor) uint8() (uint8, error) {
	value, err := cursor.take(1)
	if err != nil {
		return 0, err
	}
	return value[0], nil
}

func (cursor *debugCursor) int16() (int16, error) {
	value, err := cursor.take(2)
	if err != nil {
		return 0, err
	}
	return int16(binary.LittleEndian.Uint16(value)), nil
}

func (cursor *debugCursor) uint32() (uint32, error) {
	value, err := cursor.take(4)
	if err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint32(value), nil
}

func (cursor *debugCursor) name() (string, error) {
	start := cursor.offset
	for cursor.offset < len(cursor.data) && cursor.data[cursor.offset] != 0 {
		cursor.offset++
	}
	if cursor.offset >= len(cursor.data) {
		return "", fmt.Errorf("%w: unterminated debug name", ErrInvalidAMX)
	}
	name := string(cursor.data[start:cursor.offset])
	cursor.offset++
	return name, nil
}

func parseDebugFiles(cursor *debugCursor, count int) ([]DebugFile, error) {
	entries := make([]DebugFile, 0, count)
	for range count {
		address, err := cursor.uint32()
		if err != nil {
			return nil, err
		}
		name, err := cursor.name()
		if err != nil {
			return nil, err
		}
		entries = append(entries, DebugFile{Address: address, Name: name})
	}
	return entries, nil
}

func parseDebugLines(cursor *debugCursor, count int) ([]DebugLine, error) {
	entries := make([]DebugLine, 0, count)
	for range count {
		address, err := cursor.uint32()
		if err != nil {
			return nil, err
		}
		line, err := cursor.uint32()
		if err != nil {
			return nil, err
		}
		entries = append(entries, DebugLine{Address: address, Line: int32(line)})
	}
	return entries, nil
}

func parseDebugSymbols(cursor *debugCursor, count int) ([]DebugSymbol, error) {
	entries := make([]DebugSymbol, 0, count)
	for range count {
		symbol, dimensions, err := parseDebugSymbol(cursor)
		if err != nil {
			return nil, err
		}
		for range dimensions {
			tag, err := cursor.int16()
			if err != nil {
				return nil, err
			}
			size, err := cursor.uint32()
			if err != nil {
				return nil, err
			}
			symbol.Dimensions = append(symbol.Dimensions, DebugDimension{Tag: tag, Size: size})
		}
		entries = append(entries, symbol)
	}
	return entries, nil
}

func parseDebugSymbol(cursor *debugCursor) (DebugSymbol, int, error) {
	address, err := cursor.uint32()
	if err != nil {
		return DebugSymbol{}, 0, err
	}
	tag, err := cursor.int16()
	if err != nil {
		return DebugSymbol{}, 0, err
	}
	codeStart, err := cursor.uint32()
	if err != nil {
		return DebugSymbol{}, 0, err
	}
	codeEnd, err := cursor.uint32()
	if err != nil {
		return DebugSymbol{}, 0, err
	}
	ident, err := cursor.uint8()
	if err != nil {
		return DebugSymbol{}, 0, err
	}
	class, err := cursor.uint8()
	if err != nil {
		return DebugSymbol{}, 0, err
	}
	dimensions, err := cursor.int16()
	if err != nil {
		return DebugSymbol{}, 0, err
	}
	if dimensions < 0 {
		return DebugSymbol{}, 0, fmt.Errorf("%w: negative symbol dimensions", ErrInvalidAMX)
	}
	name, err := cursor.name()
	if err != nil {
		return DebugSymbol{}, 0, err
	}
	return DebugSymbol{
		Address: address, Tag: tag, CodeStart: codeStart, CodeEnd: codeEnd,
		Ident: ident, Class: class, Name: name,
	}, int(dimensions), nil
}

func parseDebugTags(cursor *debugCursor, count int) ([]DebugTag, error) {
	entries := make([]DebugTag, 0, count)
	for range count {
		id, err := cursor.int16()
		if err != nil {
			return nil, err
		}
		name, err := cursor.name()
		if err != nil {
			return nil, err
		}
		entries = append(entries, DebugTag{ID: id, Name: name})
	}
	return entries, nil
}

func parseDebugAutomata(cursor *debugCursor, count int) ([]DebugAutomaton, error) {
	entries := make([]DebugAutomaton, 0, count)
	for range count {
		id, err := cursor.int16()
		if err != nil {
			return nil, err
		}
		address, err := cursor.uint32()
		if err != nil {
			return nil, err
		}
		name, err := cursor.name()
		if err != nil {
			return nil, err
		}
		entries = append(entries, DebugAutomaton{ID: id, Address: address, Name: name})
	}
	return entries, nil
}

func parseDebugStates(cursor *debugCursor, count int) ([]DebugState, error) {
	entries := make([]DebugState, 0, count)
	for range count {
		id, err := cursor.int16()
		if err != nil {
			return nil, err
		}
		automaton, err := cursor.int16()
		if err != nil {
			return nil, err
		}
		name, err := cursor.name()
		if err != nil {
			return nil, err
		}
		entries = append(entries, DebugState{ID: id, Automaton: automaton, Name: name})
	}
	return entries, nil
}
