#!/bin/bash

# ASSET MANAGEMENT APP DEPLOYMENT SCRIPT

echo "Deploying Asset management App Server..."
docker-compose -p asset_management_app down -v

echo "Build container ...."
docker-compose -p asset_management_app up -d --build

echo "Deployment complete!"
