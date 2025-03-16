# Build stage
FROM golang:1.23.4-alpine AS builder

WORKDIR /app

# Install necessary build tools
RUN apk add --no-cache git

# Copy go mod files first to leverage Docker cache
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Final stage
FROM alpine:3.19

WORKDIR /app

# Add non-root user for security
RUN adduser -D appuser
USER appuser

# Copy only the binary from builder
COPY --from=builder /app/main .

# Copy the config file
COPY config/config.yaml .

# Expose the port your application runs on (adjust as needed)
EXPOSE 8080

# Run the binary
CMD ["./main"]