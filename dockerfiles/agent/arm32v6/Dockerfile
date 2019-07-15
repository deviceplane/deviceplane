FROM arm32v6/alpine:3.10
RUN apk -u add ca-certificates util-linux

FROM golang:1.12
ARG version
WORKDIR /app
COPY ./ ./
RUN GOOS=linux GOARCH=arm CGO_ENABLED=0 go build -ldflags "-s -w -X main.version=$version" -o ./agent ./cmd/agent

FROM scratch
COPY --from=0 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=0 /usr/bin/nsenter /usr/bin/
COPY --from=0 /lib/ld-musl-armhf.so.1 /lib/
COPY --from=1 /app/agent /bin/agent
ENTRYPOINT ["/bin/agent"]
