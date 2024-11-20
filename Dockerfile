  # Use an official Go image as a base
  FROM golang:1.21 AS builder

  # Set the working directory
  WORKDIR /app

  # Copy Go modules and dependencies
  COPY go.mod go.sum ./
  RUN go mod download

  # Copy the application code
  COPY . ./

  # Build the Go application
  RUN go build -o go-ticket-booking

  # Use a minimal image for the final container
  FROM alpine:latest

  # Set the working directory
  WORKDIR /root/

  # Copy the built binary from the builder
  COPY --from=builder /app/go-ticket-booking .

  # Expose the port your app runs on
  EXPOSE 8080

  # Command to run the executable
  CMD ["./go-ticket-booking"]
