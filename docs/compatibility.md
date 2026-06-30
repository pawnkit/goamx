# Compatibility

The runtime targets AMX images produced by the modern open.mp Pawn compiler with
32-bit cells. It rejects non-32-bit images and obsolete direct-native execution
paths explicitly.

Compatibility is currently guarded by unit tests that construct focused AMX
images for loader behavior, metadata tables, memory operations, compact
encoding, sleep/continuation, native dispatch, opcode decoding, and opcode
execution.

Useful future additions for this standalone repository:

- checked-in compiled AMX fixtures generated from representative Pawn source
- a small source corpus derived from the upstream open.mp compiler tests
- an optional parity harness against the canonical C AMX runtime
- CI that regenerates or verifies fixtures with a pinned compiler version

Those additions should remain development-time checks only; the library itself
is intended to stay pure Go and cgo-free.
