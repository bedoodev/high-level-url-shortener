# 🧠 High-Level URL Shortener (Go + Redis + PostgreSQL + Kafka + Docker)

This project is a **high-level, production-ready URL shortener** built in **pure Go**, with support for:

- Shortening and redirecting URLs
- DDoS protection with in-memory rate limiting
- Redis-backed caching for high-throughput URL resolution
- In-memory hot key cache for ultra-low-latency access
- Swagger documentation for full API insight
- Clean, modular project structure for scale

---

## ⚙️ Tech Stack

| Layer           | Technology                   |
| --------------- | ---------------------------- |
| Language        | Go (1.21+)                   |
| Database        | PostgreSQL (via GORM)        |
| Cache Layer     | Redis + sync.Map (hot cache) |
| Message Queue   | Kafka                        |
| API Framework   | Chi                          |
| Logging         | Uber Zap                     |
| Swagger Docs    | Swaggo + OpenAPI             |
| Dev Environment | Docker + Docker Compose      |

---

## 🛡️ Features

- `POST /shorten` – create a shortened URL
- `GET /{code}` – redirect to original URL
- `GET /ping` – health check
- **Rate limiting** – per IP, fully in-memory, no third-party lib required
- **Redis cache** – for fast slug resolution
- **In-memory hot cache** – protects Redis from hot key overload
- **Swagger UI** – view and test the API: [`/swagger/index.html`](http://localhost:8080/swagger/index.html)
- **Modular structure** with separation of concerns: `api`, `service`, `repository`, `middleware`, `cache`, `config`

---

## 🚀 Getting Started

### Prerequisites

- Docker and Docker Compose installed on your machine
- Git (to clone the repository)

### Setup

1. Clone the repository:

```bash
git clone <repository-url>
cd url-shortener
```

2. Create a .env file in the root directory with the following content:

```bash
# Database Configuration
DB_HOST=postgres
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=urlshortener
DB_PORT=5432

# Redis Configuration
REDIS_ADDR=redis:6379
```

3. Start the services using Docker Compose:

```bash
docker-compose -f docker-compose.yaml -f docker-compose.kafka.yaml up -d
```

4. Verify the service is running:

```bash
curl http://localhost:8080/ping
```

You should receive a successful response(pong).

### Accessing the API

- API Documentation : Open [swagger](http://localhost:8080/swagger/index.html) in your browser to view and test the API using Swagger UI.

#### Shortening a URL

```bash
curl -X POST http://localhost:8080/shorten \
  -H "Content-Type: application/json" \
  -d '{"url":"https://example.com"}'
```

#### Accessing a shortened URL :

Open `http://localhost:8080/{code}` in your browser, where {code} is the short code returned from the shorten endpoint.

#### Get analytic for a shortened URL

```bash
curl -X 'GET' \
  'http://localhost:8080/analytics/{code}' \
  -H 'accept: application/json'
```

#### Get top shortened URL

```bash
curl -X 'GET' \
  'http://localhost:8080/analytics/top?limit={n}' \
  -H 'accept: application/json'
```

## 🧪 Development

The project uses Air for hot reloading during development. Any changes you make to the code will automatically rebuild and restart the application.

## 📊 Architecture

The application follows a clean architecture pattern with the following components:

- API Layer : Handles HTTP requests and responses
- Service Layer : Contains business logic
- Repository Layer : Manages data persistence
- Cache Layer : Provides fast access to frequently used data
- Middleware : Implements cross-cutting concerns like rate limiting

## 🔒 Security

The application includes built-in rate limiting to protect against DDoS attacks. The rate limiter is implemented in-memory and does not require any external dependencies.

## 🔍 Monitoring

The application includes a health check endpoint at /ping that can be used for monitoring and health checks.

## 📝 License

This project is licensed under the MIT License - see the LICENSE file for details.
