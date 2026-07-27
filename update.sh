#!/usr/bin/env bash

# ==============================================================================
#  Matt File Manager Updater Script
#  Repo: https://github.com/Chintanpatel24/Matt
# ==============================================================================

set -e

BOLD="\033[1m"
GREEN="\033[32m"
YELLOW="\033[33m"
CYAN="\033[36m"
RESET="\033[0m"

echo -e "${CYAN}${BOLD}Updating Matt File Manager...${RESET}\n"

MATT_PATH=$(which matt 2>/dev/null || echo "$HOME/.local/bin/matt")

if command -v go &> /dev/null; then
    echo -e "${YELLOW}-> Pulling latest changes and rebuilding via Go...${RESET}"
    TMP_DIR=$(mktemp -d)
    git clone --depth 1 https://github.com/Chintanpatel24/Matt.git "$TMP_DIR"
    cd "$TMP_DIR"
    go build -o "$MATT_PATH" ./cmd/matt
    rm -rf "$TMP_DIR"
else
    OS=$(uname -s | tr '[:upper:]' '[:lower:]')
    ARCH=$(uname -m)
    case "$ARCH" in
        x86_64) ARCH="amd64" ;;
        aarch64|arm64) ARCH="arm64" ;;
    esac

    BINARY_URL="https://github.com/Chintanpatel24/Matt/releases/latest/download/matt-${OS}-${ARCH}"
    echo -e "${YELLOW}-> Fetching latest release binary...${RESET}"
    curl -sSL "$BINARY_URL" -o "$MATT_PATH"
fi

chmod +x "$MATT_PATH"

echo -e "\n${GREEN}${BOLD}✓ Matt updated successfully!${RESET}"
echo -e "Version info: ${CYAN}$(matt --version)${RESET}\n"
