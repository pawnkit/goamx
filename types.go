package amx

import "math"

const (
	CellBytes                   = 4
	FlagDebug            uint16 = 0x0002
	FlagCompact          uint16 = 0x0004
	FlagSleep            uint16 = 0x0008
	FlagNoChecks         uint16 = 0x0010
	FlagNoRelocation     uint16 = 0x0200
	FlagNoDirectNative   uint16 = 0x0400
	FlagSYSREQN          uint16 = 0x0800
	FlagNativeRegistered uint16 = 0x1000
	FlagJITCompiled      uint16 = 0x2000
	FlagBrowse           uint16 = 0x4000
	FlagRelocated        uint16 = 0x8000
)

type Cell int32

func CellFromFloat32(value float32) Cell { return Cell(int32(math.Float32bits(value))) }
func (cell Cell) Float32() float32       { return math.Float32frombits(uint32(cell)) }

type Public struct {
	Index int
	Name  string
}

type Native struct {
	Index      int
	Name       string
	Registered bool
}

type Library struct {
	Index int
	Name  string
}

type PublicVar struct {
	Index   int
	Name    string
	Address Cell
}

type Tag struct {
	Index int
	Name  string
	ID    Cell
}

type Info struct {
	Name           string
	FileVersion    uint8
	AMXVersion     uint8
	Flags          uint16
	DefinitionSize uint16
	CodeSize       int
	DataSize       int
	StackHeapSize  int
	NameLength     int
	HasMain        bool
}

type MemoryInfo struct {
	CodeBytes      int
	DataBytes      int
	StackHeapBytes int
}

type NativeContext interface {
	ReadString(addr Cell) (string, error)
	WriteString(addr Cell, value string) error
	ReadCell(addr Cell) (Cell, error)
	WriteCell(addr Cell, value Cell) error
	ReadBytes(addr Cell, size int) ([]byte, error)
	WriteBytes(addr Cell, value []byte) error
	CallPublic(name string, args ...Cell) (Cell, error)
}

type NativeFunc func(ctx NativeContext, params []Cell) (Cell, error)

type State struct {
	PRI, ALT           Cell
	HEA, STK, STP, FRM Cell
	CIP                int
}

type DebugEvent struct {
	Instruction Instruction
	State       State
}

type DebugHook func(DebugEvent) error

type Instruction struct {
	Offset int32
	Opcode Opcode
	Params []Cell
	Size   int
}
