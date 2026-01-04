# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Bypass Go proxy for TelemetryFlow private modules
ENV GOPRIVATE=github.com/telemetryflow/*

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Ensure all dependencies are resolved
RUN go mod tidy

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/Order-Service ./cmd/api

# Runtime stage
FROM alpine:3.21

WORKDIR /app

# Update packages to get security patches (CVE fixes) and install runtime dependencies
RUN apk upgrade --no-cache && \
    apk add --no-cache ca-certificates tzdata

# Copy binary from builder
COPY --from=builder /app/Order-Service .
COPY --from=builder /app/configs ./configs
COPY --from=builder /app/migrations ./migrations

# Copy .env file if exists (optional - can be overridden by environment variables)
COPY .env* ./

# Create non-root user
RUN addgroup -g 1000 appgroup && \
    adduser -u 1000 -G appgroup -s /bin/sh -D appuser && \
    chown -R appuser:appgroup /app

USER appuser

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Run the application
ENTRYPOINT ["./Order-Service"]
