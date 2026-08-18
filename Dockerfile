FROM golang:alpine3.10 AS builder

WORKDIR /usr/src/app

COPY go.mod ./
RUN go mod download

COPY ./ ./
RUN go build -v -o /usr/local/bin/app ./...

FROM alpine:3.10

WORKDIR /usr/src/app

COPY --from=builder ./app ./app

EXPOSE 8080
CMD ["./app"]