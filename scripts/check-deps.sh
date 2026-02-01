#!/bin/bash
# check-deps.sh - Check system dependencies for Anime-Sama Downloader
#
# Usage: ./scripts/check-deps.sh
#
# Returns exit code 0 if all required deps are present, 1 otherwise.

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo "🔍 Checking dependencies for Anime-Sama Downloader..."
echo ""

MISSING=0

# Check Go version
check_go() {
    if command -v go &> /dev/null; then
        GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
        GO_MAJOR=$(echo $GO_VERSION | cut -d. -f1)
        GO_MINOR=$(echo $GO_VERSION | cut -d. -f2)
        if [ "$GO_MAJOR" -ge 1 ] && [ "$GO_MINOR" -ge 22 ]; then
            echo -e "${GREEN}✓${NC} Go $GO_VERSION (required: 1.22+)"
        else
            echo -e "${RED}✗${NC} Go $GO_VERSION found, but 1.22+ required"
            MISSING=1
        fi
    else
        echo -e "${RED}✗${NC} Go not found (required: 1.22+)"
        MISSING=1
    fi
}

# Check ffmpeg (required for HLS/M3U8 download)
check_ffmpeg() {
    if command -v ffmpeg &> /dev/null; then
        FFMPEG_VERSION=$(ffmpeg -version 2>/dev/null | head -n1 | awk '{print $3}')
        echo -e "${GREEN}✓${NC} ffmpeg $FFMPEG_VERSION"
    else
        echo -e "${YELLOW}⚠${NC} ffmpeg not found (required for HLS/M3U8 downloads)"
        echo "  Install with: sudo apt install ffmpeg  (Debian/Ubuntu)"
        echo "                brew install ffmpeg      (macOS)"
        MISSING=1
    fi
}

# Check Node.js (optional, for frontend dev)
check_node() {
    if command -v node &> /dev/null; then
        NODE_VERSION=$(node --version)
        echo -e "${GREEN}✓${NC} Node.js $NODE_VERSION (optional, for frontend dev)"
    else
        echo -e "${YELLOW}○${NC} Node.js not found (optional, for frontend dev)"
    fi
}

# Check npm (optional, for frontend dev)
check_npm() {
    if command -v npm &> /dev/null; then
        NPM_VERSION=$(npm --version)
        echo -e "${GREEN}✓${NC} npm $NPM_VERSION (optional, for frontend dev)"
    else
        echo -e "${YELLOW}○${NC} npm not found (optional, for frontend dev)"
    fi
}

# Check Docker (optional)
check_docker() {
    if command -v docker &> /dev/null; then
        DOCKER_VERSION=$(docker --version | awk '{print $3}' | sed 's/,//')
        echo -e "${GREEN}✓${NC} Docker $DOCKER_VERSION (optional)"
    else
        echo -e "${YELLOW}○${NC} Docker not found (optional, for containerized deployment)"
    fi
}

echo "=== Required Dependencies ==="
check_go
check_ffmpeg

echo ""
echo "=== Optional Dependencies ==="
check_node
check_npm
check_docker

echo ""
if [ $MISSING -eq 0 ]; then
    echo -e "${GREEN}✓ All required dependencies are installed!${NC}"
    exit 0
else
    echo -e "${RED}✗ Some required dependencies are missing. Please install them.${NC}"
    exit 1
fi
