# Multi-stage Dockerfile for Zap
# Produces a minimal container suitable for Kubernetes deployment

# Stage 1: Build the Go binary
FROM golang:1.24-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates

# Set working directory
WORKDIR /build

# Copy go mod files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY cmd/ ./cmd/

# Build static binary
# CGO_ENABLED=0 for static linking (works with scratch base)
# -ldflags="-s -w" strips debug info for smaller binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-s -w" \
    -o zap \
    ./cmd/main.go

# Stage 2: Runtime container
FROM alpine:latest

# Add ca-certificates for HTTPS redirects
RUN apk --no-cache add ca-certificates

# Create non-root user for security
RUN addgroup -g 1000 zap && \
    adduser -D -u 1000 -G zap zap

# Create directory for config (will be mounted from ConfigMap)
RUN mkdir -p /etc/zap && \
    chown -R zap:zap /etc/zap

# Copy binary from builder
COPY --from=builder /build/zap /usr/local/bin/zap

# Switch to non-root user
USER zap

# Expose default port (8927) and standard HTTP port
EXPOSE 8927 80

# Environment variables
# Disable /etc/hosts updates (not needed/possible in containers)
ENV ZAP_DISABLE_HOSTS_UPDATE=1

# Health check using the /healthz endpoint
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8927/healthz || exit 1

# Default command - users should mount config at /etc/zap/c.yml
# or override with their own command
CMD ["zap", "-host", "0.0.0.0", "-port", "8927", "-config", "/etc/zap/c.yml"]
