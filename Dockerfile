# ── Stage 1: build ───────────────────────────────────────────────────────────
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Download dependencies first — this layer is cached unless go.mod/go.sum change
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# -w -s strips debug info → noticeably smaller binary
# CGO_ENABLED=0 gives a fully static binary (no glibc dependency)
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/main .


# ── Stage 2: run ─────────────────────────────────────────────────────────────
# Pin the version — "alpine:latest" can change under you on a rebuild
FROM alpine:3.21

# ca-certificates is required for HTTPS calls to Google Maps / Geoapify
RUN apk --no-cache add ca-certificates \
    # Create a dedicated non-root user 
    && addgroup -S appgroup \
    && adduser  -S appuser -G appgroup

WORKDIR /app

COPY --from=builder /app/main          ./main
COPY --from=builder /app/db/migrations ./db/migrations

# Make sure the non-root user owns the files
RUN chown -R appuser:appgroup /app

USER appuser

EXPOSE 8080

CMD ["/app/main"]
