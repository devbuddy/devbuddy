# Dad - Open implementation of Shopify Dev

## Install

Binaries releases will be available from Github in the future.
During the alpha development period, users will need to compile the Go binary
themselves.

★ Clone the repository:
```bash
~ $ mkdir -p ~/src/github.com/pior && cd ~/src/github.com/pior
~/src/github.com/pior $ git clone git@github.com:pior/dad.git
...
```

★ Fetch the dependencies in `vendor`
```bash
dep ensure
```

★ Install the go executable somewhere in your PATH:
```bash
~/src/github.com/pior $ make install
```

★ Install the shell integration (in `~/.bash_profile`):
```bash
type dad > /dev/null 2> /dev/null && eval "$(dad --shell-init --with-completion)"
```

## Usage

```bash
$ dad
Usage:
  dad [flags]
  dad [command]

Available Commands:
  cd          Jump to a local project
  clone       Clone a project from github.com
  help        Help about any command

Flags:
  -h, --help              help for dad
      --shell-init        Shell initialization
      --with-completion   Enable completion during initialization

Use "dad [command] --help" for more information about a command.
```
