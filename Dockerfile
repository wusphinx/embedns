# syntax=docker/dockerfile:1.4
FROM golang:1.20.5-bullseye AS build

WORKDIR /opt

COPY go.mod .
COPY go.sum .

COPY main.go .

RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod go build -o app

FROM busybox:1.35.0-glibc

WORKDIR /

COPY --from=build /opt /

EXPOSE 53
EXPOSE 2379

ENTRYPOINT ["/app"]
