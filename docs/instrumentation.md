# Instrumentation API

Install `Runtime.SetInstrumentationHook` to observe instructions, public calls, native calls, and exceptions. Each instruction event includes registers and available source location data. Returning an error stops execution.

Hooks run synchronously on the runtime goroutine. A runtime is not safe for concurrent execution. Collectors must synchronize any data shared across runtimes.
