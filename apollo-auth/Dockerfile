# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /build

RUN apk add --no-cache git

COPY shared/ /shared/

COPY apollo-auth/go.mod apollo-auth/go.sum ./
RUN go mod download

COPY apollo-auth/ .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o apollo-auth .

# Final stage
FROM alpine:latest

WORKDIR /root/

RUN apk --no-cache add ca-certificates

COPY --from=builder /build/apollo-auth .

EXPOSE 4000

CMD ["./apollo-auth"]
