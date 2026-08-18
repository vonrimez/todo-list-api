FROM golang:1.26.6-alpine3.24 AS builder

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download

COPY ./ ./
RUN go build -v -o /usr/local/bin/app ./...

FROM alpine:3.24

WORKDIR /usr/src/app

COPY --from=builder /usr/local/bin/app ./app

EXPOSE 8080
CMD ["./app"]