# syntax = docker/dockerfile:1.0-experimental

FROM golang:1.16.3-alpine

WORKDIR /app

RUN apk add --update --no-cache gcc build-base

RUN --mount=type=cache,target=/root/.cache/go-build \
  go install github.com/cosmtrek/air@v1.26.0

RUN --mount=type=cache,target=/root/.cache/go-build \
  go install github.com/go-delve/delve/cmd/dlv@v1.6.0

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod/cache \
  go mod download

ENTRYPOINT [ "air" ]
CMD [ "-c", "dev/.air.toml" ]
