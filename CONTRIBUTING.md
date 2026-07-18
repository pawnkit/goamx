# Contributing

PawnKit is maintained by volunteers, so reviews may take a little time.

Contributions are welcome, including small fixes, AMX fixtures, and clearer
documentation. If a change affects runtime behavior, include the smallest AMX
or Go test that shows it.

Run the local checks before opening a pull request:

```sh
task check
CGO_ENABLED=1 go test -race ./...
```

Keep the root package useful for normal host applications. Low-level VM details
belong in `vm`. Compatibility claims should come from a fixture or a repeatable
compiler test, not guesswork.
