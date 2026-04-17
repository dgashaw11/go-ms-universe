# go-ms-universe

Two microservices for product management:

- **Products** — REST API for creating, deleting, and listing products. Publishes lifecycle events to Kafka.
- **Notifications** — Consumes product events from Kafka and logs them.

## Prerequisites

- Go 1.26+
- Docker & Docker Compose
- golangci-lint 2.1+

## Getting Started

```sh
# Start dependencies (Postgres, Kafka, Prometheus).
make docker.up

# Copy and adjust environment.
cp .env.example .env

# Run services (in separate terminals).
make run.products
make run.notifications
```

## API

| Method | Path                    | Description              |
|--------|-------------------------|--------------------------|
| POST   | /api/v1/products        | Create a product         |
| GET    | /api/v1/products        | List products (paginated)|
| DELETE | /api/v1/products/{id}   | Delete a product         |
| GET    | /health                 | Health check             |
| GET    | /metrics                | Prometheus metrics       |
| GET    | /swagger/               | Swagger UI               |

## Testing

```sh
make test.unit      # unit tests
make test.cover     # unit tests with coverage report
```

## Linting

```sh
make lint.fix       # auto-fix formatting and tidy modules
make lint           # check for remaining issues
make modernise      # apply modern Go idioms
```

## Building

```sh
make build          # binaries in bin/

# Docker
docker build --build-arg SERVICE=products -t products .
docker build --build-arg SERVICE=notifications -t notifications .
```

## Project Structure

```
cmd/
  products/           Products service entry point + swagger spec
  notifications/      Notifications service entry point
internal/
  product/            Domain model and service layer
  httpapi/            HTTP server, router, handlers
  postgres/           PostgreSQL repository (raw SQL, pgx)
  kafka/              Kafka producer and consumer
  metrics/            Prometheus counters
  config/             Environment configuration
migrations/           PostgreSQL schema migrations
```
