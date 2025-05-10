# -------- Stage 1: Build Stage (Backend + Frontend) --------
FROM node:22.15.0-alpine3.21 AS builder

WORKDIR /app

# -------- Install System Tools --------
RUN apk update && apk add --no-cache openssl

# Copy backend and frontend source code
COPY backend/ ./backend/
COPY frontend/ ./frontend/

# Install all dependencies (dev + prod)
RUN npm install --prefix backend && npm install --prefix frontend

# Set environment
ENV VITE_API_BASE_URL="/mock/api"

# Setup Folder and Configurations
RUN npm run setup --prefix backend

# Build backend (TypeScript) and frontend (Vite)
RUN npm run build --prefix backend
RUN npm run build --prefix frontend

# Generate Prisma Client
RUN npm run db:generate --prefix backend


# -------- Stage 2: Production Stage --------
FROM node:22.15.0-alpine3.21

# Set architecture environment from container
RUN apk update
RUN apk add --no-cache curl tar gzip openssl git openssh
RUN ARCH=$(uname -m | sed 's/x86_64/amd64/' | sed 's/aarch64/arm64/') \
        && echo "Detected architecture: $ARCH" \
        && TRAEFIK_VERSION=$(curl -s https://api.github.com/repos/traefik/traefik/releases/latest | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/' | sed 's/v//') \
        && echo "Latest Traefik version: $TRAEFIK_VERSION" \
        && curl -L "https://github.com/traefik/traefik/releases/download/v${TRAEFIK_VERSION}/traefik_v${TRAEFIK_VERSION}_linux_${ARCH}.tar.gz" \
        | tar -xz -C /usr/local/bin traefik \
        && chmod +x /usr/local/bin/traefik

# Remove unnecessary tools
RUN apk del --no-network curl tar gzip

WORKDIR /app

# Set environment for production
ENV NODE_ENV=production

# Initialize folders for backend and frontend
RUN mkdir -p backend frontend

# Copy Prisma schema
COPY --from=builder /app/backend/src/prisma/schema.prisma ./backend/src/prisma/schema.prisma

# Install only production dependencies
COPY --from=builder /app/backend/package*.json ./backend/
COPY --from=builder /app/frontend/package*.json ./frontend/
RUN npm ci --prefix backend
RUN npm ci --prefix frontend

# Install mockoon CLI
RUN npm install -g @mockoon/cli

# Copy configuration files
COPY --from=builder /app/configs/ ./configs/

# Copy backend build
COPY --from=builder /app/backend/dist/ ./backend/dist/

# Copy frontend build
COPY --from=builder /app/frontend/build/ ./frontend/build/

# (Optional) Copy root package.json if needed
COPY package.json ./

# Expose Traefik default port
EXPOSE 80

# Copy the entrypoint script
COPY docker-entrypoint.sh /app/docker-entrypoint.sh

# Set execute permission for the script
RUN chmod +x /app/docker-entrypoint.sh

# Use the entrypoint script
CMD ["/app/docker-entrypoint.sh"]
