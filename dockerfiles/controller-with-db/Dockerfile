FROM alpine:3.9
RUN apk --update add ca-certificates

FROM node:10
WORKDIR /app
COPY ./ ./
WORKDIR /app/ui
RUN npm i
RUN npm run-script build

FROM golang:1.13
ARG version
WORKDIR /app
COPY ./ ./
RUN mkdir ui/build
COPY --from=1 /app/ui/dist ui/build
RUN go get github.com/rakyll/statik
RUN statik -src=./ui/build -dest=./pkg
RUN GOOS=linux CGO_ENABLED=0 go build -mod vendor -ldflags "-X main.version=$version" -o ./controller ./cmd/controller

FROM mysql:5.6
COPY --from=0 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=2 /app/controller /bin/controller
COPY ./pkg/controller/store/mysql/schema.sql /docker-entrypoint-initdb.d/
COPY ./dockerfiles/controller-with-db/entry.sh /entry.sh
VOLUME /var/lib/mysql /var/log/mysql
ENTRYPOINT ["/entry.sh"]
