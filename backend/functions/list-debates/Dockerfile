# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum* ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary from cmd/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o list-debates ./cmd/main.go

# Runtime stage
FROM gcr.io/distroless/static-debian11

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/list-debates .

# Expose port
EXPOSE 8080

# Run the binary
ENTRYPOINT ["/app/list-debates"]
