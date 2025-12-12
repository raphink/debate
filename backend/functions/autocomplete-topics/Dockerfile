FROM golang:1.24-alpine AS builder

# Install git and ca-certificates (required for go mod download and HTTPS)
RUN apk add --no-cache git ca-certificates

# Copy shared module first (required by replace directive)
COPY shared /shared

WORKDIR /app

# Copy function files
COPY functions/autocomplete-topics/go.mod functions/autocomplete-topics/go.sum* ./
RUN go mod download

# Copy source
COPY functions/autocomplete-topics/ .

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /autocomplete-topics ./cmd

# Runtime stage
FROM gcr.io/distroless/static-debian12

WORKDIR /

# Copy CA certificates from builder
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=builder /autocomplete-topics .

EXPOSE 8080

CMD ["./autocomplete-topics"]
