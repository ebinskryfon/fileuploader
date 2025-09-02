# docs/api.md
# FileUploader API Documentation

## Overview

The FileUploader microservice provides secure file upload and download capabilities with JWT authentication, file validation, and metadata management.

**Base URL:** `http://localhost:8080`

## Authentication

All API endpoints require JWT authentication. Include the token in the Authorization header:

```
Authorization: Bearer <your-jwt-token>
```

## Rate Limiting

- Default: 60 requests per minute per authenticated user
- Rate limit headers are included in responses:
  - `X-RateLimit-Limit`: Maximum requests per window
  - `X-RateLimit-Remaining`: Remaining requests in current window
  - `X-RateLimit-Reset`: Time when the rate limit resets

## Error Responses

All endpoints return errors in the following format:

```json
{
  "error": "Error message",
  "code": 400,
  "message": "Detailed error description"
}
```

Common HTTP status codes:
- `400` - Bad Request (invalid input, file too large, etc.)
- `401` - Unauthorized (missing or invalid JWT token)
- `404` - Not Found (file not found)
- `413` - Payload Too Large (file exceeds size limit)
- `415` - Unsupported Media Type (invalid file type)
- `429` - Too Many Requests (rate limit exceeded)
- `500` - Internal Server Error

## Endpoints

### Health Check

#### GET /health

Returns service health status.

**Response:**
```json
{
  "status": "healthy",
  "timestamp": "2024-01-15T10:30:00Z",
  "service": "fileuploader"
}
```

#### GET /ready

Returns service readiness status with dependency checks.

**Response:**
```json
{
  "status": "ready",
  "checks": {
    "storage": "ok"
  }
}
```

### File Upload

#### POST /api/v1/upload

Uploads a file and returns metadata.

**Request:**
- Method: `POST`
- Content-Type: `multipart/form-data`
- Authentication: Required

**Form Parameters:**
- `file` (required): The file to upload

**File Constraints:**
- Maximum size: 25MB (configurable)
- Allowed types: JPEG, PNG, PDF (configurable)
- Original filename preserved in metadata

**Example Request:**
```bash
curl -X POST http://localhost:8080/api/v1/upload \
  -H "Authorization: Bearer <token>" \
  -F "file=@document.pdf"
```

**Success Response (200 OK):**
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

**Error Responses:**
- `400` - Invalid file type or size
- `401` - Authentication required
- `413` - File too large
- `415` - Unsupported file type
- `429` - Rate limit exceeded

### File Download

#### GET /files/{id}

Downloads a file or retrieves its metadata.

**Request:**
- Method: `GET`
- Authentication: Required

**Path Parameters:**
- `id` (required): File ID returned from upload

**Headers:**
- `Accept: application/json` - Returns metadata only
- `Accept: */*` or specific MIME type - Returns file content

**Example Requests:**

Get metadata:
```bash
curl -H "Authorization: Bearer <token>" \
     -H "Accept: application/json" \
     http://localhost:8080/files/123e4567-e89b-12d3-a456-426614174000
```

Download file:
```bash
curl -H "Authorization: Bearer <token>" \
     http://localhost:8080/files/123e4567-e89b-12d3-a456-426614174000 \
     --output downloaded-file.pdf
```

**Metadata Response (200 OK):**
```json
{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "original_name": "document.pdf",
  "size": 1048576,
  "content_type": "application/pdf",
  "upload_time": "2024-01-15T10:30:00Z",
  "url": "/files/123e4567-e89b-12d3-a456-426614174000",
  "checksum": "sha256:a665a45920422f9d417e4867efdc4fb8a04a1f3fff1fa07e998e86f7f7a27ae3",
  "user_id": "user-123"
}
```

**File Content Response (200 OK):**
- Content-Type: Original file MIME type
- Content-Disposition: `attachment; filename="original-name.ext"`
- Content-Length: File size in bytes
- Body: File content (binary)

**Error Responses:**
- `401` - Authentication required
- `404` - File not found or access denied
- `429` - Rate limit exceeded

## Security Features

### File Validation
- MIME type checking against whitelist
- File extension validation
- Content sniffing for type verification
- Size limit enforcement

### Access Control
- JWT-based authentication
- User isolation (users can only access their own files)
- File ID is UUID to prevent enumeration

### Security Headers
```
Content-Security-Policy: default-src 'self'
X-Content-Type-Options: nosniff
X-Frame-Options: DENY
X-XSS-Protection: 1; mode=block
```

## Configuration

### File Size Limits
- Default: 25MB
- Configurable via `MAX_FILE_SIZE` environment variable
- Returns `413 Payload Too Large` when exceeded

### Allowed File Types
Default allowed MIME types:
- `image/jpeg`
- `image/png` 
- `application/pdf`

Configurable via application configuration.

### Storage
- Files stored with UUID names to prevent conflicts
- Metadata stored separately as JSON
- Configurable storage path
- Automatic directory creation

## Examples

### Complete Upload/Download Workflow

1. **Generate JWT Token** (development):
```bash
go run scripts/generate_token.go
```

2. **Upload File**:
```bash
TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
curl -X POST http://localhost:8080/api/v1/upload \
  -H "Authorization: Bearer $TOKEN" \
  -F "file=@example.pdf" \
  | jq .
```

3. **Get File Metadata**:
```bash
FILE_ID="123e4567-e89b-12d3-a456-426614174000"
curl -H "Authorization: Bearer $TOKEN" \
     -H "Accept: application/json" \
     http://localhost:8080/files/$FILE_ID \
     | jq .
```

4. **Download File**:
```bash
curl -H "Authorization: Bearer $TOKEN" \
     http://localhost:8080/files/$FILE_ID \
     --output downloaded.pdf
```

### Error Handling Examples

**Invalid file type:**
```bash
curl -X POST http://localhost:8080/api/v1/upload \
  -H "Authorization: Bearer $TOKEN" \
  -F "file=@malware.exe"
```

Response:
```json
{
  "error": "Invalid file type",
  "code": 400,
  "message": "Invalid file type"
}
```

**File too large:**
```bash
# Upload file larger than limit
curl -X POST http://localhost:8080/api/v1/upload \
  -H "Authorization: Bearer $TOKEN" \
  -F "file=@huge-file.pdf"
```

Response:
```json
{
  "error": "File too large",
  "code": 400,
  "message": "File too large"
}
```

**Unauthorized access:**
```bash
curl -X POST http://localhost:8080/api/v1/upload \
  -F "file=@example.pdf"
```

Response:
```json
{
  "error": "Unauthorized",
  "code": 401,
  "message": "Unauthorized"
}
```

## Development

### Running Locally
```bash
# Set environment variables
export JWT_SECRET=your-secret-key
export PORT=8080

# Run service
make run
```

### Testing with curl
```bash
# Generate development token
make generate-token

# Use token for API calls
export TOKEN="<generated-token>"
```




The service will be available at `http://localhost:8080`.