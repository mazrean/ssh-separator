# syntax = docker/dockerfile:1.0-experimental

FROM golang:1.16.3-alpine

WORKDIR /app

RUN apk add --update --no-cache gcc build-base

RUN wget https://raw.githubusercontent.com/cosmtrek/air/master/install.sh && \
  sh install.sh -b /go/bin && \
  rm install.sh

RUN --mount=type=cache,target=/root/.cache/go-build \
  go install github.com/go-delve/delve/cmd/dlv@v1.6.0

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod/cache \
  go mod download

ENTRYPOINT [ "air" ]
CMD [ "-c", "dev/.air.toml" ]
