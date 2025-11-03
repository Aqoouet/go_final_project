# Multi-stage build for Go Final Project Scheduler

# Stage 1: Build the application
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o scheduler ./cmd/scheduler

# Stage 2: Create minimal runtime image
FROM ubuntu:latest

# Install ca-certificates for HTTPS support
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/scheduler .

# Copy web directory
COPY web ./web

# Set environment variables with defaults
ENV TODO_PORT=7540
ENV TODO_DBFILE=/data/scheduler.db
ENV TODO_PASSWORD=""

# Expose the port
EXPOSE 7540

# Create data directory for database
RUN mkdir -p /data

# Run the application
CMD ["./scheduler"]

