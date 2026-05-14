# === Stage 1: Build ===
FROM golang:1.23-alpine AS builder

WORKDIR /build

# Copy go mod files first for better layer caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code and build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o gophertodo ./cmd/server

# === Stage 2: Runtime ===
FROM alpine:3.19

# Install ca-certificates for HTTPS calls and create non-root user
RUN apk --no-cache add ca-certificates tzdata && \
    adduser -D -u 10001 -s /bin/sh appuser

WORKDIR /app

# Copy binary from builder
COPY --from=builder /build/gophertodo .
COPY --from=builder /build/data ./data

# Use external data file path via env var with fallback
ENV DATA_FILE=/app/data/tasks.json
ENV PORT=8080

# Expose port
EXPOSE 8080

# Switch to non-root user
USER appuser

# Run the binary
CMD ["./gophertodo"]