package amx

import "github.com/pawnkit/goamx/vm"

var (
	ErrInvalidAMX            = vm.ErrInvalidAMX
	ErrUnsupportedCellSize   = vm.ErrUnsupportedCellSize
	ErrUnsupportedExecution  = vm.ErrUnsupportedExecution
	ErrPublicIndexOutOfRange = vm.ErrPublicIndexOutOfRange
	ErrInvalidMemoryAccess   = vm.ErrInvalidMemoryAccess
	ErrInvalidInstruction    = vm.ErrInvalidInstruction
	ErrNotSleeping           = vm.ErrNotSleeping
	ErrExecutionPaused       = vm.ErrExecutionPaused
)

type RuntimeErrorCode string

const (
	RuntimeErrorBounds       RuntimeErrorCode = "bounds"
	RuntimeErrorDivideByZero RuntimeErrorCode = "divide_by_zero"
	RuntimeErrorHalt         RuntimeErrorCode = "halt"
	RuntimeErrorSleep        RuntimeErrorCode = "sleep"
	RuntimeErrorDomain       RuntimeErrorCode = "domain"
)

type RuntimeError struct {
	Code    RuntimeErrorCode
	Message string
	CIP     Cell
}

func (e RuntimeError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return string(e.Code)
}
