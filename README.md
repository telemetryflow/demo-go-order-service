# Order-Service

<div align="center">
  <picture>
    <source media="(prefers-color-scheme: dark)" srcset="https://github.com/telemetryflow/.github/raw/main/docs/assets/tfo-logo-sdk-dark.svg">
    <source media="(prefers-color-scheme: light)" srcset="https://github.com/telemetryflow/.github/raw/main/docs/assets/tfo-logo-sdk-light.svg">
    <img src="https://github.com/telemetryflow/.github/raw/main/docs/assets/tfo-logo-sdk-light.svg" alt="TelemetryFlow Logo" width="80%">
  </picture>

[![Version](https://img.shields.io/badge/Version-1.1.2-orange.svg)](CHANGELOG.md)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8?logo=go)](https://golang.org/)
[![OpenTelemetry](https://img.shields.io/badge/OTLP-100%25%20Compliant-success?logo=opentelemetry)](https://opentelemetry.io/)
[![Docker](https://img.shields.io/badge/Docker-Ready-2496ED?style=flat&logo=docker)](https://hub.docker.com/r/telemetryflow/telemetryflow-sdk)

</div>

<p align="center">
<strong>[GENERATED TelemetryFlow SDK]</strong> Order-Service - RESTful API with DDD + CQRS Pattern
</p>

---



## Architecture

This project follows **Domain-Driven Design (DDD)** with **CQRS (Command Query Responsibility Segregation)** pattern.

```
Order-Service/
├── cmd/
│   └── api/                    # Application entry point
├── internal/
│   ├── domain/                 # Domain Layer (Core Business Logic)
│   │   ├── entity/             # Domain entities
│   │   ├── repository/         # Repository interfaces
│   │   └── valueobject/        # Value objects
│   ├── application/            # Application Layer (Use Cases)
│   │   ├── command/            # Commands (write operations)
│   │   ├── query/              # Queries (read operations)
│   │   ├── handler/            # Command & Query handlers
│   │   └── dto/                # Data Transfer Objects
│   └── infrastructure/         # Infrastructure Layer
│       ├── persistence/        # Database implementations
│       ├── http/               # HTTP server & handlers
│       └── config/             # Configuration
├── pkg/                        # Shared packages
├── telemetry/                  # TelemetryFlow integration
├── config/                     # Service configurations
│   └── otel/                   # OpenTelemetry Collector config
├── docs/                       # Documentation
│   ├── api/                    # OpenAPI/Swagger specs
│   ├── diagrams/               # ERD, DFD diagrams
│   └── postman/                # Postman collections
├── migrations/                 # Database migrations
├── configs/                    # Application configuration files
└── tests/                      # Tests
    ├── unit/
    ├── integration/
    └── e2e/
```

## Quick Start

### Prerequisites

- Go 1.22+
- PostgreSQL 16+
- Docker & Docker Compose (recommended)

### Setup

1. Clone the repository

2. Copy environment file:
   ```bash
   cp .env.example .env
   ```

3. Edit `.env` with your configuration

4. Install dependencies:
   ```bash
   make deps
   ```

5. Run migrations:
   ```bash
   make migrate-up
   ```

6. Start the server:
   ```bash
   make run
   ```

## Docker Compose

The easiest way to run the service with all dependencies:

```bash
# Start all services (PostgreSQL + API + OpenTelemetry Collector)
make docker-compose-up

# Or use profiles for selective startup
docker compose --profile all up -d        # Start everything
docker compose --profile db up -d         # Start only PostgreSQL
docker compose --profile app up -d        # Start only API
docker compose --profile monitoring up -d # Start only OTEL Collector

# Stop all services
make docker-compose-down

# View logs
docker logs -f order_service_api
docker logs -f order_service_postgres
docker logs -f order_service_otel
```

### Services

| Service | Container | Port | Description |
|---------|-----------|------|-------------|
| PostgreSQL | `order_service_postgres` | 5432 | Database |
| API | `order_service_api` | 8080 | RESTful API |
| OTEL Collector | `order_service_otel` | 4317, 4318, 8889, 13133, 55679, 1777 | OpenTelemetry Collector |
| Jaeger | `order_service_jaeger` | 16686 | Distributed Tracing UI |

### OTEL Collector Ports

| Port  | Protocol | Description           |
| ----- | -------- | --------------------- |
| 4317  | gRPC     | OTLP gRPC (v1 & v2)   |
| 4318  | HTTP     | OTLP HTTP (v1 & v2)   |
| 8889  | HTTP     | Prometheus metrics    |
| 13133 | HTTP     | Health check          |
| 55679 | HTTP     | zPages (debugging)    |
| 1777  | HTTP     | pprof (profiling)     |

### OTLP Endpoints (Dual Ingestion)

The collector supports both TelemetryFlow (v2) and OTEL Community (v1) endpoints:

**TelemetryFlow Platform (Recommended):**

```text
POST http://localhost:4318/v2/traces
POST http://localhost:4318/v2/metrics
POST http://localhost:4318/v2/logs
```

**OTEL Community (Backwards Compatible):**

```text
POST http://localhost:4318/v1/traces
POST http://localhost:4318/v1/metrics
POST http://localhost:4318/v1/logs
```

**gRPC:** `localhost:4317` (both v1 and v2)

### Network Configuration

All services run on a custom Docker network `order_service_net` with subnet `172.152.0.0/16`:

| Service | IP Address |
|---------|------------|
| API | 172.152.152.10 |
| PostgreSQL | 172.152.152.20 |
| OTEL Collector | 172.152.152.30 |

## Development

### Running locally

```bash
# Build and run
make run

# Run with hot reload
make dev

# Run tests
make test

# Build binary
make build
```

### Adding a new entity

Use the TelemetryFlow RESTful API Generator:

```bash
telemetryflow-restapi entity -n Product -f 'name:string,price:float64,stock:int'
```

This generates:
- Domain entity
- Repository interface & implementation
- CQRS commands & queries
- HTTP handlers
- Database migration

## API Documentation

| Documentation | Location |
|--------------|----------|
| OpenAPI Spec | `docs/api/openapi.yaml` |
| Swagger JSON | `docs/api/swagger.json` |
| ERD Diagram | `docs/diagrams/ERD.md` |
| DFD Diagram | `docs/diagrams/DFD.md` |
| Postman Collection | `docs/postman/collection.json` |

### API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/health` | Health check |
| GET | `/api/v1/orders` | List all orders |
| POST | `/api/v1/orders` | Create order |
| GET | `/api/v1/orders/:id` | Get order by ID |
| PUT | `/api/v1/orders/:id` | Update order |
| DELETE | `/api/v1/orders/:id` | Delete order |
| GET | `/api/v1/orderitems` | List all order items |
| POST | `/api/v1/orderitems` | Create order item |
| GET | `/api/v1/orderitems/:id` | Get order item by ID |
| PUT | `/api/v1/orderitems/:id` | Update order item |
| DELETE | `/api/v1/orderitems/:id` | Delete order item |

## Configuration

Configuration is loaded from environment variables and `.env` file.

### Application Configuration

| Variable | Description | Default |
|----------|-------------|---------|
| `SERVER_PORT` | HTTP server port | `8080` |
| `SERVER_READ_TIMEOUT` | Read timeout | `15s` |
| `SERVER_WRITE_TIMEOUT` | Write timeout | `15s` |
| `ENV` | Environment (development/production) | `development` |

### Database Configuration

| Variable | Description | Default |
|----------|-------------|---------|
| `DB_DRIVER` | Database driver | `postgres` |
| `DB_HOST` | Database host | `localhost` |
| `DB_PORT` | Database port | `5432` |
| `DB_NAME` | Database name | `orders` |
| `DB_USER` | Database user | `postgres` |
| `DB_PASSWORD` | Database password | - |
| `DB_SSL_MODE` | SSL mode | `disable` |
| `DB_MAX_OPEN_CONNS` | Max open connections | `25` |
| `DB_MAX_IDLE_CONNS` | Max idle connections | `5` |
| `DB_CONN_MAX_LIFETIME` | Connection max lifetime | `5m` |

### JWT Configuration

| Variable | Description | Default |
|----------|-------------|---------|
| `JWT_SECRET` | JWT signing secret | - |
| `JWT_REFRESH_SECRET` | JWT refresh secret | - |
| `JWT_EXPIRATION` | Token expiration | `24h` |
| `JWT_REFRESH_EXPIRATION` | Refresh token expiration | `168h` |

### TelemetryFlow / OpenTelemetry Configuration

| Variable | Description | Default |
|----------|-------------|---------|
| `TELEMETRYFLOW_API_KEY_ID` | TelemetryFlow API Key ID | - |
| `TELEMETRYFLOW_API_KEY_SECRET` | TelemetryFlow API Key Secret | - |
| `TELEMETRYFLOW_ENDPOINT` | OTLP endpoint | `localhost:4317` |
| `TELEMETRYFLOW_SERVICE_NAME` | Service name | `Order-Service` |
| `TELEMETRYFLOW_SERVICE_VERSION` | Service version | `1.1.1` |

### Docker Compose Configuration

| Variable | Description | Default |
|----------|-------------|---------|
| `POSTGRES_VERSION` | PostgreSQL image version | `16-alpine` |
| `OTEL_VERSION` | OTEL Collector image version | `latest` |
| `CONTAINER_POSTGRES` | PostgreSQL container name | `order_service_postgres` |
| `CONTAINER_API` | API container name | `order_service_api` |
| `CONTAINER_OTEL` | OTEL container name | `order_service_otel` |
| `PORT_OTEL_GRPC` | OTEL gRPC port | `4317` |
| `PORT_OTEL_HTTP` | OTEL HTTP port | `4318` |
| `PORT_OTEL_METRICS` | OTEL metrics port | `8889` |

## Testing

```bash
# Run all tests
make test

# Run unit tests only
make test-unit

# Run integration tests
make test-integration

# Generate coverage report
make test-coverage
```

## Docker

```bash
# Build image
make docker-build

# Run container
make docker-run

# Start all services (app + database + monitoring)
make docker-compose-up

# Stop all services
make docker-compose-down
```

## Observability

The service is instrumented with OpenTelemetry for:
- **Traces**: Distributed tracing for request flows
- **Metrics**: Application and runtime metrics
- **Logs**: Structured logging

The OpenTelemetry Collector receives telemetry data and can export to:
- Jaeger (tracing)
- Prometheus (metrics)
- Any OTLP-compatible backend

### Prometheus Metrics

Access metrics at: `http://localhost:8889/metrics`

## License

Copyright (c) 2024-2026 TelemetryFlow. All rights reserved.
