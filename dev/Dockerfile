# syntax = docker/dockerfile:1.0-experimental

FROM golang:1.16.3-alpine

WORKDIR /app

RUN apk add --update --no-cache gcc build-base

RUN wget https://raw.githubusercontent.com/cosmtrek/air/master/install.sh && \
  sh install.sh -b /go/bin && \
  rm install.sh

COPY dev/air.toml ./

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod/cache \
  go mod download

ENTRYPOINT [ "air" ]
