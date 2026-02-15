# DevBuddy

[![Go Report Card](https://goreportcard.com/badge/github.com/devbuddy/devbuddy)](https://goreportcard.com/report/github.com/devbuddy/devbuddy)
[![tests](https://github.com/devbuddy/devbuddy/workflows/tests/badge.svg)](https://github.com/devbuddy/devbuddy/actions?query=workflow%3Atests)
[![GitHub Release](https://img.shields.io/github/release/devbuddy/devbuddy.svg)](https://github.com/devbuddy/devbuddy/releases/latest)
[![GitHub Release Date](https://img.shields.io/github/release-date/devbuddy/devbuddy.svg)](https://github.com/devbuddy/devbuddy/releases/latest)
[![Gitter](https://img.shields.io/badge/Discussions%20on-Gitter-crimson.svg?logo=gitter&style=flat)](https://gitter.im/devbuddy)
![GitHub Release](https://img.shields.io/badge/license-MIT-green.svg)

Contents:
- [Install DevBuddy](#install)
- [Usage](#usage)
- [Tasks](docs/Config.md)
- [Contributing to DevBuddy](docs/CONTRIBUTING.md)

#### Note: contributors are welcome to join the project! 👩‍💻👨‍💻🤖

The project is evolving slowly, mostly because DevBuddy covers my current needs.

I would love to help people implement their languages/environments/tools/dev-flow.

## What is DevBuddy?

**DevBuddy** is an open-source implementation of an internal tool developed at
[Shopify](https://engineering.shopify.com) called "**Dev**".

The first goal of this tool is to automate the **setup** tasks required to work on a project.

With **DevBuddy**, pushing a change to a project you have never touched looks like this:

- `bud clone devbuddy/devbuddy`
- `bud up`
- `git commit`
- `bud test`
- `git push`

<br>
<p align="center"><img src="/docs/demo.gif?raw=true"/></p>

## Status and progress

**DevBuddy** is more useful for Python, Go and Node projects. More languages will be natively
supported. DevBuddy is still useful for languages without native support thanks to the **custom** task.

See the project config [documentation](docs/Config.md).

### Tasks

Python:
- Python environment (pyenv + virtualenv)
- Pip (requirements file)
- Pipenv (support for [Pipfile](https://github.com/pypa/pipfile))

Go:
- Go environment (with GOPATH)
- Go modules
- Dep (support for [Go Dep](https://github.com/golang/dep))

Node:
- Node environment
- Dependencies with NPM
- Dependencies with Yarn **planned**

Ruby: **planned**

Rust: **planned**

Others:
- Homebrew
- Apt
- Custom (conditional shell command)
- Docker Compose (manage a docker-compose setup): **planned**

### Features

- Notification when important files (e.g. `requirements.txt`) are updated locally
  (e.g. by `git pull`)
- A `help` command to guide a new developer based on `dev.yml`
- An `upgrade` command to auto-upgrade **DevBuddy**

### Supported code hosting platforms

- GitHub
- GitLab
- Bitbucket (with Git)

### Shell integration

- Bash
- Zsh

## Install

### Homebrew

Install the Homebrew tap:
```bash
$ brew install devbuddy/devbuddy/devbuddy
```

Install the latest release:
```bash
$ brew install devbuddy
```

Install from the master branch:
```bash
$ brew install --HEAD devbuddy
```

### Go install

Note: use this if your PATH includes the GOBIN path.

Latest release:
```bash
$ go install github.com/devbuddy/devbuddy/cmd/bud@latest
```

Specify a version, like `v0.11.1`:
```bash
$ go install github.com/devbuddy/devbuddy/cmd/bud@v0.11.1
```

### Manual

#### Select the version to download:

Latest:
```bash
$ VERSION=$(curl -Ls -o /dev/null -w %{url_effective} "https://github.com/devbuddy/devbuddy/releases/latest" | grep -oE "[^/]+$")
```

[Previous releases](https://github.com/devbuddy/devbuddy/releases):

```bash
$ VERSION="v0.12.0"
```

#### Download the binary

Linux Intel:
```bash
$ curl -L https://github.com/devbuddy/devbuddy/releases/download/${VERSION}/bud-linux-amd64 > /tmp/bud
```

Linux ARM:
```bash
$ curl -L https://github.com/devbuddy/devbuddy/releases/download/${VERSION}/bud-linux-arm64 > /tmp/bud
```

macOS Intel:
```bash
$ curl -L https://github.com/devbuddy/devbuddy/releases/download/${VERSION}/bud-darwin-amd64 > /tmp/bud
```

macOS Apple Silicon:
```bash
$ curl -L https://github.com/devbuddy/devbuddy/releases/download/${VERSION}/bud-darwin-arm64 > /tmp/bud
```

#### Install the binary

```bash
$ sudo install /tmp/bud /usr/local/bin/bud
```

## Setup

★ Install shell integration (in `~/.bash_profile` or `~/.zshrc`):
```bash
eval "$(bud --shell-init --with-completion)"
```

A safer version:
```bash
type bud > /dev/null 2> /dev/null && eval "$(bud --shell-init --with-completion)"
```

### Configuration

If you usually work with repos from the same organization (like your personal one), you can set it as a default:

```
export BUD_DEFAULT_ORG="google"
```

Then you can use it directly to create, clone, and jump to those repos:
```bash
$ bud clone pytruth
```
Rather than:
```bash
$ bud clone google/pytruth
```

## Usage

★ Add a `dev.yml` file in your project:
```yaml
up:
  - go: 1.9.2
  - golang_dep
  - python: 3.6.4rc1
  - pip:
    - requirements.txt

commands:
  test:
    desc: Run the tests
    run: script/test
  lint:
    desc: Lint the project
    run: script/lint

open:
  staging: https://staging.myapp.com
  doc: https://godoc.org/github.com/org/myapp
```
See DevBuddy's own [dev.yml](dev.yml).

```bash
$ bud
Usage:
  bud [flags]
  bud [command]

Available Commands:
  cd          Jump to a local project
  clone       Clone a project from github.com
  create      Create a new project
  help        Help about any command
  inspect     Inspect the project and its tasks
  open        Open a link about your project
  up          Ensure the project is up and running
  upgrade     [experimental] Upgrade DevBuddy to the latest available release.

Flags:
  -h, --help              help for bud
      --shell-init        Shell initialization
      --version           version for bud
      --with-completion   Enable completion during initialization

Use "bud [command] --help" for more information about a command.
```

## License

[MIT](https://github.com/devbuddy/devbuddy/blob/master/LICENSE)

Authors:
- Pior Bastida (pior@pbastida.net)
- Mathieu Leduc-Hamel (mathieu.leduchamel@shopify.com)
- Emmanuel Milou <manumilou@mykolab.com>
- John Duff @jduff
