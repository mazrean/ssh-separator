# Repository Guidelines — ssh-separator

Go server that separates / multiplexes SSH connections. Uses `wire` for DI and
ships a Dockerfile for deployment.

> Agent configuration is managed via [apm](https://github.com/microsoft/apm).
> Common conventions live in `mazrean/apm-plackage/common`; Go-specific rules
> come from `mazrean/apm-plackage/go`. Run `apm install` to materialise locally.

## Build & Test

- `go test -v ./...`
- `go generate ./...` — regenerate wire bindings
- `docker compose up` — local stack
- `golangci-lint run`

## Conventions

- Specs go under `specs/`; use `mazrean/agent-skills/skills/writing-*`.
- Commit using Conventional Commits (`committing-code` skill).
- Wire-based DI: regenerate after editing injector files.
- Use the Go 1.24+ `tool` directive for build tools; see `using-go-tool-directive` skill.
