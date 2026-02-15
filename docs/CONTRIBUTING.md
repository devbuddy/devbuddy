# Contributing to DevBuddy

## Quickstart

Install DevBuddy by following the steps in the [README](../README.md#install).

### Clone the project

This will clone a repository from GitHub: `github.com/devbuddy/devbuddy`.

The repo will be cloned at `~/src/github.com/devbuddy/devbuddy`.

```shell
$ bud clone devbuddy/devbuddy
```

### Set up the project with DevBuddy

This is the core feature of DevBuddy: the `up` command prepares everything needed to start developing on the project.

```shell
~/src/github.com/devbuddy/devbuddy $ bud up
```

The command will sequentially evaluate the *up* tasks defined in [dev.yml](../dev.yml).
Some will set up the working environment (`go`, `python`), some will install dependencies (`golang_dep`, `pip`),
some will conditionally execute an arbitrary command for specific situations.

### Run the tests

Project specific commands can be defined in the `commands` section of the [dev.yml](../dev.yml).

Typical commands are `test`, `lint`, `clean`, and `release`.

```shell
~/src/github.com/devbuddy/devbuddy $ bud test
```

### Run the linters

```shell
~/src/github.com/devbuddy/devbuddy $ bud lint
```

### Run the integration tests

```shell
~/src/github.com/devbuddy/devbuddy $ TEST_SHELL=bash go test -v -count=1 ./tests
~/src/github.com/devbuddy/devbuddy $ TEST_SHELL=zsh go test -v -count=1 ./tests
```

Make sure Docker is running, and set `TEST_DOCKER_IMAGE` when needed.

### Install DevBuddy from your branch

```shell
~/src/github.com/devbuddy/devbuddy $ bud install-dev
```

### Reinstall a release

```shell
~/src/github.com/devbuddy/devbuddy $ bud install-release
```

Or simply:

```shell
~ $ bud upgrade
```

### Debugging

You can enable debug logging with:

```bash
$ bud-enable-debug
```

Which is equivalent to:

```bash
export BUD_DEBUG=1
```

## Release

Create a release locally with the command:
```bash
$ bud release
```

This command will create a tag and push it to the origin.
The CI process will build and upload the distributions on GitHub.

If release tooling changes, keep install instructions in [README](../README.md#install) in sync.
