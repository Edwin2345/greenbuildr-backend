# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /build

# Install build dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o auth-service .

# Final stage
FROM alpine:latest

WORKDIR /root/

# Install ca-certificates for HTTPS support
RUN apk --no-cache add ca-certificates

# Copy binary from builder
COPY --from=builder /build/auth-service .

EXPOSE 4001

CMD ["./auth-service"]
