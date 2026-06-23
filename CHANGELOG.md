# Order-Service

<div align="center">
  <picture>
    <source media="(prefers-color-scheme: dark)" srcset="https://github.com/telemetryflow/.github/raw/main/docs/assets/tfo-logo-sdk-dark.svg">
    <source media="(prefers-color-scheme: light)" srcset="https://github.com/telemetryflow/.github/raw/main/docs/assets/tfo-logo-sdk-light.svg">
    <img src="https://github.com/telemetryflow/.github/raw/main/docs/assets/tfo-logo-sdk-light.svg" alt="TelemetryFlow Logo" width="80%">
  </picture>

[![Version](https://img.shields.io/badge/Version-1.2.0-orange.svg)](CHANGELOG.md)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Go Version](https://img.shields.io/badge/Go-1.26+-00ADD8?logo=go)](https://golang.org/)
[![OpenTelemetry](https://img.shields.io/badge/OTLP-100%25%20Compliant-success?logo=opentelemetry)](https://opentelemetry.io/)
[![Docker](https://img.shields.io/badge/Docker-Ready-2496ED?style=flat&logo=docker)](https://hub.docker.com/r/telemetryflow/telemetryflow-sdk)

</div>

<p align="center">
<strong>[GENERATED TelemetryFlow SDK]</strong> Order-Service - RESTful API with DDD + CQRS Pattern
</p>

---

# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.2.0] - 2026-06-23

### Added

- **Dependabot Configuration**: Added `.github/dependabot.yml` for automated dependency updates (gomod, docker, github-actions ecosystems)
- **TFO Platform Profile**: Added `platform` profile to docker-compose with TFO-Backend, TFO-Viz, and infrastructure services (PostgreSQL, ClickHouse, Redis, NATS) for end-to-end local observability

### Changed

- **Go Version**: Upgraded from Go 1.24 to Go 1.26
  - `go.mod`: `go 1.24.0` → `go 1.26.0`, `toolchain go1.24.11` → `toolchain go1.26.3`
  - `Dockerfile`: `golang:1.24-alpine` → `golang:1.26-alpine`
- **TelemetryFlow Go SDK**: Upgraded from v1.1.2 to v1.2.0
- **OpenTelemetry SDK**: Upgraded from v1.39.0 to v1.43.0
- **TFO-Collector**: Updated from v1.1.2 to v1.2.1 in docker-compose
- **GitHub Actions CI/CD**: Upgraded all workflow action versions
  - `actions/checkout` v4 → v6
  - `actions/setup-go` v5 → v6
  - `actions/upload-artifact` v4 → v7
  - `actions/download-artifact` v4 → v8
  - `docker/metadata-action` v5 → v6
  - `docker/setup-qemu-action` v3 → v4
  - `docker/setup-buildx-action` v3 → v4
  - `docker/login-action` v3 → v4
  - `docker/build-push-action` v6 → v7
  - `softprops/action-gh-release` v2 → v3
  - `golangci/golangci-lint-action` v7 → v9
- **Docker Compose**: Removed Grafana and Jaeger services; telemetry visualization now handled by TFO-Viz (via `platform` profile)
- **golangci-lint**: Added `-ST1000` to staticcheck exclusions (package comments present but not recognized through license header blocks)
- **File Headers**: Refactored all 65 Go source files to standard Apache 2.0 license header format
- **Version**: Updated from 1.1.2 to 1.2.0

### Removed

- **Grafana**: Removed from docker-compose (replaced by TFO-Viz via `platform` profile)
- **Jaeger**: Removed from docker-compose (replaced by TFO-Backend distributed tracing via `platform` profile)
