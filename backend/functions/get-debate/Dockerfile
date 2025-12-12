# Multi-stage build for get-debate Cloud Function

FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum* ./
COPY ../../shared ../shared

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

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
