#!/bin/sh
set -eu

# Install and activate DevBuddy (bud).
#
# Usage:
#   curl -sSL https://raw.githubusercontent.com/devbuddy/devbuddy/main/install.sh > /tmp/install-bud.sh
#   sh /tmp/install-bud.sh          # install the binary
#   bud up                           # setup project dependencies
#   source /tmp/install-bud.sh       # activate the environment
#
# Environment variables:
#   VERSION      - version to install (e.g. "v0.15.0"), defaults to latest
#   INSTALL_DIR  - directory to install to, defaults to /usr/local/bin

REPO="devbuddy/devbuddy"
INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"

detect_os() {
    case "$(uname -s)" in
        Linux)  echo "linux" ;;
        Darwin) echo "darwin" ;;
        *)      echo "Unsupported OS: $(uname -s)" >&2; exit 1 ;;
    esac
}

detect_arch() {
    case "$(uname -m)" in
        x86_64|amd64)  echo "amd64" ;;
        aarch64|arm64) echo "arm64" ;;
        *)             echo "Unsupported architecture: $(uname -m)" >&2; exit 1 ;;
    esac
}

get_latest_version() {
    curl -sSL -o /dev/null -w "%{url_effective}" "https://github.com/${REPO}/releases/latest" | grep -oE "[^/]+$"
}

install_bud() {
    if command -v bud > /dev/null 2>&1; then
        return
    fi

    OS=$(detect_os)
    ARCH=$(detect_arch)
    VERSION="${VERSION:-$(get_latest_version)}"
    BINARY="bud-${OS}-${ARCH}"
    URL="https://github.com/${REPO}/releases/download/${VERSION}/${BINARY}"

    echo "Installing DevBuddy ${VERSION} (${OS}/${ARCH})..."

    tmpfile=$(mktemp)
    trap 'rm -f "$tmpfile"' EXIT

    if ! curl -sSfL -o "$tmpfile" "$URL"; then
        echo "Failed to download ${URL}" >&2
        exit 1
    fi

    chmod +x "$tmpfile"

    if [ -w "$INSTALL_DIR" ]; then
        mv "$tmpfile" "${INSTALL_DIR}/bud"
    else
        sudo mv "$tmpfile" "${INSTALL_DIR}/bud"
    fi

    echo "DevBuddy installed to ${INSTALL_DIR}/bud"
}

install_bud
eval "$("${INSTALL_DIR}/bud" --shell-hook)"
