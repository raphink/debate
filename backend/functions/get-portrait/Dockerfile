# Multi-stage build for get-portrait Cloud Function
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum* ./
RUN go mod download

# Copy source code
COPY *.go ./

# Build the function
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o get-portrait .

# Final stage - minimal runtime image
FROM gcr.io/distroless/static-debian11

WORKDIR /

# Copy binary from builder
COPY --from=builder /app/get-portrait /get-portrait

# Set environment
ENV PORT=8083

# Run the function
ENTRYPOINT ["/get-portrait"]
