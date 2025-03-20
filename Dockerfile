# Use the official Golang image for building
FROM golang:1.24.1 AS builder

WORKDIR /app

# Copy Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire project
COPY . .

# Build the Go application statically
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Use a minimal base image for production
FROM alpine:latest

WORKDIR /root/

# Install required dependencies
RUN apk add --no-cache ca-certificates

# Copy the built binary from the builder stage
COPY --from=builder /app/main .
COPY .env .

# Ensure the binary has execute permissions
RUN chmod +x ./main

# Expose the application port
EXPOSE 8080

# Run the application
CMD ["./main"]
