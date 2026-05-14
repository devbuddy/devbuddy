# DevBuddy

[![Go Report Card](https://goreportcard.com/badge/github.com/devbuddy/devbuddy)](https://goreportcard.com/report/github.com/devbuddy/devbuddy)
[![tests](https://github.com/devbuddy/devbuddy/workflows/tests/badge.svg)](https://github.com/devbuddy/devbuddy/actions?query=workflow%3Atests)
[![GitHub Release](https://img.shields.io/github/release/devbuddy/devbuddy.svg)](https://github.com/devbuddy/devbuddy/releases/latest)
![License](https://img.shields.io/badge/license-MIT-green.svg)

**DevBuddy** is a command-line tool that automates development environment setup and provides project-specific commands. Define your project's requirements in a `dev.yml` file and let `bud` handle the rest.

With DevBuddy, getting started on a project you've never touched looks like this:

```bash
bud clone devbuddy/devbuddy
bud up
# hack hack hack
bud test
```

<br>
<p align="center"><img src="/docs/demo.gif?raw=true"/></p>

## Supported Tasks

DevBuddy manages environments and dependencies through tasks defined in `dev.yml`:

**Languages:**
- **Python** — version management (pyenv + virtualenv), pip, Pipfile
- **Go** — version management, Go modules
- **Node.js** — version management (including Apple Silicon), npm

**System & Environment:**
- **Homebrew** — install macOS packages
- **Apt** — install Debian/Ubuntu packages
- **Environment variables** — set project-specific env vars
- **Env files** — load variables from `.env` files
- **Custom tasks** — conditional shell commands (met?/meet pattern)

See the full [task documentation](docs/Config.md) for details.

## Features

- Automatic environment activation/deactivation as you `cd` between projects
- Notification when important files (e.g. `requirements.txt`) are updated locally
- `bud open` to jump to the GitHub repo page, or `bud open <name>` for project URLs
- `bud clone` / `bud cd` for fast project navigation
- Shell completion (bash and zsh)

### Supported platforms

- macOS (Intel and Apple Silicon)
- Linux (amd64 and arm64)

### Supported shells

- Bash
- Zsh

## Install

### Quick install (CI-friendly)

```bash
curl -sSL https://raw.githubusercontent.com/devbuddy/devbuddy/main/install.sh | sh
```

Pin a specific version:
```bash
curl -sSL https://raw.githubusercontent.com/devbuddy/devbuddy/main/install.sh | VERSION=v0.15.0 sh
```

Choose a custom install directory (defaults to `/usr/local/bin`):
```bash
curl -sSL https://raw.githubusercontent.com/devbuddy/devbuddy/main/install.sh | INSTALL_DIR=./bin sh
```

### Homebrew (macOS)

```bash
brew install devbuddy/devbuddy/devbuddy
```

### Go install

Requires Go and `GOPATH/bin` in your PATH:

```bash
go install github.com/devbuddy/devbuddy/cmd/bud@latest
```

## Setup

Add shell integration to your `~/.bash_profile` or `~/.zshrc`:

```bash
eval "$(bud --shell-init --with-completion)"
```

A safer version that only activates if `bud` is installed:
```bash
type bud > /dev/null 2> /dev/null && eval "$(bud --shell-init --with-completion)"
```

### Configuration

If you usually work with repos from the same GitHub organization, set it as a default:

```bash
export BUD_DEFAULT_ORG="myorg"
```

Then `bud clone myrepo` is equivalent to `bud clone myorg/myrepo`.

## Using in CI

DevBuddy can run your project commands in CI without full shell integration. Install the binary, then activate the environment with `eval "$(bud --shell-hook)"` before running any `bud` command:

```bash
curl -sSL https://raw.githubusercontent.com/devbuddy/devbuddy/main/install.sh | sh
eval "$(bud --shell-hook)"
bud up
bud test
```

Example GitHub Actions step:
```yaml
- name: Run tests
  run: |
    curl -sSL https://raw.githubusercontent.com/devbuddy/devbuddy/main/install.sh | sh
    eval "$(bud --shell-hook)"
    bud up
    bud test
```

## Usage

Add a `dev.yml` file to your project root:

```yaml
env:
  DATABASE_URL: postgres://localhost:5432/myapp_dev

up:
  - go:
      version: '1.22'
  - python: 3.12.0
  - pip:
    - requirements.txt
  - homebrew:
    - curl
  - custom:
      name: Download GeoIP db
      met?: test -e GeoIP.dat
      meet: curl -sL http://example.com/GeoIP.dat.gz | gunzip > GeoIP.dat

commands:
  test:
    desc: Run the tests
    run: script/test
  lint:
    desc: Lint the project
    run: script/lint

open:
  staging: https://staging.myapp.com
  ci: https://github.com/org/myapp/actions
```

See DevBuddy's own [dev.yml](dev.yml) for a real-world example.

## Contributing

Contributions are welcome! See [CONTRIBUTING.md](docs/CONTRIBUTING.md) for details.

## License

[MIT](https://github.com/devbuddy/devbuddy/blob/master/LICENSE)

Authors:
- Pior Bastida (pior@pbastida.net)
- Mathieu Leduc-Hamel (mathieu.leduchamel@shopify.com)
- Emmanuel Milou <manumilou@mykolab.com>
- John Duff @jduff
