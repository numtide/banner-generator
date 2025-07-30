# Build stage
FROM golang:1.23-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git make

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the API server
RUN go build -o banner-api cmd/banner-api/main.go

# Runtime stage
FROM alpine:latest

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata rsvg-convert

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/banner-api .

# Copy static assets
COPY --from=builder /app/deploy/fonts ./deploy/fonts
COPY --from=builder /app/deploy/templates ./deploy/templates

# Copy configuration file
COPY --from=builder /app/deploy/banner-generator.toml ./deploy/

# Create a non-root user
RUN adduser -D -g '' appuser
USER appuser

# Expose the default port
EXPOSE 8080

# Run the API server
CMD ["./banner-api"]
