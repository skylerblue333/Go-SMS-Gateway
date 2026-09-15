# Contributing to Sky SMS Envelope

Thanks for contributing to the Go SMS envelope validator used in the SKYCOIN4444 ecosystem.

## Requirements

- Go 1.26.x
- Docker for container and smoke-test validation

## Development setup

No Node.js, npm, Python, or Rust toolchain is required for this repository.

```bash
go version
go test ./...
```

## Required verification

Run the same native Go checks enforced by CI before opening a pull request:

```bash
test -z "$(gofmt -l .)"
go vet ./...
go test -v ./...
go test -race ./...
CGO_ENABLED=0 go build -trimpath -o sky-sms-envelope .
```

For vulnerability scanning, install and run the Go vulnerability tool:

```bash
go install golang.org/x/vuln/cmd/govulncheck@latest
govulncheck ./...
```

Container verification:

```bash
docker build -t sky-sms-envelope:local .
docker run --rm -p 8080:8080 sky-sms-envelope:local
```

Then verify the health endpoint from another shell:

```bash
curl --fail http://127.0.0.1:8080/healthz
```

## Code style

- Format Go source with `gofmt`.
- Keep `go vet` clean.
- Prefer standard-library solutions unless an external dependency has a clear benefit.
- Add deterministic tests for validation rules, malformed input, boundary conditions, and failure paths.
- Keep provider dispatch, credentials, consent/compliance workflows, and durable delivery concerns outside this validator unless their implementation is explicitly in scope and tested.

## Commits and pull requests

Use small, reviewable changes and descriptive Conventional Commit-style messages where practical, for example:

- `feat(validation): add sender identifier rule`
- `fix(http): reject trailing JSON payload data`
- `test(validation): cover UTF-8 byte limits`
- `docs: clarify provider integration boundary`

Pull requests should explain the behavior changed, the tests added or updated, and any remaining limitations.

## License

See `LICENSE`.
