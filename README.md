# XIS Data Aggregator

A high-performance data aggregation service built in Go that processes, aggregates, and serves data through both REST and gRPC APIs. The service includes real-time data processing, Redis-based storage, and comprehensive monitoring capabilities.

## 🚀 Features

- **Dual API Support**: REST API (Gin) and gRPC API for flexible client integration
- **Real-time Data Processing**: Concurrent worker-based data processing pipeline
- **Redis Storage**: Fast, in-memory data storage with persistence
- **Swagger Documentation**: Auto-generated API documentation
- **Metrics Collection**: Built-in metrics and monitoring
- **Docker Support**: Containerized deployment
- **Mock Data Generation**: Simulated data input for testing and development
- **Configurable Architecture**: Tunable worker counts, batch sizes, and intervals

## 📋 Prerequisites

- Go 1.23.0 or higher
- Redis server
- Docker (optional, for containerized deployment)
- Protocol Buffers compiler (for gRPC development)

## 🛠️ Installation

### Local Development

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd xis-data-aggregator
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Generate Protocol Buffer code** (if working with gRPC)
   ```bash
   make proto
   ```

4. **Generate Swagger documentation**
   ```bash
   make swag
   ```

5. **Start Redis server**
   ```bash
   # Using Docker
   docker run -d -p 6379:6379 redis:alpine
   
   # Or install Redis locally
   # Follow Redis installation guide for your OS
   ```

### Docker Deployment

1. **Build the Docker image**
   ```bash
   make docker-build
   ```

2. **Run the container**
   ```bash
   make docker-run
   ```

## 🚀 Usage

### Running the Service

```bash
# Run with default configuration
go run cmd/xis-data-aggregator/main.go

# Run with custom parameters
go run cmd/xis-data-aggregator/main.go -workersCount=10 -r=9090 -g=50052
```

### Command Line Flags

| Flag | Description | Default |
|------|-------------|---------|
| `-workersCount` | Number of worker goroutines | 5 |
| `-b` | Metrics batch size | 10 |
| `-r` | REST API port | 8080 |
| `-g` | gRPC port | 50051 |
| `-n` | Input interval (ms) | 555 |
| `-l` | Input pack length | 10 |

## 📡 API Documentation

### REST API

The service exposes a REST API on port 8080 (default) with the following endpoints:

#### Get Data by ID
```http
GET /api/v1/data/{id}
```

**Parameters:**
- `id` (path): UUID of the data record

**Response:**
```json
{
  "id": "uuid-string",
  "ts": 1640995200,
  "max": 42
}
```

#### List Data by Time Range
```http
GET /api/v1/data?from={timestamp}&to={timestamp}
```

**Parameters:**
- `from` (query): Start timestamp (Unix timestamp)
- `to` (query): End timestamp (Unix timestamp)

**Response:**
```json
[
  {
    "id": "uuid-string",
    "ts": 1640995200,
    "max": 42
  }
]
```

### gRPC API

The service also provides a gRPC API on port 50051 (default). See the generated protobuf files in `pb/` directory for detailed service definitions.

### Swagger Documentation

Interactive API documentation is available at:
```
http://localhost:8080/swagger/index.html
```

## 🏗️ Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Mock Input    │    │   Worker Pool   │    │   Redis Store   │
│   Generator     │───▶│   (5 workers)   │───▶│                 │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                                │
                                ▼
                       ┌─────────────────┐
                       │  Metrics        │
                       │  Collector      │
                       └─────────────────┘
                                │
                                ▼
                       ┌─────────────────┐    ┌─────────────────┐
                       │   REST API      │    │   gRPC API      │
                       │   (Gin)         │    │   (Port 50051)  │
                       │   (Port 8080)   │    └─────────────────┘
                       └─────────────────┘
```

## 📁 Project Structure

```
xis-data-aggregator/
├── cmd/xis-data-aggregator/    # Application entry point
├── config/                     # Configuration management
├── docs/                       # Swagger and gRPC documentation
├── examples/                   # Example clients
├── gen/proto/                  # Protocol Buffer definitions
├── internal/                   # Internal application code
│   ├── api/                   # API handlers (REST/gRPC)
│   ├── metrics/               # Metrics collection
│   ├── mocks/                 # Mock data generators
│   ├── models/                # Data models
│   ├── repository/            # Data access layer
│   └── service/               # Business logic
├── pb/                        # Generated protobuf code
├── pkg/utils/                 # Utility functions
├── Dockerfile                 # Docker configuration
├── Makefile                   # Build automation
└── README.md                  # This file
```

## 🔧 Development

### Building

```bash
# Build the application
go build -o bin/xis-data-aggregator ./cmd/xis-data-aggregator

# Run tests
go test ./...

# Run with race detection
go run -race cmd/xis-data-aggregator/main.go
```

### Code Generation

```bash
# Generate Protocol Buffer code
make proto

# Generate Swagger documentation
make swag

# Generate gRPC documentation
make grpcdoc
```

### Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific test
go test ./internal/api/...
```

## 📊 Monitoring

The service includes built-in metrics collection that tracks:
- Data processing throughput
- Worker performance
- Error rates
- Processing latency

Metrics are collected in batches and can be exposed through the metrics collector.

## 🐳 Docker

### Building the Image

```bash
make docker-build
```

### Running the Container

```bash
make docker-run
```

### Custom Docker Run

```bash
docker run --rm -it \
  --memory="128m" \
  --name=xis-data-aggregator \
  --volume=./logs:/tmp \
  -p 8080:8080 \
  -p 50051:50051 \
  xis-data-aggregator:latest
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🆘 Support

For support and questions:
- Create an issue in the repository
- Check the API documentation at `/swagger/index.html`
- Review the example client in `examples/`

## 🔄 Version History

- **v0.1.0**: Initial release with REST and gRPC APIs, Redis storage, and worker-based processing

---

**Note**: This service is designed for high-throughput data processing and aggregation. The default configuration is tuned for development environments. For production deployment, consider adjusting worker counts, batch sizes, and intervals based on your specific requirements and infrastructure capabilities. 


 