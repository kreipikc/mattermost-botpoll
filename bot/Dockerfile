FROM golang:1.21

WORKDIR /bot

COPY go.mod go.sum ./

RUN go mod download

COPY . .

CMD ["go", "run", "cmd/main.go"]