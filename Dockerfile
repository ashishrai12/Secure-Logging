# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install dependencies
COPY go.mod ./
RUN go mod download

# Copy source
COPY src/*.go ./src/

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /secure-logging ./src/*.go

# Run stage
FROM alpine:latest

WORKDIR /

COPY --from=builder /secure-logging /secure-logging

EXPOSE 5000

CMD ["/secure-logging"]
