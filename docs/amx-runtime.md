# AMX runtime

## Instrumentation

`SetInstrumentationHook` emits instruction, public, native, and exception events. Hooks are synchronous and can stop execution by returning an error. The instruction path performs no event work when instrumentation is disabled.

The pure-Go AMX runtime currently loads AMX headers, publics, natives, data
memory, and code segments. Compact-encoded AMX code/data blocks are expanded at
load time. It decodes and executes the active 32-bit compiler instruction set:

- constants and register moves
- integer arithmetic, division, bitwise ops, shifts, and comparisons
- signed and unsigned branches
- public entry execution
- basic stack frames, locals, `call`, `call.pri`, `proc`, `ret`, and `retn`
- global/local memory loads and stores
- array indexing helpers
- block memory operations: `movs`, `cmps`, and `fill`
- `switch` / `casetbl`
- `halt`, including typed nonzero halt errors
- `sysreq.c` and `sysreq.n` native dispatch
- macro push opcodes and load-both opcodes
- alignment and debug metadata opcodes
- standard Pawn floating-point helper natives

Obsolete direct-native placeholders and non-32-bit AMX images fail explicitly.
Focused loader, memory, decoder, and execution fixtures guard compatibility.

See `docs/amx-opcode-support.md` for the current coverage table.

The reusable host API, including native registration, metadata, memory,
debugging, cloning, and suspended execution, is documented in `go-amx.md`.
