# nexo

## Feasibility check

The framework is a Go module and should build and test cleanly when your environment matches the required Go version.

```bash
# Run the test suite
go test ./...

# Build all packages (or use `task build` if Task is installed)
go build ./...
```

If both commands succeed, the codebase is functional and the CLI compiles. For linting, install a `golangci-lint` binary built with Go 1.25.x to match the `go.mod` version.
