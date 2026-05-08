# Build stage
FROM golang:1.25-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the binaries
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o 4g-proxy ./cmd/app.go
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o dropctl ./cmd/dropctl/main.go

# Runtime stage
FROM alpine:3.19

# Install runtime dependencies
RUN apk add --no-cache ca-certificates

# Create non-root user
RUN adduser -D -g '' appuser

# Set working directory
WORKDIR /app

# Copy binaries from builder
COPY --from=builder /app/4g-proxy .
COPY --from=builder /app/dropctl .

# Copy default config
COPY --from=builder /app/config/config.yaml ./config/

# Change ownership
RUN chown -R appuser:appuser /app

# Switch to non-root user
USER appuser

# Expose ports
# S1AP SCTP port
EXPOSE 36412/sctp
# HTTP API port
EXPOSE 8080/tcp

# Default environment variables
ENV MME_ADDRESS=127.0.0.1
ENV MME_PORT=36412
ENV API_PORT=8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Run the proxy
ENTRYPOINT ["./4g-proxy"]
