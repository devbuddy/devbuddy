#!/usr/bin/env bash
set -euo pipefail

BUDPATH=${GOPATH}/bin/bud

go build -ldflags "-X main.Version=$(git describe --tags --dirty --always)" -o $BUDPATH ./cmd/bud

echo "Installed in ${BUDPATH}"

[ -e "/usr/local/bin/bud" ] && echo "WARNING: another bud command is installed in /usr/local/bin/bud" || true
[ -e "/usr/bin/bud" ] && echo "WARNING: another bud command is installed in /usr/bin/bud" || true
