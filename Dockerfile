FROM golang:1.16-alpine

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /roulette-service

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN go build -o roulette-service cmd/main.go

EXPOSE 8080
