# Multi-stage Dockerfile for daily-driver (Go + Templ + Tailwind static assets)
# Builder stage: use Go toolchain and run Templ codegen
FROM golang:1.24-bookworm AS builder

# Disable CGO for a static binary
ENV CGO_ENABLED=0

WORKDIR /src

# Install templ CLI matching go.mod version
RUN go install github.com/a-h/templ/cmd/templ@v0.3.943

# Cache modules first
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source
COPY . .

# Generate templ files (safe no-op if already generated)
RUN /go/bin/templ generate

# Build the server binary
RUN go build -trimpath -ldflags="-s -w" -o /out/server ./cmd

# Runtime stage: minimal image with certs and tzdata, non-root user
FROM alpine:3.20 AS runtime
RUN adduser -D -H -u 10001 appuser \
  && apk add --no-cache ca-certificates tzdata

WORKDIR /app

# Copy the binary and static assets
COPY --from=builder /out/server /app/server
COPY --from=builder /src/web/static /app/web/static

# Create writable log file location for the app
RUN touch /app/application.log && chown -R appuser:appuser /app

USER appuser

# Railway sets PORT; the app defaults to 8080 when PORT is not set
ENV PORT=8080
EXPOSE 8080

ENTRYPOINT ["/app/server"]

