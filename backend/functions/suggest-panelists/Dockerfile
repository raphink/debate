# Build stage
FROM golang:1.23-alpine AS builder

# Install git (required for go mod download)
RUN apk add --no-cache git

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum* ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the function
RUN CGO_ENABLED=0 GOOS=linux go build -o /suggest-panelists .

# Runtime stage
FROM gcr.io/distroless/static-debian12

WORKDIR /

# Copy binary from builder
COPY --from=builder /suggest-panelists /suggest-panelists

# Expose port
EXPOSE 8080

# Run the function
ENTRYPOINT ["/suggest-panelists"]
