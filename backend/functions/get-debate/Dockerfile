# Multi-stage build for get-debate Cloud Function

FROM golang:1.24-alpine AS builder

# Copy shared module first (required by replace directive)
COPY shared /shared

WORKDIR /app

# Copy go mod files
COPY functions/get-debate/go.mod functions/get-debate/go.sum* ./

# Download dependencies
RUN go mod download

# Copy source code
COPY functions/get-debate/ .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/get-debate ./cmd/main.go

# Runtime image
FROM alpine:latest

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/get-debate .

# Expose port
EXPOSE 8080

# Run
CMD ["./get-debate"]
