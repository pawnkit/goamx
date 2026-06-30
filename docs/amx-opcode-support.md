# AMX Opcode Support

This table tracks pure-Go AMX runtime coverage. Decode support means the bytecode
walker can size and name the instruction. Execute support means the interpreter
implements the runtime behavior. Parity-tested means coverage has been compared
with the canonical C AMX runtime.

| Opcode area | Decode | Execute | Parity-tested | Notes |
| --- | --- | --- | --- | --- |
| Core active opcode table | yes | yes | smoke | Fixed and variable instruction decoding is covered. |
| Constants and memory | yes | yes | fixture | Global, frame, indirect, byte, block, array, heap, and stack operations execute with bounds checks. |
| Integer arithmetic | yes | yes | fixture | Signed/unsigned arithmetic, division, bitwise operations, shifts, sign extension, and comparisons execute. |
| Branches and switch | yes | yes | fixture | Absolute/relative branches and `switch`/`casetbl` execute. |
| Calls and returns | yes | yes | fixture | Public calls, indirect calls, frames, `ret`, and `retn` execute. |
| Native calls | yes | yes | fixture | `sysreq.pri`, `sysreq.c`, and `sysreq.n` dispatch Go natives. |
| Macro push/load ops | yes | yes | fixture | Repeated push, macro push, swap, and load-both opcodes execute. |
| Runtime errors | n/a | yes | fixture | Bounds, divide-by-zero, invalid memory, and nonzero halt errors are typed. |
| Compact AMX encoding | n/a | yes | fixture | Compact code/data blocks are expanded before decoding. |
| Alignment/debug metadata | yes | yes | unit | Alignment and debug records have their canonical runtime no-op behavior. |
| Obsolete `sysreq.d` placeholders | yes | rejected | n/a | The upstream interpreter marks these enum slots unused. |
