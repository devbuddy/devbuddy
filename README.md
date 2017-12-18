# Dad 👴

[![Go Report Card](https://goreportcard.com/badge/github.com/pior/dad)](https://goreportcard.com/report/github.com/pior/dad)
[![Build Status](https://travis-ci.org/pior/dad.svg?branch=master)](https://travis-ci.org/pior/dad)

## What is this?

**Dad** is an open-source implementation of an amazing internal tool developed at
[Shopify](https://engineering.shopify.com) called "**Dev**".

The first goal of this tools is to automate the **setup** tasks required to work on a project.

With **Dad**, pushing a change on a project you never touched look like this:

- `dad clone Shopify/sarama`
- `dad up`
- **Commit your changes**
- `dad test`
- `git push`

## Scope and Status

**Dad** will start as a highly compatible implementation of **Dev** (when possible).

**Dad** is currently a _working prototype_. Core features are being implemented.

Current limitations:

- Github is the only _hosting platform_ implemented
- Bash is the only shell supported

## Install

Binaries releases will be available from Github soon.

In the mean time, you will need to compile/install the Go binary yourself:

★ Clone the repository:
```bash
~ $ mkdir -p ~/src/github.com/pior
~ $ cd ~/src/github.com/pior
~/src/github.com/pior $ git clone git@github.com:pior/dad.git
```

★ Fetch the dependencies (in `vendor`)
```bash
~/src/github.com/pior $ dep ensure
```

★ Install the go executable somewhere in your PATH:
```bash
~/src/github.com/pior $ make install
```

## Setup

★ Install the shell integration (in `~/.bash_profile`):
```bash
eval "$(dad --shell-init --with-completion)"
```

A safer version:
```bash
type dad > /dev/null 2> /dev/null && eval "$(dad --shell-init --with-completion)"
```

## Usage

★ Add a `dev.yml` file in your project:
```yaml
up:
  - custom:
      met?: brew info upx 2> /dev/null > /dev/null
      meet: brew install upx

commands:
  test:
    run:
      make test

  lint:
    run: make lint
```
See Dad own [dev.yml](dev.yml)

```bash
$ dad
Usage:
  dad [flags]
  dad [command]

Available Commands:
  cd          Jump to a local project
  clone       Clone a project from github.com
  help        Help about any command
  lint        Custom
  test        Custom: Run tests
  up          Ensure the project is up and running

Flags:
  -h, --help              help for dad
      --shell-init        Shell initialization
      --with-completion   Enable completion during initialization

Use "dad [command] --help" for more information about a command.
```

## Development

### Release

1. Push the release code to the `master` branch
2. Tag the commit: `git tag v1.2.3 && git push --tags`
3. Admire the release on [Github](https://github.com/pior/dad/releases)
