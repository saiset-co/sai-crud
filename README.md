# SAI CRUD Service

A microservice for CRUD operations built with Go and the SAI framework. This service provides a RESTful API for creating, reading, updating, and deleting documents in collections with configurable storage backend.

## Features

- **Full CRUD Operations**: Create, Read, Update, Delete documents
- **Collection Prefixing**: Support for prefixed collections for multi-tenancy
- **Flexible Storage**: Integration with SAI Storage service
- **Configurable**: Environment-based configuration with templates
- **Docker Ready**: Containerized deployment with Docker Compose
- **Health Checks**: Built-in health monitoring
- **API Documentation**: Automatic OpenAPI documentation generation
- **Middleware Support**: CORS, logging, recovery, and authentication
- **Validation**: Request validation with comprehensive error handling

## Quick Start

### Prerequisites

- Go 1.21+ (for local development)
- Docker and Docker Compose (for containerized deployment)
- SAI Storage service running

### Local Development

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd sai-crud
   ```

2. **Set up environment**
   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

3. **Install dependencies**
   ```bash
   make deps
   ```

4. **Generate configuration**
   ```bash
   make config
   ```

5. **Run the service**
   ```bash
   make run
   ```

The service will start on `http://localhost:8081` (configurable via `SERVER_PORT`)

### Docker Deployment

1. **Start with Docker Compose**
   ```bash
   make up
   ```

2. **View logs**
   ```bash
   make logs
   ```

3. **Stop services**
   ```bash
   make down
   ```

## API Endpoints

### Base URL
```
http://localhost:8081/api/v1
```

### Create Documents
```http
POST /api/v1/
Content-Type: application/json

{
  "prefix": "user",
  "data": [
    {"name": "John", "age": 30},
    {"name": "Jane", "age": 25}
  ]
}
```

### Read Documents
```http
GET /api/v1/
Content-Type: application/json

{
  "prefix": "user",
  "filter": {"age": {"$gte": 25}},
  "sort": {"name": 1},
  "limit": 10,
  "skip": 0
}
```

### Update Documents
```http
PUT /api/v1/
Content-Type: application/json

{
  "prefix": "user",
  "filter": {"name": "John"},
  "data": {"$set": {"age": 31, "status": "active"}}
}
```

**Update operators:**
- `$set`: Set field values
- `$unset`: Remove fields
- `$inc`: Increment numeric values
- `$push`: Add to arrays

**Examples:**
```json
// Set fields
{"data": {"$set": {"age": 31, "email": "john@example.com"}}}

// Remove fields  
{"data": {"$unset": {"tempField": ""}}}

// Increment value
{"data": {"$inc": {"views": 1}}}
```

### Delete Documents
```http
DELETE /api/v1/
Content-Type: application/json

{
  "prefix": "user",
  "filter": {"age": {"$lt": 18}}
}
```

## Configuration

The service uses environment variables for configuration. Key settings include:

### Service Configuration
- `SERVICE_NAME`: Service name (default: `sai-service`)
- `SERVICE_VERSION`: Service version (default: `1.0.0`)
- `COLLECTION_NAME`: Default collection name (default: `dictionary`)

### Server Configuration
- `SERVER_HOST`: Server host (default: `127.0.0.1`)
- `SERVER_PORT`: Server port (default: `8081`)
- `SERVER_READ_TIMEOUT`: Read timeout in seconds (default: `30`)
- `SERVER_WRITE_TIMEOUT`: Write timeout in seconds (default: `30`)
- `SERVER_IDLE_TIMEOUT`: Idle timeout in seconds (default: `120`)

### Storage Configuration
- `STORAGE_URL`: SAI Storage service URL (default: `http://localhost:8080`)
- `STORAGE_USERNAME`: Storage service username
- `STORAGE_PASSWORD`: Storage service password

### Logging Configuration
- `LOG_LEVEL`: Log level (default: `debug`)
- `LOG_OUTPUT`: Log output (default: `stdout`)
- `LOG_FORMAT`: Log format (default: `console`)

## Development

### Available Make Commands

```bash
# Development
make deps          # Download Go dependencies
make config        # Generate config from template
make run           # Run the application locally
make build         # Build the application binary

# Docker
make docker-build  # Build Docker image
make docker-run    # Run Docker container
make up            # Start all services with docker-compose
make down          # Stop all services
make logs          # Show logs from all services
make restart       # Restart all services
make rebuild       # Rebuild and restart all services

# Code Quality
make fmt           # Format Go code
make vet           # Run go vet
make lint          # Run linter (requires golangci-lint)

# Cleanup
make clean         # Clean build artifacts
make clean-docker  # Clean Docker resources
make clean-all     # Clean everything
```

### Project Structure

```
.
├── cmd/                    # Application entry points
│   └── main.go            # Main application
├── internal/              # Internal packages
│   ├── handler.go         # HTTP handlers
│   └── service.go         # Business logic
├── types/                 # Type definitions
│   ├── request.go         # Request types
│   ├── response.go        # Response types
│   └── types.go           # Service types
├── scripts/               # Deployment scripts
├── config.template.yml   # Configuration template
├── .env                   # Environment variables
├── Dockerfile             # Docker configuration
├── docker-compose.yml     # Docker Compose configuration
└── Makefile              # Build automation
```

## Health Checks

The service includes built-in health checks:

```bash
curl http://localhost:8081/health
```

## API Documentation

When `DOCS_ENABLED=true`, interactive API documentation is available at:

```
http://localhost:8081/docs
```

## Authentication

The service supports two types of authentication:

### Incoming Request Authentication
Configure authentication for incoming API requests:

```env
USERNAME=your-username
PASSWORD=your-password
```

### Outgoing Client Authentication
Configure authentication for outgoing requests to other services (like SAI Storage). Each client can have its own authentication method and parameters configured in `config.yml`:

```yaml
clients:
  enabled: true
  services:
    storage:
      url: "http://localhost:8080"
      auth:
        provider: "basic"           # Authentication method
        payload:
          username: "storage-user"
          password: "storage-pass"
```

**Supported authentication providers:**
- `basic`: HTTP Basic Authentication
- `bearer`: Bearer token authentication
- `api-key`: API key authentication
- Custom providers as configured

**Environment variables for client authentication:**
```env
STORAGE_URL=http://localhost:8080
STORAGE_AUTH_PROVIDER=basic
STORAGE_USERNAME=storage-user
STORAGE_PASSWORD=storage-pass
```

## Error Handling

All API responses follow a consistent format:

**Success Response:**
```json
{
  "data": [...],
  "created": 2
}
```

**Error Response:**
```json
{
  "error": "validation failed",
  "details": "specific error details"
}
```

## Monitoring

The service includes:

- **Health Checks**: `/health` endpoint
- **Metrics**: Built-in metrics collection
- **Logging**: Structured logging with configurable levels
- **Request Tracing**: Request/response logging middleware

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run tests and linting
5. Submit a pull request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Support

For issues and questions, please create an issue in the repository or contact the development team.