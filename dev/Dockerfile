# syntax = docker/dockerfile:1.3.0

FROM golang:1.16.6-alpine

WORKDIR /app

RUN apk add --update --no-cache gcc build-base

RUN --mount=type=cache,target=/root/.cache/go-build \
  go install github.com/cosmtrek/air@v1.27.3

RUN --mount=type=cache,target=/root/.cache/go-build \
  go install github.com/go-delve/delve/cmd/dlv@v1.7.0

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod/cache \
  go mod download

ENTRYPOINT [ "air" ]
CMD [ "-c", "dev/.air.toml" ]
