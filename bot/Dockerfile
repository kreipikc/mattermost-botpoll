FROM golang:1.24.1

WORKDIR /bot

COPY ./bot/go.mod ./bot/go.sum ./

RUN go mod download

COPY . .