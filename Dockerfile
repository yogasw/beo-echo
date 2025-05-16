# Accept NODE_VERSION as a build argument with a default value
ARG NODE_VERSION=22.14.0

# Stage 1: Build Frontend
FROM node:${NODE_VERSION}-alpine AS frontend-builder

# Set the working directory for the frontend build
ENV VITE_API_BASE_URL="/mock/api"

WORKDIR /app/frontend

# Copy package files first to leverage Docker layer caching
COPY frontend/package*.json ./

RUN npm install

# Copy the rest of the frontend code
COPY frontend/ ./

ENV NODE_ENV=production
# Build the frontend
RUN npm run build -- --no-sourcemap

# Stage 2: Build Backend
FROM golang:1.24-alpine AS backend-builder

WORKDIR /app/backend

# Install build dependencies
RUN apk add --no-cache gcc musl-dev

# Copy go.mod and go.sum first to leverage Docker layer caching
COPY backend/go.mod backend/go.sum ./
RUN go mod download

# Copy the rest of the backend code
COPY backend/ ./

# Build the Go application with optimizations
RUN go build -ldflags="-s -w" -o beo-echo .

# Stage 3: Production
FROM alpine:3.19 as production

WORKDIR /app

# Install only the necessary packages in a single layer
RUN apk add --no-cache ca-certificates caddy && \
    mkdir -p /app/frontend /app/backend /app/configs /app/logs /data/caddy /config/caddy

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

# # Set a non-root user for better security
# RUN addgroup -S beoecho && adduser -S -G beoecho beoecho && \
#     chown -R beoecho:beoecho /app /data /config && \
#     chmod -R 755 /app/logs

# # Use the non-root user
# USER beoecho

ENTRYPOINT ["/app/entrypoint.sh"]
