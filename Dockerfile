# Build stage
FROM golang:1.25-trixie AS builder

# Install git and ca-certificates (needed for fetching dependencies)
RUN apt-get update && apt-get install -y git ca-certificates && rm -rf /var/lib/apt/lists/*

# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/api

# Final stage
FROM debian:trixie-slim

# Install ca-certificates for HTTPS requests
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

# Create non-root user
RUN groupadd --gid 1001 appgroup && \
    useradd --uid 1001 --gid appgroup --shell /bin/bash --create-home appuser

# Create app directory
RUN mkdir -p /app && chown appuser:appgroup /app

# Copy the binary from builder stage
COPY --from=builder /app/main /app/main

# Set execute permissions
RUN chmod +x /app/main

# Switch to non-root user and set working directory
USER appuser
WORKDIR /app

# Expose port
EXPOSE 8080

# Run the application
CMD ["./main"]