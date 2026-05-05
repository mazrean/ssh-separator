# Repository Guidelines — ssh-separator

Go server that separates / multiplexes SSH connections. Ships a Dockerfile for deployment.

> Agent configuration is managed via [apm](https://github.com/microsoft/apm).
> Common conventions live in `mazrean/apm-plackage/common`; Go-specific rules
> come from `mazrean/apm-plackage/go`. Run `apm install` to materialise locally.

## Build & Test

- `go test -v ./...`
- `go generate ./...` — regenerate kessoku injectors
- `go tool lint ./...` — run the per-repo `tools/lint` linter
- `docker compose up` — local stack

## Conventions

- Specs go under `specs/`; use `mazrean/agent-skills/skills/writing-*`.
- Commit using Conventional Commits (`committing-code` skill).
- **DI**: compile-time via `mazrean/kessoku` (consumed via `go.mod` `tool` directive,
  invoked as `go tool kessoku $GOFILE` from `//go:generate` markers). The historical
  wire-based setup should be migrated via `go tool kessoku migrate ./...`.
- **Linter**: per-repo `tools/lint` Go module invoked as `go tool lint ./...`.
  Do not introduce `golangci-lint`.
- Use the Go 1.24+ `tool` directive for build tools; see `using-go-tool-directive` skill.
