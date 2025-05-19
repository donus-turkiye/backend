#!/bin/bash
set -e
cd ~/backend
cat .env
docker compose -f docker-compose.prod.yaml --env-file .env.prod pull
docker compose -f docker-compose.prod.yaml --env-file .env.prod down
docker compose -f docker-compose.prod.yaml --env-file .env.prod up -d

echo "Backend updated successfully."
