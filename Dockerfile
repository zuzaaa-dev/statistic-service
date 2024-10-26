FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

EXPOSE 8003

RUN go build -o main ./cmd/statistic/main/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/. .

EXPOSE 8003

CMD ["./main"]