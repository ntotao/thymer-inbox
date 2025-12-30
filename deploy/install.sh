#!/bin/bash
set -e

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${BLUE}=== Thymer Inbox Installer ===${NC}"

# Check for Git
if ! command -v git &> /dev/null; then
    echo -e "${RED}Error: git is not installed.${NC}"
    echo "Please install git first."
    exit 1
fi

# Set install directory
INSTALL_DIR="$HOME/thymer-inbox"
BRANCH="${1:-main}"
REPO_OWNER="${2:-riclib}"
REPO_URL="https://github.com/$REPO_OWNER/thymer-inbox.git"

echo -e "Installing to: ${GREEN}$INSTALL_DIR${NC}"

if [ -d "$INSTALL_DIR" ] && [ -d "$INSTALL_DIR/.git" ]; then
    echo "Updating existing installation..."
    cd "$INSTALL_DIR"
    git fetch origin
    git checkout "$BRANCH"
    git pull origin "$BRANCH"
else
    if [ -d "$INSTALL_DIR" ]; then
        echo "Cleaning up previous incomplete installation..."
        rm -rf "$INSTALL_DIR"
    fi
    echo "Cloning repository..."
    git clone -b "$BRANCH" "$REPO_URL" "$INSTALL_DIR"
    cd "$INSTALL_DIR"
fi

# Configure ports if set
if [ ! -z "$HTTP_PORT" ] || [ ! -z "$HTTPS_PORT" ]; then
    echo "Configuring custom ports..."
    echo "HTTP_PORT=${HTTP_PORT:-80}" > .env
    echo "HTTPS_PORT=${HTTPS_PORT:-443}" >> .env
fi

# Start services
echo "Starting services..."
docker-compose up -d --build

# Find local IP
IP=$(hostname -I | cut -d' ' -f1)
PORT="${HTTPS_PORT:-443}"
URL="https://$IP"
if [ "$PORT" != "443" ]; then
    URL="https://$IP:$PORT"
fi

echo -e "\n${GREEN}=== Installation Complete! ===${NC}"
echo -e "Server is running at: ${BLUE}$URL${NC}"
echo -e "\n${RED}IMPORTANT:${NC}"
echo "1. Open $URL in your browser"
echo "2. Accept the self-signed certificate warning"
echo "3. Configure your tokens in the Web GUI"
echo "4. Use $URL as the plugin URL in Thymer"
