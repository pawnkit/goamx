// Package vm exposes the low-level AMX virtual machine.
//
// Most applications should import github.com/pawnkit/goamx and use the root
// package. Use this package when building bytecode tooling, compatibility
// harnesses, custom hosts, or integrations that need direct access to VM state,
// instruction decoding, runtime errors, debug hooks, and raw AMX memory.
package vm
