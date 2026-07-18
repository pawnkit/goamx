package vm

import "errors"

var (
	ErrInvalidAMX            = errors.New("invalid AMX")
	ErrUnsupportedCellSize   = errors.New("unsupported AMX cell size")
	ErrUnsupportedExecution  = errors.New("unsupported AMX execution operation")
	ErrPublicIndexOutOfRange = errors.New("public index out of range")
	ErrInvalidMemoryAccess   = errors.New("invalid AMX memory access")
	ErrInvalidInstruction    = errors.New("invalid AMX instruction")
	ErrNotSleeping           = errors.New("AMX runtime is not sleeping")
	ErrExecutionPaused       = errors.New("AMX execution paused")
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
