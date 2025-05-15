# Stage 1: Build Frontend
FROM node:20-alpine AS frontend-builder

WORKDIR /app/frontend

# Copy frontend files
COPY frontend/package*.json ./
RUN npm install

COPY frontend/ ./
# Increase Node.js memory limit for building
ENV NODE_OPTIONS="--max-old-space-size=2048"
# Use production mode to reduce memory usage
ENV NODE_ENV=production
# Use a more memory-efficient build command with lower concurrency
RUN npm run build -- --no-sourcemap

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
