FROM golang:1.24

RUN apk add --no-cache curl git

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .
RUN go install github.com/air-verse/air@latest

# Healthcheck (Kubernetes & Docker support)
HEALTHCHECK --interval=10s --timeout=3s --start-period=5s --retries=3 \
  CMD curl -f http://localhost:8080/healthz || exit 1

# Use Air for dev (change for production later)
CMD ["air"]