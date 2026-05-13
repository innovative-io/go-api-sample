# go-api-sample
Sample Go API Project — a reference implementation of a Go REST API using the Go standard library (`net/http`), GORM, and Swagger.

[![Build Status](https://drone.onebytedata.net/api/badges/one-byte-data/go-api-sample/status.svg)](https://drone.onebytedata.net/one-byte-data/go-api-sample)

## Requirements

- Go 1.26+
- PostgreSQL 17+

## Usage

Run all tests `go test ./...`

Run benchmarks against the API `cd internal/controllers && go test -benchmem -run="^#" -bench . && cd ../..`

Run benchmarks against the Database `cd internal/services && go test -benchmem -run="^#" -bench . && cd ../..`

Install swagger spec generate tool `go install github.com/swaggo/swag/cmd/swag@latest`

Generate swagger spec `swag init --generalInfo cmd/server/main.go --output docs`

Build API `go build -v -a -o build/docker/go-api-sample cmd/server/main.go`

## Endpoints

| Method | Path | Description |
|--------|------|-------------|
| GET | /health | Server health check |
| GET | /cats | List all cats |
| GET | /cats/:id | Get cat by ID |
| GET | /cats/count | Total count of cats |
| POST | /cats | Create a cat |
| PUT | /cats/:id | Update a cat |
| DELETE | /cats/:id | Delete a cat |
| GET | /dogs | List all dogs |
| GET | /dogs/:id | Get dog by ID |
| GET | /dogs/count | Total count of dogs |
| POST | /dogs | Create a dog |
| PUT | /dogs/:id | Update a dog |
| DELETE | /dogs/:id | Delete a dog |
| GET | /swagger/* | Swagger UI |

## Runtime Environment Variables

`CONNECTION_STRING="postgresql://go_api:go_api@localhost:5432/animals?sslmode=disable"`

`PORT=8080` (optional, defaults to `8080`)

