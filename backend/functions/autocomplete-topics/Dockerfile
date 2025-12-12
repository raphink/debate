FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum* ./
RUN go mod download

# Copy shared module
COPY ../../shared /app/backend/shared

# Copy source
COPY . .

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o autocomplete-topics ./cmd

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/autocomplete-topics .

EXPOSE 8080

CMD ["./autocomplete-topics"]
