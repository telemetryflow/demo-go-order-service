# Data Flow Diagram (DFD)

## order-service

This document describes the data flow for order-api.

## Level 0 - Context Diagram

```mermaid
graph LR
    subgraph External
        Client[Client Application]
        TF[TelemetryFlow Platform]
    end

    subgraph System
        API[order-service API]
        DB[(Database)]
    end

    Client -->|HTTP Requests| API
    API -->|HTTP Responses| Client
    API -->|CRUD Operations| DB
    DB -->|Query Results| API
    API -->|Telemetry Data| TF
```

## Level 1 - System Diagram

```mermaid
graph TB
    subgraph Client Layer
        WEB[Web Client]
        MOBILE[Mobile App]
        CLI[CLI Tool]
    end

    subgraph API Gateway
        LB[Load Balancer]
        RATE[Rate Limiter]
        AUTH[Auth Middleware]
    end

    subgraph Application Layer
        subgraph Presentation
            HANDLER[HTTP Handlers]
        end

        subgraph Application
            CMD[Command Handlers]
            QRY[Query Handlers]
        end

        subgraph Domain
            ENT[Entities]
            REPO[Repository Interfaces]
        end

        subgraph Infrastructure
            PERSIST[Persistence Layer]
            CACHE[Cache Layer]
        end
    end

    subgraph Data Layer
        DB[(PostgreSQL)]
        REDIS[(Redis Cache)]
    end
    subgraph Observability
        TF[TelemetryFlow]
    end

    WEB --> LB
    MOBILE --> LB
    CLI --> LB
    LB --> RATE
    RATE --> AUTH
    AUTH --> HANDLER

    HANDLER --> CMD
    HANDLER --> QRY
    CMD --> REPO
    QRY --> REPO
    REPO --> PERSIST
    PERSIST --> DB
    PERSIST -.-> CACHE
    CACHE -.-> REDIS
    HANDLER --> TF
    CMD --> TF
    PERSIST --> TF
```

## Level 2 - CQRS Flow

### Command Flow (Write Operations)

```mermaid
sequenceDiagram
    participant C as Client
    participant H as HTTP Handler
    participant V as Validator
    participant CH as Command Handler
    participant R as Repository
    participant DB as Database
    participant T as Telemetry

    C->>H: POST /api/v1/resource
    H->>V: Validate Request
    V-->>H: Validation Result

    alt Validation Failed
        H-->>C: 400 Bad Request
    else Validation Passed
        H->>CH: Execute Command
        CH->>T: Start Span
        CH->>R: Create/Update Entity
        R->>DB: INSERT/UPDATE
        DB-->>R: Result
        R-->>CH: Entity
        CH->>T: End Span
        CH->>T: Log & Metrics
        CH-->>H: Command Result
        H-->>C: 201 Created / 200 OK
    end
```

### Query Flow (Read Operations)

```mermaid
sequenceDiagram
    participant C as Client
    participant H as HTTP Handler
    participant QH as Query Handler
    participant R as Repository
    participant CACHE as Cache
    participant DB as Database
    participant T as Telemetry

    C->>H: GET /api/v1/resource
    H->>QH: Execute Query
    QH->>T: Start Span
    QH->>R: Find Entity
    R->>CACHE: Check Cache

    alt Cache Hit
        CACHE-->>R: Cached Data
    else Cache Miss
        R->>DB: SELECT
        DB-->>R: Query Result
        R->>CACHE: Store in Cache
    end

    R-->>QH: Entity/List
    QH->>T: End Span
    QH-->>H: Query Result
    H-->>C: 200 OK
```

## Data Transformations

| Layer | Input | Output | Transformation |
|-------|-------|--------|----------------|
| HTTP Handler | HTTP Request | DTO | Parse JSON, Validate |
| Command Handler | DTO | Domain Event | Apply Business Rules |
| Repository | Entity | DB Row | Serialize to DB Format |
| Query Handler | Query Params | DTO | Projection, Filtering |

## Error Handling Flow

```mermaid
graph TD
    REQ[Request] --> VALID{Validation}
    VALID -->|Invalid| E400[400 Bad Request]
    VALID -->|Valid| AUTH{Authorization}
    AUTH -->|Unauthorized| E401[401 Unauthorized]
    AUTH -->|Forbidden| E403[403 Forbidden]
    AUTH -->|OK| BIZ{Business Logic}
    BIZ -->|Not Found| E404[404 Not Found]
    BIZ -->|Conflict| E409[409 Conflict]
    BIZ -->|Error| E500[500 Internal Error]
    BIZ -->|Success| OK[200/201 Success]
```

## Notes

1. All requests go through rate limiting
2. Commands modify state, Queries are read-only (CQRS)
3. Cache is checked before database on read operations
4. Telemetry captures metrics, logs, and traces at each layer
