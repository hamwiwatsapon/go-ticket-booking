FROM golang:1.23-alpine AS builder

# Set destination for COPY
WORKDIR /app

# Install git if you have dependencies from git repositories
RUN apk add --no-cache git

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Copy the source code
COPY . .

# Verbose build with error checking
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v \
    -ldflags="-w -s" \
    -o main .

# Debug: list contents and check binary
RUN ls -la
RUN file main

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS connections
RUN apk add --no-cache ca-certificates tzdata

WORKDIR /root/

# Copy the pre-built binary file from the previous stage
COPY --from=builder /app/main .

# Optional: Add a non-root user for security
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
USER appuser

# Optional: Set a default port if your app uses one
EXPOSE 8080

# Add health check if applicable
# HEALTHCHECK --interval=30s --timeout=10s --start-period=5s \
#   CMD wget -q -O- http://localhost:8080/health || exit 1

# Command to run the executable
CMD ["./main"]