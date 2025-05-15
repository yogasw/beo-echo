# Stage 1: Build Frontend
FROM node:20-alpine AS frontend-builder

WORKDIR /app/frontend

# Copy frontend files
COPY frontend/package*.json ./
RUN npm install

COPY frontend/ ./
RUN npm run build

# Stage 2: Build Backend
FROM golang:1.24-alpine AS backend-builder

WORKDIR /app/backend

# Install build dependencies
RUN apk add --no-cache gcc musl-dev

# Copy backend files
COPY backend/ ./

# Build the Go application
RUN go mod download
RUN go build -o beo-echo .

# Stage 3: Production
FROM alpine:3.19

WORKDIR /app

# Install Caddy and necessary packages
RUN apk add --no-cache \
    ca-certificates \
    caddy

# Create necessary directories
RUN mkdir -p /app/frontend /app/backend /app/configs /data/caddy /config/caddy

# Copy built artifacts from previous stages
COPY --from=frontend-builder /app/frontend/build /app/frontend
COPY --from=backend-builder /app/backend/beo-echo /app/backend/

# Copy configuration files
COPY configs/ /app/configs/

# Create Caddy configuration file
COPY Caddyfile /etc/caddy/Caddyfile

# Expose ports
EXPOSE 80

# Create entry point script
COPY entrypoint.sh /app/entrypoint.sh
RUN chmod +x /app/entrypoint.sh

# Set environment variables
ENV GIN_MODE=release

ENTRYPOINT ["/app/entrypoint.sh"]
