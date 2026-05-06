#!/bin/sh
set -eu

# Install DevBuddy (bud) binary.
#
# Usage:
#   curl -sSL https://raw.githubusercontent.com/devbuddy/devbuddy/main/install.sh | sh
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

main() {
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

main
