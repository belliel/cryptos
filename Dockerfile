FROM golang:1.22-alpine as builder

WORKDIR /builder

COPY . .

RUN go build -o ./build/cryptos cmd/app/main.go

FROM alpine:latest as application

WORKDIR /app

COPY --from=builder /builder/build/cryptos .

RUN ["./cryptos"]
