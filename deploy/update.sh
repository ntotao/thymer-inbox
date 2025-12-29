#!/bin/bash
set -e

GREEN='\033[0;32m'
NC='\033[0m'

echo -e "${GREEN}Updating Thymer Inbox...${NC}"

cd "$HOME/thymer-inbox"

# Pull latest images
docker-compose pull

# Recreate containers
docker-compose up -d

echo -e "${GREEN}Update complete!${NC}"
