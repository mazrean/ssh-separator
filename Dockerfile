# syntax = docker/dockerfile:1

FROM golang:1.24.0-bookworm AS build

WORKDIR /app

RUN --mount=type=bind,source=go.mod,target=/app/go.mod,readonly \\
  --mount=type=bind,source=go.sum,target=/app/go.sum,readonly \\
  --mount=type=cache,target=/go/pkg/mod/cache \\
  go mod download -x

RUN --mount=type=bind,source=.,target=/app,readonly \\
  --mount=type=cache,target=/root/.cache/go-build \\
  --mount=type=cache,target=/go/pkg/mod/cache \\
  GOOS=${TARGETOS} GOARCH=${TARGETARCH} CGO_ENABLED=0 go build -o ssh-server .

FROM gcr.io/distroless/base-debian12

WORKDIR /app
COPY --from=build /app/ssh-server ./

ENTRYPOINT [ "/app/ssh-server" ]
