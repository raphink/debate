# Build stage
FROM golang:1.24-alpine AS builder

# Install git and ca-certificates (required for go mod download and HTTPS)
RUN apk add --no-cache git ca-certificates

# Copy shared module first (required by replace directive)
COPY shared /shared

WORKDIR /app

# Copy go mod files
COPY functions/generate-debate/go.mod functions/generate-debate/go.sum* ./

# Download dependencies
RUN go mod download

# Copy source code
COPY functions/generate-debate/ .

# Build the function from cmd/ directory
RUN CGO_ENABLED=0 GOOS=linux go build -o /generate-debate ./cmd

# Runtime stage
FROM gcr.io/distroless/static-debian12

WORKDIR /

# Copy CA certificates from builder
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy binary from builder
COPY --from=builder /generate-debate /generate-debate

# Expose port
EXPOSE 8080

# Run the function
ENTRYPOINT ["/generate-debate"]
