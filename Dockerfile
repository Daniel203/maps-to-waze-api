# Step 1: Build the Go app in a Golang base image
FROM golang:1.21.0 AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main .

# Step 2: Create a minimal runtime image using Alpine
FROM alpine:latest
RUN apk --no-cache add ca-certificates 
WORKDIR /root/
COPY --from=builder /app/main .
COPY --from=builder /app/db/migrations /root/db/migrations
EXPOSE 8080
CMD ["/root/main"]

