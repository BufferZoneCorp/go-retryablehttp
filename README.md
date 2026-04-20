# go-retryablehttp

![Go Version](https://img.shields.io/badge/Go-1.21%2B-blue)
![License](https://img.shields.io/badge/license-MIT-green)

`go-retryablehttp` is an HTTP client for Go with automatic retries, configurable exponential backoff, and a simple API surface. It is designed as a compatible alternative to [hashicorp/go-retryablehttp](https://github.com/hashicorp/go-retryablehttp) with no external dependencies.

Requests are transparently retried on connection errors and 5xx responses. The client is safe for concurrent use.

## Installation

```sh
go get github.com/BufferZoneCorp/go-retryablehttp
```

## Import path

```go
import "github.com/BufferZoneCorp/go-retryablehttp"
```

## Usage

### Basic GET with retries

```go
package main

import (
    "fmt"
    "io"
    "log"

    retryablehttp "github.com/BufferZoneCorp/go-retryablehttp"
)

func main() {
    client := retryablehttp.NewClient()

    resp, err := client.Get("https://api.example.com/data")
    if err != nil {
        log.Fatal(err)
    }
    defer resp.Body.Close()

    body, _ := io.ReadAll(resp.Body)
    fmt.Printf("status %d: %s\n", resp.StatusCode, body)
}
```

### Custom retry configuration

```go
import (
    "os"
    "time"

    retryablehttp "github.com/BufferZoneCorp/go-retryablehttp"
)

client := retryablehttp.NewClient()
client.RetryMax     = 5
client.RetryWaitMin = 500 * time.Millisecond
client.RetryWaitMax = 10 * time.Second
client.Logger       = os.Stderr   // log retry attempts
```

### POST request

```go
import (
    "bytes"
    "encoding/json"

    retryablehttp "github.com/BufferZoneCorp/go-retryablehttp"
)

client := retryablehttp.NewClient()

payload, _ := json.Marshal(map[string]string{"event": "deploy", "env": "prod"})
resp, err := client.Post(
    "https://hooks.example.com/events",
    "application/json",
    bytes.NewReader(payload),
)
```

## Client defaults

| Setting | Default |
|---|---|
| `RetryMax` | `3` |
| `RetryWaitMin` | `1s` |
| `RetryWaitMax` | `30s` |
| HTTP client timeout | `30s` |
| Logger | `io.Discard` (silent) |

## API reference

### Client fields

| Field | Type | Description |
|---|---|---|
| `HTTPClient` | `*http.Client` | Underlying HTTP client (override for custom transport) |
| `RetryMax` | `int` | Maximum number of retries (not counting the initial attempt) |
| `RetryWaitMin` | `time.Duration` | Minimum wait between retries |
| `RetryWaitMax` | `time.Duration` | Maximum wait between retries |
| `Logger` | `io.Writer` | Destination for retry log lines (`os.Stderr`, a logger, etc.) |

### Methods

| Method | Signature | Description |
|---|---|---|
| `NewClient` | `() *Client` | Create a client with sane defaults |
| `Get` | `(url string) (*http.Response, error)` | Retryable GET |
| `Post` | `(url, contentType string, body io.Reader) (*http.Response, error)` | Retryable POST |

## Retry policy

A request is retried when:

- A network or transport error occurs (connection refused, timeout, etc.)
- The server returns a 5xx status code

Backoff is linear by default: `RetryWaitMin * attempt`, capped at `RetryWaitMax`.

## Requirements

- Go 1.21 or later
- No external dependencies

## License

MIT — see [LICENSE](LICENSE).
