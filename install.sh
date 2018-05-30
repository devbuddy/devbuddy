#!/usr/bin/env bash
set -eu

VERSION="v0.4.0"
DEST="/usr/local/bin"
SHELL_LINE='eval "$(dad --shell-init --with-completion)"'

YELLOW="\033[1;33m"
BLUE="\033[1;34m"
WHITE="\033[1;37m"
CODE="\033[44m\033[1;37m"
LINK="\033[4m\033[1;34m"
RESET="\033[0m"

make_variant() {
    ARCH=$(uname -m)
    if [[ "${ARCH}" != "x86_64" ]]; then
        echo >&2 "unsupported architecture: ${ARCH}"
        return 1
    fi

    case "$OSTYPE" in
    darwin*)
        echo "darwin-amd64"
        ;;
    linux*)
        echo "linux-amd64"
        ;;
    *)
        echo >&2 "unsupported OS: $OSTYPE"
        return 1
        ;;
    esac
}

header() {
    echo ""
    echo -e "${BLUE}$1${RESET}"
}

banner() {
    echo ""
    echo -e "${YELLOW}Welcome to Dad installer!${RESET}"
}

instructions() {
    echo ""
    echo -e "${YELLOW}Good!${RESET}\n"
    echo -e "${WHITE}Now, all you need is to add this to your bash .profile:${RESET}\n"
    echo -e "   ${CODE}" ${SHELL_LINE} "${RESET}\n"
    echo -e "Report any issue to ${LINK}https://github.com/pior/dad/issues${RESET}\n"
}

main() {
    banner

    TMPDIR=$(mktemp -d)
    cd "${TMPDIR}"

    header "Downloading binary from Github"
    BINARY="dad-$(make_variant)"
    URL="https://github.com/pior/dad/releases/download/${VERSION}/${BINARY}"
    curl -L -# --fail "${URL}" -o "${BINARY}"
    curl -L -# --fail "${URL}.sha256" -o "${BINARY}.sha256"

    header "Verify SHA256 checksum"
    shasum -c "${BINARY}.sha256"

    header "Installing to ${DEST}"
    sudo install "${BINARY}" "${DEST}/dad"

    instructions
}

main
