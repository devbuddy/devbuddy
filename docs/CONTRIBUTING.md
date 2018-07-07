# Contributing to DevBuddy

## Quickstart

Install DevBuddy following the steps in the [README](README.md#install).

### Clone the project

```shell
$ bud clone devbuddy/devbuddy
```

```shell
~/src/github.com/devbuddy/devbuddy $ bud up
```

### Run the tests

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
$ bud release 1.0.0
```

This command will create a tag and push it to the origin.
The CI process will build and upload the distributions on Github.

Updating the version defined in the [install.sh](https://github.com/devbuddy/devbuddy/blob/master/install.sh)
script is probably a good idea.
