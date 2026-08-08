# keepalivectl

[![Go Version](https://img.shields.io/github/go-mod/go-version/fishfinal/keepalivectl)](https://golang.org/)
[![Go Report Card](https://goreportcard.com/badge/github.com/fishfinal/keepalivectl)](https://goreportcard.com/report/github.com/fishfinal/keepalivectl)
[![License](https://img.shields.io/github/license/fishfinal/keepalivectl)](LICENSE)
[![Release](https://img.shields.io/github/v/release/fishfinal/keepalivectl)](https://github.com/fishfinal/keepalivectl/releases)

**keepalivectl** is a command-line tool for testing gRPC server Keepalive and connection pooling support.

## Features

- Verify gRPC server Keepalive support
- Test connection pooling stability
- Concurrent connection testing
- Real-time connection state monitoring
- Cross-platform support (Linux, macOS, Windows)
- Colorful terminal output with aurora
- Structured logging with gologger
- Detailed statistics and summary reports

## Quick Start

### Installation

#### Using Go

```bash
go install github.com/fishfinal/keepalivectl/cmd/keepalivectl@latest
```

#### From Releases

```bash
# Linux
wget https://github.com/fishfinal/keepalivectl/releases/download/v1.0.0/keepalivectl_linux_amd64
chmod +x keepalivectl_linux_amd64
sudo mv keepalivectl_linux_amd64 /usr/local/bin/keepalivectl

# macOS
wget https://github.com/fishfinal/keepalivectl/releases/download/v1.0.0/keepalivectl_darwin_amd64
chmod +x keepalivectl_darwin_amd64
sudo mv keepalivectl_darwin_amd64 /usr/local/bin/keepalivectl

# Windows (download and rename to keepalivectl.exe)
```

### Usage

```bash
# Basic test - connect to endpoint, run for 2 minutes
keepalivectl --endpoint localhost:9600 --duration 2m

# Concurrent test - 10 connections, 3 minutes
keepalivectl --concurrency 10 --duration 3m

# Custom keepalive parameters
keepalivectl --keepalive-time 30s --keepalive-timeout 10s --duration 5m

# Quick smoke test
keepalivectl --duration 30s --check-interval 1s

# High concurrency test
keepalivectl --concurrency 50 --duration 5m --check-interval 5s
```

## Command Line Options

### Target Options

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--endpoint` | `-e` | gRPC server endpoint address | `localhost:9600` |

### Test Behavior Options

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--duration` | `-d` | Test duration | `2m0s` |
| `--check-interval` | `-i` | Connection state check interval | `3s` |
| `--concurrency` | `-c` | Number of concurrent connections | `1` |

### Keepalive Settings

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--keepalive-time` | `-t` | Keepalive ping interval | `10s` |
| `--keepalive-timeout` | `-T` | Keepalive ping timeout | `5s` |

### Global Options

| Flag | Description |
|------|-------------|
| `-h, --help` | Show help message |
| `--version` | Show version information |

## Example Output

```bash
$ keepalivectl --endpoint localhost:9600 --duration 2m --concurrency 5
```

```
[INF] 2026-09-03T15:17:32+08:00 [Connection 0] connected successfully, monitoring (keepalive: 10s, timeout: 5s)
[INF] 2026-09-03T15:17:35+08:00 [Connection 0] state changed: READY (change #1)
[INF] 2026-09-03T15:17:38+08:00 [Connection 0] state: READY
[INF] 2026-09-03T15:17:41+08:00 [Connection 0] state: READY
[INF] 2026-09-03T15:17:44+08:00 [Connection 0] state: READY
...
[INF] 2026-09-03T15:19:32+08:00 Test completed, duration: 2m0s

============================================================
📊 Test Summary
============================================================
Endpoint:           localhost:9600
Keepalive Interval: 10s
Keepalive Timeout:  5s
Test Duration:      2m0s
Concurrency:        5
============================================================
Total Checks:       200
Ready State Count:  198
State Changes:      5
✅ All connections healthy, server Keepalive support is good
============================================================
```

## Use Cases

### 1. Verify Server Keepalive Support

Test whether your gRPC server properly handles keepalive pings:

```bash
keepalivectl --endpoint grpc-server:9600 --duration 5m --keepalive-time 30s
```

### 2. Load Testing

Simulate real-world workload with multiple concurrent connections:

```bash
keepalivectl --concurrency 100 --duration 10m --check-interval 5s
```

### 3. Debug Connection Issues

Quickly check if connection stability issues are related to keepalive:

```bash
keepalivectl --duration 30s --check-interval 1s
```

### 4. CI/CD Integration

Use in automated testing pipelines:

```yaml
# GitHub Actions example
- name: Test gRPC Keepalive
  run: |
    ./keepalivectl --endpoint ${{ secrets.GRPC_ENDPOINT }} --duration 1m
```

## Troubleshooting

### Error: `too_many_pings`

```
ERROR: Client received GoAway with error code ENHANCE_YOUR_CALM 
and debug data equal to ASCII "too_many_pings".
```

**Solution**: The server has a minimum keepalive interval (`MinTime`). Increase the client's keepalive interval:

```bash
keepalivectl --keepalive-time 30s
```

### Error: Connection Refused

```
connection failed: dial tcp 127.0.0.1:9600: connect: connection refused
```

**Solution**: Ensure the gRPC server is running on the specified endpoint:

```bash
# Check if server is running
netstat -an | grep 9600
# or
lsof -i :9600
```

## Short Option Reference

For quick command-line usage, short options are also available:

| Long Option | Short Option | Description |
|-------------|--------------|-------------|
| `--endpoint` | `-e` | gRPC server endpoint address |
| `--duration` | `-d` | Test duration |
| `--check-interval` | `-i` | Connection state check interval |
| `--concurrency` | `-c` | Number of concurrent connections |
| `--keepalive-time` | `-t` | Keepalive ping interval |
| `--keepalive-timeout` | `-T` | Keepalive ping timeout |

```bash
# Example using short options (equivalent to long option examples above)
keepalivectl -e localhost:9600 -d 2m -c 10
keepalivectl -t 30s -T 10s -d 5m
```

## Development

### Prerequisites

- Go 1.21 or higher
- Make (optional)

### Building from Source

```bash
# Clone the repository
git clone https://github.com/fishfinal/keepalivectl.git
cd keepalivectl

# Build
make build

# Or build directly
go build -o keepalivectl ./cmd/keepalivectl

# Build for all platforms
make build-all

# Run tests
make test

# Run linting
make lint

# Format code
make fmt
```

### Project Structure

```
keepalivectl/
├── cmd/
│   └── keepalivectl/
│       ├── main.go          # Entry point
│       └── helper/
│           └── flagshelper/ # Flag grouping utilities
├── internal/
│   ├── config/              # Configuration management
│   ├── monitor/             # Connection monitoring
│   └── reporter/            # Result reporting
├── go.mod
├── go.sum
├── Makefile
├── README.md
└── LICENSE
```

### Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

Please ensure your code passes all tests and linting checks:

```bash
make test
make lint
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- [grpcurl](https://github.com/fullstorydev/grpcurl) - Inspiration for gRPC testing
- [aurora](https://github.com/logrusorgru/aurora) - Terminal color support
- [gologger](https://github.com/shaichunfeng/gologger) - Structured logging
- [cobra](https://github.com/spf13/cobra) - CLI framework

Made with ❤️ by [fishfinal](https://github.com/fishfinal)

