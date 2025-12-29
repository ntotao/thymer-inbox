#!/bin/bash
set -e

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${BLUE}=== Thymer Inbox Installer ===${NC}"

# Check for Docker
if ! command -v docker &> /dev/null; then
    echo -e "${RED}Error: Docker is not installed.${NC}"
    echo "Please install Docker and Docker Compose first."
    exit 1
fi

# Set install directory
INSTALL_DIR="$HOME/thymer-inbox"
BRANCH="${1:-main}"
REPO_OWNER="${2:-riclib}"
REPO_URL="https://raw.githubusercontent.com/$REPO_OWNER/thymer-inbox/$BRANCH"

echo -e "Installing to: ${GREEN}$INSTALL_DIR${NC}"
mkdir -p "$INSTALL_DIR"
cd "$INSTALL_DIR"

# Download files
echo "Downloading configuration..."
curl -sSL "$REPO_URL/docker-compose.yml" -o docker-compose.yml
curl -sSL "$REPO_URL/Caddyfile" -o Caddyfile

# Start services
echo "Starting services..."
docker-compose up -d

# Find local IP
IP=$(hostname -I | cut -d' ' -f1)

echo -e "\n${GREEN}=== Installation Complete! ===${NC}"
echo -e "Server is running at: ${BLUE}https://$IP${NC}"
echo -e "\n${RED}IMPORTANT:${NC}"
echo "1. Open https://$IP in your browser"
echo "2. Accept the self-signed certificate warning"
echo "3. Configure your tokens in the Web GUI"
echo "4. Use https://$IP as the plugin URL in Thymer"
