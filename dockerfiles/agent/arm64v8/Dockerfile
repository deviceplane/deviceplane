FROM ubuntu:18.04
RUN apt-get update && apt-get install -y qemu-user-static

FROM arm64v8/alpine:3.10
COPY --from=0 /usr/bin/qemu-aarch64-static /usr/bin/
RUN apk -u add ca-certificates util-linux

FROM golang:1.13
ARG version
WORKDIR /app
COPY ./ ./
RUN GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -mod vendor -ldflags "-s -w -X main.version=$version" -o ./agent ./cmd/agent

FROM ubuntu:18.04
RUN apt-get update && apt-get install -y qemu-user-static

FROM scratch
COPY --from=1 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=1 /usr/bin/nsenter /usr/bin/
COPY --from=1 /lib/ld-musl-aarch64.so.1 /lib/
COPY --from=2 /app/agent /bin/agent
ENTRYPOINT ["/bin/agent"]
