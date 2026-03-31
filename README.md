# Order Packs Calculator

A REST API that calculates the optimal combination of pack sizes to fulfill orders, minimizing the number of packs while meeting or exceeding the requested quantity.

## Features

- CRUD operations for managing available pack sizes
- Dynamic programming algorithm for optimal pack calculation
- Web UI for managing packs and calculating orders
- SQLite database with automatic migration and seeding
- Dockerized deployment

## Getting Started

### Prerequisites

- Go 1.24+
- CGO enabled (for SQLite)
- Docker (optional)

## Live Demo

**https://order-packs-calculator.onrender.com**

> **Note:** Hosted on Render's free tier. The instance spins down after inactivity, so the first request may take 50 seconds or more to respond.

### Run Locally

```bash
go run ./src/cmd/api -port=8080
```

### Run with Docker

```bash
docker compose up --build
```

Clean up containers, images:

```bash
docker compose down --rmi all
```

The application will be available at `http://localhost:8080`.

### Run Tests

```bash
go test ./...
```

## API Endpoints

Base URL: `/v1/order_packs`

| Method   | Endpoint | Description                      |
|----------|----------|----------------------------------|
| `GET`    | `/`      | List all pack sizes              |
| `POST`   | `/`      | Add a new pack size              |
| `PUT`    | `/:id`   | Update a pack size               |
| `DELETE` | `/:id`   | Remove a pack size               |
| `POST`   | `/solve` | Calculate optimal packs for order|

### Request / Response Examples

**Add a pack size**
```bash
curl -X POST http://localhost:8080/v1/order_packs \
  -H "Content-Type: application/json" \
  -d '{"items": 250}'
```

**Calculate packs for an order**
```bash
curl -X POST http://localhost:8080/v1/order_packs/solve \
  -H "Content-Type: application/json" \
  -d '{"order": 501}'
```

Response:
```json
{
  "result": true,
  "data": {
    "250": 1,
    "500": 1
  }
}
```

## Algorithm

The calculator uses a dynamic programming approach (coin change variant) to find the minimum number of packs that meet or exceed the order quantity.

Given packs `[250, 500, 1000, 2000, 5000]` and order `12001`:

| Pack | Quantity |
|------|----------|
| 5000 | 2        |
| 2000 | 1        |
| 250  | 1        |

Total: 12,250 items (minimum possible >= 12,001)

## Project Structure

```
src/
├── cmd/api/
│   ├── main.go              # Entry point, server setup
│   ├── static/               # Frontend (index.html)
│   ├── router/               # HTTP handlers & routing
│   └── requests/             # Request DTOs
├── internal/
│   ├── model/                # Database entities
│   ├── service/              # Business logic & algorithm
│   └── repository/           # Database operations
└── shared/
    ├── db/                   # SQLite initialization & seeding
    └── logging/              # Logger utilities
```

## Default Seed Data

On first run, the following pack sizes are created automatically:

`250, 500, 1000, 2000, 5000, 10000`

## Tech Stack

- **Go** with Gin (HTTP) and GORM (ORM)
- **SQLite** for persistence
- **Docker** for containerization
