# syntax = docker/dockerfile:1

FROM golang:1.25.0-bookworm AS build

WORKDIR /app

RUN --mount=type=bind,source=go.mod,target=/app/go.mod \
  --mount=type=bind,source=go.sum,target=/app/go.sum \
  --mount=type=cache,target=/go/pkg/mod/,sharing=locked \
  go mod download -x

RUN --mount=type=cache,target=/go/pkg/mod/ \
  --mount=type=bind,target=. \
  GOOS=${TARGETOS} GOARCH=${TARGETARCH} CGO_ENABLED=0 go build -o /bin/ssh-server .

FROM gcr.io/distroless/base-debian12

WORKDIR /app
COPY --from=build /bin/ssh-server ./

ENTRYPOINT [ "/app/ssh-server" ]
