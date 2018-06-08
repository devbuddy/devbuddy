# DevBuddy

[![Go Report Card](https://goreportcard.com/badge/github.com/devbuddy/devbuddy)](https://goreportcard.com/report/github.com/devbuddy/devbuddy)
[![CircleCI](https://circleci.com/gh/devbuddy/devbuddy.svg?style=svg)](https://circleci.com/gh/devbuddy/devbuddy)
[![GitHub Release](https://img.shields.io/github/release/devbuddy/devbuddy.svg)](https://github.com/devbuddy/devbuddy/releases/latest)
[![GitHub Release Date](https://img.shields.io/github/release-date/devbuddy/devbuddy.svg)](https://github.com/devbuddy/devbuddy/releases/latest)

## What is this?

**DevBuddy** (previously known as "Dad") is an open-source implementation of an amazing internal tool developed at
[Shopify](https://engineering.shopify.com) called "**Dev**".

The first goal of this tools is to automate the **setup** tasks required to work on a project.

With **DevBuddy**, pushing a change on a project you never touched look like this:

- `bud clone devbuddy/devbuddy`
- `bud up`
- `git commit`
- `bud test`
- `git push`

## Status and progress

**DevBuddy** is mostly useful for Python and Go projects. More languages will be natively
supported. Additional automatic tasks will also be implemented, making **DevBuddy** also
useful for languages without native support.

See the project config [documentation](docs/Config.md).

### Tasks:

Python:
- Python (version + virtualenv): **working**
- Pip (requirements file): **simple**
- Pipenv (support for [Pipfile](https://github.com/pypa/pipfile)): **simple**

Go:
- Go (version): **working**
- Dep (support for [Go Dep](https://github.com/golang/dep)): **simple**

Others
- Custom (conditional shell command): **working**
- Packages (Homebrew, Apt...): **planned**
- Docker Compose (manage a docker-compose setup): **planned**

### Automatic environment:

- Virtualenv: **working**
- Go: **working**

### Features:

- Notification when important files (eg: `requirements.txt`) are updated locally
  (eg: by `git pull`)
- A `help` command to guide a new developer based on `dev.yml`
- A `upgrade` command to auto-upgrade **DevBuddy**

### Hosting platform:

- Github: **working**
- Any git remote-url: **planned**
- Gitlab: **planned**
- Bitbucket (with Git): **planned**

### Shell integration

- Bash: **working**
- Zsh: **working**
- Fish: **planned**

## Install

```bash
$ bash -c "$(curl -sL https://raw.githubusercontent.com/devbuddy/devbuddy/master/install.sh)"
```

Uninstall it:
```bash
$ sudo rm /usr/local/bin/bud
```

## Setup

★ Install the shell integration (in `~/.bash_profile`):
```bash
eval "$(bud --shell-init --with-completion)"
```

A safer version:
```bash
type bud > /dev/null 2> /dev/null && eval "$(bud --shell-init --with-completion)"
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
See DevBuddy own [dev.yml](dev.yml)

```bash
$ bud
Usage:
  bud [flags]
  bud [command]

Available Commands:
  cd              Jump to a local project
  ci              Custom: Run all tests as CI would do
  clone           Clone a project from github.com
  create          Create a new project
  godoc           Custom: Starting GoDoc server on http://0.0.0.0:6060
  help            Help about any command
  inspect         Inspect the project and its tasks
  install-dev     Custom: Install bud in the GOROOT
  install-release Custom: Remove all bud binaries and install a released version
  integration     Custom: Run the integration tests
  lint            Custom: Lint the project
  lint-shell      Custom: Lint the shell scripts
  open            Open a link about your project
  release         Custom: Create a new release (bud release [VERSION])
  test            Custom: Run the unittests
  testup          Custom: Build, install and run the development version
  up              Ensure the project is up and running
  upgrade         [experimental] Upgrade DevBuddy to the latest available release.

Flags:
  -h, --help              help for bud
      --shell-init        Shell initialization
      --version           version for bud
      --with-completion   Enable completion during initialization

Use "bud [command] --help" for more information about a command.
```

## Development

★ Clone the repository:
```bash
~ $ mkdir -p ~/src/github.com/pior
~ $ cd ~/src/github.com/pior
~/src/github.com/pior $ git clone git@github.com:devbuddy/devbuddy.git
```

★ Fetch the dependencies (in `vendor`)
```bash
~/src/github.com/pior $ dep ensure
```

★ Install the go executable somewhere in your PATH:
```bash
~/src/github.com/pior $ go install
```

### Debugging

To show the debug messages, add this before in your environment:
```bash
export BUD_DEBUG=1
```
Or:
```bash
$ bud-enable-debug  # bud-disable-debug
```

### Release

```bash
$ bud release 1.0.0
```

The big idea:
1. Create a release commit to ensure the release is visible in the git log
2. Create an annotated tag
3. Push the commit and tag upstream

Expected:
1. The CI process will test the release
2. The CI process will publish macOS/Linux binaries to the Github Releases page

Updating the [install.sh](https://github.com/devbuddy/devbuddy/blob/master/install.sh) script is probably a good idea.

## License

[MIT](https://github.com/devbuddy/devbuddy/blob/master/LICENSE)
