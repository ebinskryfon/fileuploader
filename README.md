# FileUploader 📁

A secure, JWT-authenticated microservice for file upload and download with validation, metadata management, and rate limiting.

## Features ✨

- **Secure Authentication**: JWT-based authentication with user isolation
- **File Validation**: MIME type checking, size limits, and content verification
- **Metadata Management**: Automatic file metadata extraction and storage
- **Rate Limiting**: Configurable rate limiting (60 requests/minute by default)
- **RESTful API**: Clean REST endpoints with comprehensive error handling
- **Docker Support**: Ready-to-deploy Docker configuration
- **Health Checks**: Built-in health and readiness endpoints

## Quick Start 🚀

### Prerequisites

- Go 1.19+

### Installation

1. Clone the repository:
```bash
git clone https://github.com/yourusername/fileuploader.git
cd fileuploader
```

2. Install dependencies:
```bash
make dev-setup
```

3. Set environment variables:
```bash
export JWT_SECRET=your-secret-key
export PORT=8080
```

4. Run the service:
```bash
make run
```

The service will be available at `http://localhost:8080`

## Configuration ⚙️

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `JWT_SECRET` | Secret key for JWT token signing | Required |
| `PORT` | Server port | `8080` |
| `MAX_FILE_SIZE` | Maximum file size in bytes | `25MB` |
| `STORAGE_PATH` | File storage directory | `./storage` |

### File Constraints

- **Maximum Size**: 25MB (configurable)
- **Allowed Types**: JPEG, PNG, PDF (configurable)
- **Security**: MIME type validation, content sniffing, extension checking

## API Usage 📡

### Authentication

All endpoints require JWT authentication:

```bash
curl -H "Authorization: Bearer <your-jwt-token>" \
     http://localhost:8080/api/v1/upload
```

### Upload a File

```bash
curl -X POST http://localhost:8080/api/v1/upload \
  -H "Authorization: Bearer $TOKEN" \
  -F "file=@document.pdf"
```

**Response:**
```json
{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "url": "/files/123e4567-e89b-12d3-a456-426614174000",
  "size": 1048576,
  "content_type": "application/pdf",
  "upload_time": "2024-01-15T10:30:00Z",
  "checksum": "sha256:a665a45920422f9d417e4867efdc4fb8a04a1f3fff1fa07e998e86f7f7a27ae3"
}
```

### Download a File

```bash
# Download file content
curl -H "Authorization: Bearer $TOKEN" \
     http://localhost:8080/files/$FILE_ID \
     --output downloaded-file.pdf

# Get metadata only
curl -H "Authorization: Bearer $TOKEN" \
     -H "Accept: application/json" \
     http://localhost:8080/files/$FILE_ID
```

### Health Check

```bash
curl http://localhost:8080/health
```

## Development 🛠️

### Available Make Commands

```bash
# Build the application
make build

# Run the application
make run

# Run tests
make test

# Run tests with coverage
make test-coverage

# Format code
make fmt

# Lint code
make lint

# Generate development JWT token
make generate-token

# Clean build artifacts
make clean
```

### Development Workflow

1. **Setup development environment:**
```bash
make dev-setup
```

2. **Generate a test JWT token:**
```bash
make generate-token
export TOKEN="<generated-token>"
```

3. **Test the API:**
```bash
# Upload a file
curl -X POST http://localhost:8080/api/v1/upload \
  -H "Authorization: Bearer $TOKEN" \
  -F "file=@example.pdf"

# Check health
curl http://localhost:8080/health
```

### Running Tests

```bash
# Run all tests
make test

# Run tests with coverage report
make test-coverage
# Opens coverage.html in browser
```



## Security 🔒

### Built-in Security Features

- **JWT Authentication**: Stateless token-based authentication
- **File Validation**: Multiple layers of file type and content validation
- **User Isolation**: Users can only access their own files
- **Rate Limiting**: Prevents abuse with configurable limits
- **Secure Headers**: CSP, CSRF protection, and other security headers
- **UUID File Names**: Prevents file enumeration attacks

### Security Headers

```
Content-Security-Policy: default-src 'self'
X-Content-Type-Options: nosniff
X-Frame-Options: DENY
X-XSS-Protection: 1; mode=block
```

## API Reference 📚

### Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/health` | Service health check |
| `GET` | `/ready` | Service readiness check |
| `POST` | `/api/v1/upload` | Upload a file |
| `GET` | `/files/{id}` | Download file or get metadata |

### Status Codes

| Code | Description |
|------|-------------|
| `200` | Success |
| `400` | Bad Request (invalid input, file too large) |
| `401` | Unauthorized (invalid JWT) |
| `404` | File not found |
| `413` | File too large |
| `415` | Unsupported file type |
| `429` | Rate limit exceeded |
| `500` | Internal server error |

### Rate Limiting

- **Default**: 60 requests per minute per user
- **Headers**: 
  - `X-RateLimit-Limit`: Max requests per window
  - `X-RateLimit-Remaining`: Remaining requests
  - `X-RateLimit-Reset`: Reset time



## Contributing 🤝

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Development Guidelines

- Run `make test` before submitting PRs
- Follow Go formatting standards (`make fmt`)
- Add tests for new features
- Update documentation as needed

## License 📄

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support 💬

- 📧 **Email**: your-email@example.com
- 🐛 **Issues**: [GitHub Issues](https://github.com/yourusername/fileuploader/issues)
- 📖 **Documentation**: [API Documentation](docs/api.md)

---

**Built with ❤️ using Go**