# Compatibility

The runtime loads 32-bit AMX images produced by the modern open.mp Pawn compiler. Non-32-bit images and obsolete direct-native execution paths are rejected.

Focused fixtures cover headers, metadata, memory, compact encoding, suspended execution, native dispatch, decoding, and opcode execution. The [opcode table](amx-opcode-support.md) records the current interpreter coverage.

## Gaps

The repository still needs a larger set of checked-in programs compiled from representative Pawn source and an optional parity harness against the canonical C runtime. Any generated fixture should pin its compiler version and remain a development dependency; the Go library itself stays free of cgo.
