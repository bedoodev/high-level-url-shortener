# Stage 1: Build binary
FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Giriş noktası olarak main.go build ediliyor
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main ./cmd/api

# Stage 2: Minimal runtime image
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/main .
COPY .env .

# Sağlık kontrolü (opsiyonel ama önerilir)
HEALTHCHECK --interval=10s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --spider -q http://localhost:8080/healthz || exit 1

EXPOSE 8080

CMD ["./main"]
