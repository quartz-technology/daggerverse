FROM golang:1.23.2-alpine

WORKDIR /app

COPY go.mod go.mod
COPY go.sum go.sum
COPY main.go main.go

RUN go mod download

RUN go build -o /app/main .