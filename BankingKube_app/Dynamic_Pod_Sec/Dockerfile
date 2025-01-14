FROM golang:1.23 AS builder

WORKDIR /app

# Copy only the Dynamic_Pod_Sec folder into /app/BankingKube_app/Dynamic_Pod_Sec
COPY BankingKube_app/Dynamic_Pod_Sec /app/BankingKube_app/Dynamic_Pod_Sec

# Debugging steps
RUN ls -la /app
RUN ls -la /app/BankingKube_app
RUN ls -la /app/BankingKube_app/Dynamic_Pod_Sec
RUN ls -la /app/BankingKube_app/Dynamic_Pod_Sec/cmd

WORKDIR /app/BankingKube_app/Dynamic_Pod_Sec

# Download deps, build, etc.
RUN go mod download
RUN go build -o webhook ./cmd/main.go

FROM alpine:latest
WORKDIR /app

COPY --from=builder /app/BankingKube_app/Dynamic_Pod_Sec/webhook /app/webhook
COPY --from=builder /app/BankingKube_app/Dynamic_Pod_Sec/configs /configs

EXPOSE 8443
CMD ["/app/webhook"]
