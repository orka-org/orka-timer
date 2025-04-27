FROM golang:1.21-alpine AS builder

# Install git and build dependencies
RUN apk add --no-cache git

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Final stage
FROM alpine:latest

WORKDIR /app

# Add basic security settings
RUN adduser -D appuser && \
    apk add --no-cache ca-certificates tzdata

# Copy binary and config files
COPY --from=builder /app/main .
COPY --from=builder /app/.env .

# Use non-root user
RUN chown -R appuser:appuser /app && \
    chmod +x /app/main
USER appuser

# Expose the port the app runs on
EXPOSE 3000

# Wait for MongoDB to be ready
CMD ["./main"]
