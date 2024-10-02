# DateTime Client

This project implements a client for fetching date and time information from any server exposing a datetime API. The client supports both JSON and plain text responses.

## Table of Contents

- [Features](#features)
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Usage](#usage)
- [API](#api)
- [Error Handling](#error-handling)
- [Docker Support](#docker-support)
- [Makefile Commands](#makefile-commands)

## Features

- Fetch date and time from any server exposing a datetime API
- Support for JSON and plain text response formats
- Command-line flag configuration
- Exponential backoff retry mechanism
- Logging with Zap logger
- Docker support for easy deployment
- Comprehensive error handling

## Prerequisites

- Go 1.16 or higher
- Docker (optional, for containerization)
- Make (optional, for using Makefile commands)

## Installation

1. Clone the repository:
   ```
   git clone https://github.com/codescalersinternships/Datetime-Client-MohamedFadel.git
   cd Datetime-Client-MohamedFadel
   ```

2. Install dependencies:
   ```
   go mod download
   ```

## Usage

To use the DateTime Client, run the main program with the following flags:

```
go run main.go -url=http://localhost:8000 -type=application/json
```

Flags:
- `-url`: The URL of the datetime server (default: "http://localhost:8080")
- `-type`: The desired content type for the response ("application/json" or "text/plain", default: "application/json")

## API

### `NewDateTimeClient(serverURL string) *DateTimeClient`

Creates a new DateTimeClient with the specified server URL.

### `(c *DateTimeClient) GetDateTime(contentType string) (string, error)`

Retrieves the current date and time from the server.

- `contentType`: The desired content type for the response ("application/json" or "text/plain")

Returns the retrieved datetime as a string and an error if the request fails.

## Error Handling

The client defines several custom errors in `errors.go`:

- `ErrUnsupportedContentType`: Returned when an unsupported content type is specified
- `ErrCreatingRequest`: Returned when there's an error creating the HTTP request
- `ErrMakingRequest`: Returned when there's an error making the HTTP request
- `ErrReadingResponseBody`: Returned when there's an error reading the response body
- `ErrUnmarshallingJSON`: Returned when there's an error unmarshalling JSON response

## Docker Support

The project includes a Dockerfile for containerization. To build and run the Docker image:

```bash
docker build -t datetime-client .
docker run datetime-client -url=http://host.docker.internal:8000 -type=application/json
```

Note: Use `host.docker.internal` to access the host machine's localhost from within the Docker container.

## Makefile Commands

The project includes a Makefile with the following commands:

- `make all`: Run all tasks (deps, fmt, lint, build, test, docker-build, docker-run)
- `make build`: Build the client binary
- `make test`: Run tests for the client package
- `make clean`: Clean up build artifacts
- `make run`: Build and run the client
- `make deps`: Download dependencies
- `make fmt`: Format the Go code
- `make lint`: Run golangci-lint
- `make docker-build`: Build the Docker image
- `make docker-run`: Run the client in a Docker container

To use these commands, run `make <command>` in the project root directory.