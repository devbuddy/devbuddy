# Contributing to DevBuddy

## Quickstart

Install DevBuddy by following the steps in the [README](../README.md#install).

### Clone the project

This will clone a repository from Github: `github.com/devbuddy/devbuddy`.

The repo will be cloned at `~/src/github.com/devbuddy/devbuddy`.

```shell
$ bud clone devbuddy/devbuddy
```

### Setup the project with DevBuddy

This is the core feature of DevBuddy, the `up` command is supposed to prepare/install/setup everything needed to
start developing on the project.

```shell
~/src/github.com/devbuddy/devbuddy $ bud up
```

The command will sequentially evaluate the *up* tasks defined in [dev.yml](../dev.yml).
Some will setup the working environment (`go`, `python`), some will install dependencies (`golang_dep`, `pip`),
some will conditionally execute an arbitrary command for specific situations.

### Run the tests

Project specific commands can be defined in the `commands` section of the [dev.yml](../dev.yml).

Typical commands are `test`, `lint`, `clean`, `release`

```shell
~/src/github.com/devbuddy/devbuddy $ bud test
```

### Run the linters

```shell
~/src/github.com/devbuddy/devbuddy $ bud lint
```

### Run the integration tests

```shell
~/src/github.com/devbuddy/devbuddy $ bud integration
```

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

You can enable the debug messages with:

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
The CI process will build and upload the distributions on Github.

Updating the version defined in the [install.sh](https://github.com/devbuddy/devbuddy/blob/master/install.sh)
script is probably a good idea.
