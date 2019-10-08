FROM golang:1.13
ARG version
WORKDIR /app
COPY ./ ./
RUN GOOS=linux CGO_ENABLED=0 go build -mod vendor -ldflags "-s -w -X main.version=$version" -o ./deviceplane ./cmd/deviceplane

FROM alpine:3.9
RUN apk --update add git openssh tar gzip ca-certificates \
  bash curl
COPY --from=0 /app/deviceplane /bin/deviceplane
ENTRYPOINT ["/bin/deviceplane"]
