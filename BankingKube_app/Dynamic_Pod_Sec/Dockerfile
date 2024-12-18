# Use the official Golang image for building the binary
FROM golang:1.19 AS builder

WORKDIR /app

# Copy the source code
COPY . .

# Download dependencies and build the binary
RUN go mod tidy
RUN go build -o webhook ./cmd/main.go

# Use a minimal image for running the binary
FROM alpine:latest

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/webhook /app/webhook

# Copy the TLS certificates (these will be mounted as volumes)
COPY configs /configs

# Expose the port the webhook listens on
EXPOSE 8443

# Command to run the webhook
CMD ["/app/webhook"]
