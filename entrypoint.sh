#!/bin/sh
set -e

# Start the backend service in the background
cd /app/backend
echo "Starting Beo Echo backend service..."
./beo-echo server &
BACKEND_PID=$!

# Give a moment for the backend to initialize
sleep 2

# Start Caddy in the foreground
echo "Starting Caddy server..."
caddy run --config /app/configs/caddy/Caddyfile &
CADDY_PID=$!

echo "Server URL: http://0.0.0.0:80"

# Handle shutdown signals
trap 'kill $BACKEND_PID $CADDY_PID; exit 0' SIGTERM SIGINT

# Keep the container running
wait
