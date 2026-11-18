# Yandex Cloud Traefik HTTP Provider

A dynamic HTTP provider for Traefik that integrates with Yandex Cloud services to automatically configure routing and SSL certificates.

## Features

- **Dynamic HTTP Provider for Traefik** - Serves Traefik configuration via HTTP endpoint
- **Yandex Cloud Certificate Manager Integration** - Automatically fetches and updates SSL certificates from Yandex Cloud Certificate Manager
- **Automatic Certificate Rotation** - Periodically refreshes certificates (configurable interval, default: 5 minutes)
- **Flexible Authentication** - Supports service account key file or instance metadata authentication
- **Built with Go** - Uses uber-go/fx dependency injection framework for clean architecture
- **Easy Deployment** - Available as Docker container or standalone binary

## Installation

### Using Docker

```bash
docker pull ghcr.io/nikolaymatrosov/yc-traefik-http-provider:latest
```

### From Source

```bash
go install github.com/nikolaymatrosov/yc-traefik-http-provider/cmd/yc-provider@latest
```

### Download Binary

Download the latest release from the [releases page](https://github.com/nikolaymatrosov/yc-traefik-http-provider/releases).

## Configuration

Create a `config.yaml` file:

```yaml
port: 8080
logLevel: info
saKeyPath: ./key.json  # Optional: Path to Yandex Cloud service account key
certFetchIntervalSeconds: 300  # How often to refresh certificates (default: 300)

features:
  certificateManager:
    enabled: true
    certificateIds:  # Yandex Cloud Certificate Manager IDs
      - "your-certificate-id"
      - "another-certificate-id"
```

### Configuration Options

| Option | Environment Variable | Default | Description |
|--------|---------------------|---------|-------------|
| `port` | `PORT` | `8080` | HTTP server port |
| `logLevel` | `LOG_LEVEL` | `info` | Logging level |
| `saKeyPath` | `SA_KEY_PATH` | - | Path to Yandex Cloud service account key file |
| `certFetchIntervalSeconds` | `CERT_FETCH_INTERVAL_SECONDS` | `300` | Certificate refresh interval in seconds |
| `features.certificateManager.enabled` | - | `false` | Enable Certificate Manager integration |
| `features.certificateManager.certificateIds` | `CERTIFICATE_IDS` (comma-separated) | - | Yandex Cloud Certificate IDs |
| `features.certificateManager.certificateId` | `CERTIFICATE_ID` | - | Single Certificate ID (legacy, merged into the list) |

### Environment Variables

Environment variables take precedence over configuration file values. Setting `CERTIFICATE_ID` or `CERTIFICATE_IDS` enables the Certificate Manager integration. You can also use Yandex Cloud instance metadata for authentication when running on Yandex Cloud Compute instances.
## Usage

### Running with Docker

```bash
docker run -p 8080:8080 \
  -v $(pwd)/config.yaml:/app/config.yaml \
  -v $(pwd)/key.json:/app/key.json \
  ghcr.io/nikolaymatrosov/yc-traefik-http-provider:latest
```

Or using environment variables:

```bash
docker run -p 8080:8080 \
  -e PORT=8080 \
  -e CERTIFICATE_IDS=your-cert-id,another-cert-id \
  -e SA_KEY_PATH=/app/key.json \
  -e LOG_LEVEL=debug \
  -v $(pwd)/key.json:/app/key.json \
  ghcr.io/nikolaymatrosov/yc-traefik-http-provider:latest
```

### Running the Binary

```bash
yc-provider
```

### Integrating with Traefik

Configure Traefik to use this provider by adding an HTTP provider in your Traefik configuration:

```yaml
providers:
  http:
    endpoint: "http://localhost:8080/config"
    pollInterval: "5s"
```

The provider returns dynamic configuration in Traefik format, including:
- TLS certificates fetched from Yandex Cloud Certificate Manager
- Certificate chains and private keys automatically formatted for Traefik

Example response:
```json
{
  "tls": {
    "certificates": [
      {
        "certFile": "-----BEGIN CERTIFICATE-----\n...",
        "keyFile": "-----BEGIN PRIVATE KEY-----\n..."
      }
    ]
  }
}
```

All configured certificates are returned. If a refresh fails, the last successfully fetched content is kept. Until every configured certificate has been fetched at least once, `/config` responds with `503` — Traefik then keeps its previous configuration instead of dropping TLS.

## Development

### Prerequisites

- Go 1.27 or later
- Docker (optional)

### Building

```bash
go build -o yc-provider ./cmd/yc-provider
```

### Running Locally

```bash
go run ./cmd/yc-provider/provider.go
```

## Project Structure

- `cmd/` - Application entry points
  - `cli/` - CLI commands
  - `yc-provider/` - Main HTTP provider service
- `pkg/` - Library code
  - `cert/` - Certificate management
  - `config/` - Configuration handling
  - `logger/` - Logging setup
  - `remote/` - Remote data sources
  - `server/` - HTTP server implementation
- `tf/` - Terraform infrastructure as code

## Release Process

Releases are automatically built and published to GitHub Container Registry when a new GitHub release is created. The release workflow uses GoReleaser to:

- Build binaries for linux/amd64
- Create Docker images
- Generate release artifacts and checksums

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

