# Entity Relationship Diagram (ERD)

## order-service

This document describes the database schema for order-api.

## Diagram

```mermaid
erDiagram
    %% Base entity fields are included in all entities
    %% id: UUID (PK)
    %% created_at: TIMESTAMP
    %% updated_at: TIMESTAMP
    %% deleted_at: TIMESTAMP (nullable, for soft delete)
    USERS {
        uuid id PK
        string email UK
        string password_hash
        string name
        string role
        boolean is_active
        timestamp created_at
        timestamp updated_at
        timestamp deleted_at
    }

    REFRESH_TOKENS {
        uuid id PK
        uuid user_id FK
        string token UK
        timestamp expires_at
        timestamp created_at
    }

    USERS ||--o{ REFRESH_TOKENS : has

    %% Add your entities here
    %% Example:
    %% PRODUCTS {
    %%     uuid id PK
    %%     string name
    %%     string description
    %%     decimal price
    %%     integer stock
    %%     timestamp created_at
    %%     timestamp updated_at
    %%     timestamp deleted_at
    %% }

```

## Tables

### Base Fields

All entities include the following base fields:

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | UUID | PRIMARY KEY | Unique identifier |
| created_at | TIMESTAMP | NOT NULL | Creation timestamp |
| updated_at | TIMESTAMP | NOT NULL | Last update timestamp |
| deleted_at | TIMESTAMP | NULLABLE | Soft delete timestamp |

### Users

User accounts for authentication.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | UUID | PRIMARY KEY | Unique identifier |
| email | VARCHAR(255) | UNIQUE, NOT NULL | User email |
| password_hash | VARCHAR(255) | NOT NULL | Bcrypt password hash |
| name | VARCHAR(255) | NOT NULL | User display name |
| role | VARCHAR(50) | NOT NULL | User role (admin, user) |
| is_active | BOOLEAN | DEFAULT true | Account active status |
| created_at | TIMESTAMP | NOT NULL | Creation timestamp |
| updated_at | TIMESTAMP | NOT NULL | Last update timestamp |
| deleted_at | TIMESTAMP | NULLABLE | Soft delete timestamp |

### Refresh Tokens

JWT refresh tokens for session management.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | UUID | PRIMARY KEY | Unique identifier |
| user_id | UUID | FOREIGN KEY | Reference to users |
| token | VARCHAR(500) | UNIQUE, NOT NULL | Refresh token value |
| expires_at | TIMESTAMP | NOT NULL | Token expiration |
| created_at | TIMESTAMP | NOT NULL | Creation timestamp |

## Indexes

```sql
-- Add indexes for frequently queried columns
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_role ON users(role);
CREATE INDEX idx_refresh_tokens_user_id ON refresh_tokens(user_id);
CREATE INDEX idx_refresh_tokens_expires_at ON refresh_tokens(expires_at);
```

## Notes

1. All UUIDs are generated using UUID v4
2. Soft delete is implemented via `deleted_at` column
3. Timestamps are stored in UTC
4. Indexes should be added based on query patterns
