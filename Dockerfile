# Build stage
FROM golang:1.24-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags='-w -s' -o /app/bin/xis-data-aggregator ./cmd/xis-data-aggregator

# Final stage
FROM alpine:3.18

# Install CA certificates for SSL
RUN apk --no-cache add ca-certificates

# Set working directory
WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/bin/xis-data-aggregator .

# Expose ports for gRPC and HTTP
EXPOSE 50051 8080

# Set the entry point
ENTRYPOINT ["/app/xis-data-aggregator"]

