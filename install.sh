#!/usr/bin/env bash
set -eu

VERSION="v0.0.9"
DEST="/usr/local/bin"
SHELL_LINE='eval "$(dad --shell-init --with-completion)"'

YELLOW="\033[1;33m"
BLUE="\033[1;34m"
CYAN="\033[1;36m"
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

    header "Downloading binary from Github"
    URL="https://github.com/pior/dad/releases/download/${VERSION}/dad-$(make_variant)"
    TMPFILE=`mktemp`
    curl -L -# --fail "${URL}" -o "${TMPFILE}"

    header "Downloading SHA from Github"
    expected_hash=$(curl -LsS --fail "${URL}.sha256")
    downloaded_hash=$(shasum -a 256 "${TMPFILE}")
    echo "Expected hash   : ${expected_hash}"
    echo "Downloaded hash : ${downloaded_hash}"
    read -p "Correct? [enter]"

    header "Installing to ${DEST}"
    sudo install "${TMPFILE}" "${DEST}/dad"

    [[ -e "${TMPFILE}" ]] && unlink "${TMPFILE}"

    instructions
}

main
