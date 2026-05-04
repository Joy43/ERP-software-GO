#!/bin/bash

# Quick Redis Cache Clear
# Fastest way to clear Redis cache

export $(cat .env | grep -v '#' | xargs) 2>/dev/null

REDIS_CONTAINER=$(docker ps --filter "name=redis" --format "{{.Names}}" | head -1)

if [ -z "$REDIS_CONTAINER" ]; then
    echo "❌ Redis container not found"
    exit 1
fi

echo "🔄 Clearing Redis cache..."
docker exec "$REDIS_CONTAINER" redis-cli -a "$REDIS_PASSWORD" FLUSHDB 2>/dev/null
echo "✅ Done! Cache cleared."
