
# DateTime Client

This project implements a client for fetching date and time information from two different server types: a standard HTTP server and a Gin server. The client supports both JSON and plain text responses.

## Table of Contents

- [Features](#features)
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Usage](#usage)
- [Configuration](#configuration)
- [API](#api)
- [Error Handling](#error-handling)
- [Docker Support](#docker-support)
- [Makefile Commands](#makefile-commands)

## Features

- Fetch date and time from standard HTTP and Gin servers
- Support for JSON and plain text response formats
- Environment variable configuration
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

To use the DateTime Client, you need to set up the environment variables and then run the main program:

```go
client.SetupEnv("http://localhost", "8000", "9000")
dateTime, err := client.GetDateTime("standard", "application/json")
if err != nil {
    log.Fatal(err)
}
fmt.Println(dateTime)
```

## Configuration

The client requires the following environment variables to be set:

- `SERVER_URL`: The base URL of the datetime server (e.g., "http://localhost")
- `SERVER_PORT`: The port for the standard server
- `SERVER_PORT_GIN`: The port for the Gin server

You can set these variables using the `SetupEnv` function or by setting them directly in your environment.

## API

### `SetupEnv(serverURL, serverPort, serverPortGin string)`

Sets up the environment variables for the datetime client.

### `GetDateTime(serverType, contentType string) (string, error)`

Retrieves the current date and time from the specified server.

- `serverType`: The type of server to query ("standard" or "gin")
- `contentType`: The desired content type for the response ("application/json" or "text/plain")

Returns the retrieved datetime as a string and an error if the request fails.

## Error Handling

The client defines several custom errors in `errors.go`:

- `ErrURLandPortMustBeSet`: Returned when required environment variables are not set
- `ErrCreatingRequest`: Returned when there's an error creating the HTTP request
- `ErrMakingRequest`: Returned when there's an error making the HTTP request
- `ErrReadingResponseBody`: Returned when there's an error reading the response body
- `ErrUnmarshallingJSON`: Returned when there's an error unmarshalling JSON response
- `ErrUnsupportedContentType`: Returned when an unsupported content type is specified

## Docker Support

The project includes a Dockerfile for containerization. To build and run the Docker image:

```bash
docker build -t datetime-client .
docker run --network host -e SERVER_URL=http://localhost -e SERVER_PORT=8000 -e SERVER_PORT_GIN=9000 datetime-client
```

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
