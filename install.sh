#!/usr/bin/env bash

# ==============================================================================
#  Matt File Manager Installer Script
#  Repo: https://github.com/Chintanpatel24/Matt
# ==============================================================================

set -e

BOLD="\033[1m"
GREEN="\033[32m"
WHITE="\033[97m"
CYAN="\033[36m"
RED="\033[31m"
RESET="\033[0m"

echo -e "${WHITE}${BOLD}"
echo '  ███▄ ▄███▓ ▄▄▄     ▄▄▄█████▓▄▄▄█████▓'
echo '  ▓██▒▀█▀ ██▒▒████▄   ▓  ██▒ ▓▒▓  ██▒ ▓▒'
echo '  ▓██    ▓██░▒██  ▀█▄ ▒ ▓██░ ▒░▒ ▓██░ ▒░'
echo '  ▒██    ▒██ ░██▄▄▄▄██░ ▓██▓ ░ ░ ▓██▓ ░ '
echo '  ▒██▒   ░██▒ ▓█   ▓██▒ ▒██▒ ░   ▒██▒ ░ '
echo '  ░ ▒░   ░  ░ ▒▒   ▓▒█░ ▒ ░░     ▒ ░░   '
echo '  ░  ░      ░  ▒   ▒▒ ░   ░        ░    '
echo '  ░      ░     ░   ▒    ░        ░      '
echo '         ░         ░  ░'
echo -e "${RESET}"
echo -e "${BOLD}Installing Matt Black Terminal File Manager...${RESET}\n"

INSTALL_DIR="/usr/local/bin"

if [ ! -w "$INSTALL_DIR" ]; then
    INSTALL_DIR="$HOME/.local/bin"
    mkdir -p "$INSTALL_DIR"
fi

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}" 2>/dev/null || echo ".")" && pwd)"

# Case 1: Local repository installation
if [ -f "$SCRIPT_DIR/cmd/matt/main.go" ] && command -v go &> /dev/null; then
    echo -e "${WHITE}-> Local repository clone detected. Building Matt from local source...${RESET}"
    go build -o "$INSTALL_DIR/matt" "$SCRIPT_DIR/cmd/matt"
# Case 2: Remote execution - build from Go
elif command -v go &> /dev/null; then
    echo -e "${WHITE}-> Go detected. Cloning and building Matt from GitHub source...${RESET}"
    TMP_DIR=$(mktemp -d)
    git clone --depth 1 https://github.com/Chintanpatel24/Matt.git "$TMP_DIR"
    cd "$TMP_DIR"
    go build -o "$INSTALL_DIR/matt" ./cmd/matt
    rm -rf "$TMP_DIR"
# Case 3: Prebuilt binary download
else
    OS=$(uname -s | tr '[:upper:]' '[:lower:]')
    ARCH=$(uname -m)

    case "$ARCH" in
        x86_64) ARCH="amd64" ;;
        aarch64|arm64) ARCH="arm64" ;;
        *) echo -e "${RED}Unsupported architecture: $ARCH${RESET}"; exit 1 ;;
    esac

    BINARY_URL="https://github.com/Chintanpatel24/Matt/releases/latest/download/matt-${OS}-${ARCH}"
    echo -e "${WHITE}-> Downloading Matt prebuilt binary (${OS}/${ARCH})...${RESET}"
    
    if command -v curl &> /dev/null; then
        curl -sSL "$BINARY_URL" -o "$INSTALL_DIR/matt"
    elif command -v wget &> /dev/null; then
        wget -qO "$INSTALL_DIR/matt" "$BINARY_URL"
    else
        echo -e "${RED}Neither curl nor wget found. Please install one to proceed.${RESET}"
        exit 1
    fi
fi

chmod +x "$INSTALL_DIR/matt"

echo -e "\n${GREEN}${BOLD}✓ Matt successfully installed to ${INSTALL_DIR}/matt${RESET}"
echo -e "Run ${CYAN}matt${RESET} in your terminal to start!\n"
