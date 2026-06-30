package vm

type Cell int32

type Public struct {
	Index int
	Name  string
}

type NativeFunc func(ctx NativeContext, params []Cell) (Cell, error)

type NativeContext interface {
	ReadString(addr Cell) (string, error)
	WriteString(addr Cell, value string) error
	ReadCell(addr Cell) (Cell, error)
	WriteCell(addr Cell, value Cell) error
}

// PublicCaller is implemented by runtimes that support invoking a Pawn public
// from a native callback.
type PublicCaller interface {
	CallPublic(name string, args ...Cell) (Cell, error)
}

type MemorySnapshotter interface {
	SnapshotMemory() []byte
	RestoreMemory(snapshot []byte) error
}

type InstructionLimiter interface {
	SetInstructionLimit(limit int)
}

type DebugLocator interface {
	DebugLocation(cip Cell) (file string, line int, function string, ok bool)
}

type CoverageLocation struct {
	File     string
	Line     int
	Function string
}

type CoverageInstrumenter interface {
	CoverageLocations() []CoverageLocation
	SetInstructionObserver(observer func(cip Cell))
}
