.PHONY: build run test clean docker-build docker-run

# Build the application
build:
	go build -o bin/fileuploader cmd/server/main.go

# Run the application
run:
	go run cmd/server/main.go

# Run tests
test:
	go test -v ./...

# Run tests with coverage
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Clean build artifacts
clean:
	rm -rf bin/
	rm -f coverage.out coverage.html

# Format code
fmt:
	go fmt ./...

# Lint code
lint:
	golangci-lint run

# Build Docker image
docker-build:
	docker build -f docker/Dockerfile -t fileuploader:latest .

# Run Docker container
docker-run:
	docker run -p 8080:8080 \
		-e JWT_SECRET=your-secret-key \
		-e STORAGE_PATH=/app/storage \
		-v $(PWD)/storage:/app/storage \
		fileuploader:latest

# Development setup
dev-setup:
	go mod tidy
	mkdir -p storage
	mkdir -p bin

# Generate test token (for development)
generate-token:
	go run scripts/generate_token.go
