# 1) Build stage
FROM golang:1.23 AS builder

WORKDIR /app

# Copy go.mod and go.sum to leverage caching
COPY go.mod go.sum ./
RUN go mod download

# Copy all source code from the current folder (Dynamic_Pod_Sec) to /app
COPY . .

# (Optional) Debug: list what we have in /app
RUN ls -la /app
RUN ls -la /app/cmd

# Build the webhook binary
RUN go build -o webhook ./cmd/main.go

# 2) Minimal runtime image
FROM alpine:latest

WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /app/webhook /app/webhook

# Copy configs if needed
COPY --from=builder /app/configs /configs

EXPOSE 8443
CMD ["/app/webhook"]
