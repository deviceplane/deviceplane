FROM alpine:3.9
RUN apk --update add ca-certificates

FROM golang:1.13
ARG version
WORKDIR /app
COPY ./ ./
RUN GOOS=linux CGO_ENABLED=0 go build -mod vendor -ldflags "-X main.version=$version" -o ./controller ./cmd/controller

FROM scratch
COPY --from=0 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=1 /app/controller /bin/controller
ENTRYPOINT ["/bin/controller"]
