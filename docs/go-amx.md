# Pure-Go AMX Runtime

`goamx` is a reusable, cgo-free AMX loader and interpreter. The public package
does not depend on a Pawn compiler, a server process, or any application-level
test runner.

## Loading and execution

```go
runtime, err := amx.LoadFile("program.amx")
if err != nil {
    return err
}
defer runtime.Close()

public, ok := runtime.FindPublic("calculate")
if !ok {
    return errors.New("calculate is not public")
}
value, err := runtime.ExecPublic(public.Index, 10, 20)
```

The runtime supports main/public execution, nested public callbacks, sleep and
continuation, configurable instruction limits, and instruction-level debug
hooks.

## Host natives

```go
err = runtime.RegisterNative("host_add", func(ctx amx.NativeContext, params []amx.Cell) (amx.Cell, error) {
    if len(params) != 2 {
        return 0, errors.New("host_add expects two arguments")
    }
    return params[0] + params[1], nil
})
```

Native contexts provide checked cell, byte, and packed/unpacked string memory
access plus nested public invocation. A native must be declared in the loaded
image before it can be registered.

## Metadata and memory

The public API exposes:

- image versions, flags, segment sizes, and name limits
- public functions, native declarations, public variables, tags, and libraries
- decoded typed instructions and opcode metadata
- debug files, lines, symbols, dimensions, tags, automata, and states
- checked byte/cell/string access
- heap allocation and release, including array and string helpers
- runtime cloning, execution reset, and complete data-memory reset
- host user data keyed by an integer tag

## Package Layout

Application code should usually import `github.com/pawnkit/goamx` and use
the root package. It provides typed cells, native callbacks, string/memory
helpers, public lookup by name, user data, debug hooks, cloning, and translated
runtime errors.

Advanced tooling can import `github.com/pawnkit/goamx/vm` directly. That
subpackage exposes the low-level AMX virtual machine, direct loader, instruction
decoder, opcode metadata, debug tables, raw memory access, execution state, and
runtime error values.
