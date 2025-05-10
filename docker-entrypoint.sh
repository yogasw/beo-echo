#!/bin/sh

# Check if DATABASE_URL is set
if [ -z "$DATABASE_URL" ]; then
  echo "DATABASE_URL is not set. Using default SQLite database."
  export DATABASE_URL="file:/app/configs/db/db.sqlite"
fi

# Update schema.prisma provider based on DATABASE_URL
if echo "$DATABASE_URL" | grep -q "postgresql://"; then
  PROVIDER="postgresql"
elif echo "$DATABASE_URL" | grep -q "mysql://"; then
  PROVIDER="mysql"
elif echo "$DATABASE_URL" | grep -q "mongodb://"; then
  PROVIDER="mongodb"
else
  PROVIDER="sqlite"
fi

echo "Updating schema.prisma with provider: $PROVIDER"
sed -i 's/provider = "sqlite"/provider = "'"$PROVIDER"'"/' /app/backend/src/prisma/schema.prisma

# Run database migration
echo "Running database migration..."
npm run db:migrate --prefix backend

# Start Traefik
traefik --configFile=/app/configs/traefik/traefik.yml &

# Start Backend
npm run start --prefix backend &

# Start Frontend
npx serve -s frontend/build -l 3005

# Wait for all background processes
wait
