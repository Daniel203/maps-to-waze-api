# syntax=docker/dockerfile:1

# Step 1: Build the Go app in a Golang base image
FROM golang:1.21.0 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files for Go dependency management
COPY go.mod go.sum ./

# Download Go dependencies
RUN go mod download

# Copy the entire source code into the container
COPY . .

# Build the Go app as a statically linked binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main .

# Step 2: Create a minimal runtime image using Alpine
FROM alpine:latest

# Install SSL certificates, which might be needed for Go's HTTPS support
RUN apk --no-cache add ca-certificates

# Set the working directory for the app in the runtime image
WORKDIR /root/

# Copy the statically built binary from the builder image
COPY --from=builder /app/main .

# Expose port 8080 (Cloud Run expects this port)
EXPOSE 8080

# Set the entrypoint to run the Go binary
CMD ["/root/main"]

