# Multi-stage Dockerfile for Aperture Science Data Portal
# Stage 1: Build SvelteKit frontend
FROM node:20-alpine AS frontend-builder

WORKDIR /frontend

# Copy package files and install dependencies
COPY frontend/package*.json ./
RUN npm ci

# Copy frontend source and build
COPY frontend/ ./
RUN npm run build

# Stage 2: Build Go backend
FROM golang:1.25-alpine AS backend-builder

WORKDIR /backend

# Install build dependencies
RUN apk add --no-cache git

# Copy go mod files and download dependencies
COPY backend/go.mod backend/go.sum ./
RUN go mod download

# Copy backend source
COPY backend/ ./

# Build the Go binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server ./cmd/api

# Stage 3: Final runtime image
FROM alpine:latest

# Install ca-certificates for HTTPS requests to external APIs
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy the built Go binary from backend-builder
COPY --from=backend-builder /backend/server .

# Copy the built frontend from frontend-builder
COPY --from=frontend-builder /frontend/build ./frontend/build

# Copy database migrations
COPY backend/db/migrations ./db/migrations

# Expose port (Railway will inject PORT env variable)
EXPOSE 8080

# Run the server
CMD ["./server"]
