# goamx

`goamx` is a pure-Go loader and runtime for AMX/Pawn bytecode.

It can load AMX images, inspect public functions/native declarations/debug
metadata, execute public functions, register Go native callbacks, read and
write AMX memory, clone runtime state, and decode instructions.

The root package is the friendly host API. The `vm` subpackage exposes the
lower-level virtual machine for bytecode tooling and advanced embedding.

## Install

```sh
go get github.com/pawnkit/goamx
```

## Root package example

```go
package main

import (
	"fmt"
	"log"

	"github.com/pawnkit/goamx"
)

func main() {
	runtime, err := amx.LoadFile("gamemode.amx")
	if err != nil {
		log.Fatal(err)
	}
	defer runtime.Close()

	public, ok := runtime.FindPublic("OnGameModeInit")
	if !ok {
		log.Fatal("public not found")
	}

	value, err := runtime.ExecPublic(public.Index)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(value)
}
```

## Low-level VM example

```go
package main

import (
	"log"

	"github.com/pawnkit/goamx/vm"
)

func main() {
	machine, err := vm.LoadFile("gamemode.amx")
	if err != nil {
		log.Fatal(err)
	}

	instructions, err := machine.Decode()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("decoded %d instructions", len(instructions))
}
```

## Development

```sh
task fmt
task check
task lint
```

Run `go test -race ./...` before changing runtime state, hooks, or native dispatch. See the [compatibility policy](docs/compatibility.md) for supported AMX images.

The public API is split into two packages:

- `github.com/pawnkit/goamx`: ergonomic runtime API for hosts.
- `github.com/pawnkit/goamx/vm`: lower-level VM, loader, decoder, and memory API.

## Contributing

This is community tooling built in spare time. Bug reports, small fixes, and
AMX compatibility fixtures are welcome. See [CONTRIBUTING.md](CONTRIBUTING.md).
