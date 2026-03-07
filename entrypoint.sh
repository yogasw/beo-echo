#!/bin/sh
set -e

# Start the backend service in the background
cd /app/backend
echo "Starting Beo Echo backend service..."
./beo-echo server &
BACKEND_PID=$!

# Give a moment for the backend to initialize
echo "Waiting for backend to generate Caddyfile..."
attempts=0
while [ ! -f /app/configs/caddy/Caddyfile ] && [ $attempts -lt 30 ]; do
  sleep 1
  attempts=$((attempts + 1))
done

if [ ! -f /app/configs/caddy/Caddyfile ]; then
  echo "Error: Caddyfile was not generated within 30 seconds."
  exit 1
fi

# Start Caddy in the foreground
echo "Starting Caddy server..."
caddy run --config /app/configs/caddy/Caddyfile &
CADDY_PID=$!

PORT_USED=${PORT:-80}
echo "Server URL: http://0.0.0.0:$PORT_USED"

# Handle shutdown signals
trap 'kill $BACKEND_PID $CADDY_PID; exit 0' SIGTERM SIGINT

# Keep the container running
wait
