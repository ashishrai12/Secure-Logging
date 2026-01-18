# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY src/ ./src/

# Build server
RUN CGO_ENABLED=0 GOOS=linux go build -o /secure-logging-server ./src/cmd/server/main.go

# Build client (optional but good to have)
RUN CGO_ENABLED=0 GOOS=linux go build -o /secure-logging-client ./src/cmd/client/main.go

# Run stage
FROM alpine:latest

WORKDIR /

COPY --from=builder /secure-logging-server /secure-logging-server
COPY --from=builder /secure-logging-client /secure-logging-client

EXPOSE 5000

CMD ["/secure-logging-server"]
