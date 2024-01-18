FROM golang:1.21.4-bullseye
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY main.go .
CMD go run main.go
