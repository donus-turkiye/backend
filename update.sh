#!/bin/bash

cd ~/backend
docker compose -f docker-compose.prod.yaml --env-file .env pull
docker compose -f docker-compose.prod.yaml --env-file .env down
docker compose -f docker-compose.prod.yaml --env-file .env up -d

echo "Backend updated successfully."
