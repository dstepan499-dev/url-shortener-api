# High-Performance URL Shortener & Analytics API

A scalable and production-ready RESTful URL Shortener service built with Go (Golang) and PostgreSQL.

## Features

- URL Shortening: Generates cryptographically secure, randomized 6-character aliases or accepts custom aliases.
- Fast Redirects: HTTP 302 redirection with optimized database index lookups.
- Click Analytics: Tracks total visits and metadata (creation date, original target URL).
- Atomic Increments: Database counters are safely incremented in a single SQL operation to prevent race conditions.
- Clean Architecture: Distinct separation of concerns (Handler -> Service -> Repository).
- Automated Testing: Includes unit tests for alias generation and input validation.
- Containerization: Fully orchestrated using multi-stage Docker and Docker Compose.

## Tech Stack

- Language: Go (Golang) 1.23+
- Router: go-chi/chi/v5
- Database: PostgreSQL 16
- DB Driver: jackc/pgx/v5
- DevOps: Docker, Docker Compose

## API Reference

### 1. Shorten URL
POST /api/shorten
Content-Type: application/json

Payload:
{
  "url": "https://example.com",
  "custom_alias": "custom"
}

Response (201 Created):
{
  "alias": "custom",
  "short_url": "http://localhost:8080/custom"
}

### 2. Redirect
GET /{alias}
Response: HTTP 302 Found (Redirects to original destination)

### 3. Get Analytics
GET /api/analytics/{alias}
Response (200 OK):
{
  "id": 1,
  "alias": "custom",
  "original_url": "https://example.com",
  "clicks": 14,
  "created_at": "2026-09-18T12:00:00Z"
}

## Quick Start

1. Clone the repository:
git clone https://github.com/dstepan499-dev/url-shortener-api.git
cd url-shortener-api

2. Run with Docker Compose (Recommended):
docker compose up -d --build

3. Run Tests:
go test -v ./...