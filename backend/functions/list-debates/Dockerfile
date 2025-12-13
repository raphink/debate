# Build stage
FROM golang:1.24-alpine AS builder

# Copy shared module first (required by replace directive)
COPY shared /shared

WORKDIR /app

# Copy go mod files
COPY functions/list-debates/go.mod functions/list-debates/go.sum* ./
RUN go mod download

# Copy source code
COPY functions/list-debates/ .

# Build the binary from cmd
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/list-debates ./cmd/main.go

# Runtime stage
FROM gcr.io/distroless/static-debian11

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/list-debates .

# Expose port
EXPOSE 8080

# Run the binary
ENTRYPOINT ["/app/list-debates"]
