# Use the official Golang image for building the binary
FROM golang:1.23 AS builder

WORKDIR /app

# Copy the go.mod and go.sum files first to leverage Docker cache
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the source code
COPY . .

# Debugging step: List the contents of the /app directory
RUN ls -la /app

# Debugging step: List the contents of the /app/BankingKube_app directory
RUN ls -la /app/BankingKube_app

# Debugging step: List the contents of the /app/BankingKube_app/Dynamic_Pod_Sec directory
RUN ls -la /app/BankingKube_app/Dynamic_Pod_Sec

# Debugging step: List the contents of the /app/BankingKube_app/Dynamic_Pod_Sec/cmd directory
RUN ls -la /app/BankingKube_app/Dynamic_Pod_Sec/cmd

# Set the working directory for the Go module
WORKDIR /app/BankingKube_app/Dynamic_Pod_Sec

# Build the binary
RUN go build -o webhook ./cmd/main.go

# Use a minimal image for running the binary
FROM alpine:latest

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/BankingKube_app/Dynamic_Pod_Sec/webhook /app/webhook

# Copy the TLS certificates (these will be mounted as volumes)
COPY configs /configs

# Expose the port the webhook listens on
EXPOSE 8443

# Command to run the webhook
CMD ["/app/webhook"]