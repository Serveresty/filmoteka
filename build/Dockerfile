# Используем базовый образ Golang
FROM golang:latest

ENV GO111MODULE=on

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o filmoteka ./cmd/app/

CMD ["./filmoteka"]