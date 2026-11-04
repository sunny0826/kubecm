# Contributing guidelines

We love your input! We want to make contributing to this project as easy and transparent as possible, whether it's:

- Reporting a bug
- Discussing the current state of the code
- Submitting a fix
- Proposing new features
- Becoming a maintainer

Feel free to open [issue](https://github.com/sunny0826/kubecm/issues/new) and [pull requests](https://github.com/sunny0826/kubecm/pulls). Any feedback is highly appreciated!

Be sure to follow our [Community Guidelines](https://kubecm.cloud/en-us/contribute) when submitting your PR and interacting with other folks in this repository.

## Development environment

Install [mise](https://mise.jdx.dev/getting-started.html), then run:

```bash
mise trust
mise install
mise exec -- go version
```

The project pins Go and validation tools in [mise.toml](mise.toml).
Use `mise exec --` for development commands, or activate mise in your shell
and configure your editor to use its Go installation. `go.mod` declares the
minimum Go version; mise pins the patch version used locally and in CI.

## Testing

### Unit Tests

Run unit tests with:
```bash
mise exec -- make test
mise exec -- go test -race ./pkg/...
mise exec -- golangci-lint run ./...
mise exec -- actionlint
```

### E2E Tests

End-to-end tests validate the complete workflow of kubecm. To run e2e tests:

```bash
# Run e2e tests (builds kubecm automatically)
mise exec -- make e2e-test

# Or use the test runner script
mise exec -- ./test/run-e2e.sh
```

For more information on writing and running e2e tests, see [test/e2e/README.md](test/e2e/README.md).