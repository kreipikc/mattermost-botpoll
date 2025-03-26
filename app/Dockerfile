FROM golang:1.24.1

WORKDIR /app

COPY ./app/go.mod ./app/go.sum ./

RUN go mod download

COPY . .