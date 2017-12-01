# Dad - Open implementation of Shopify Dev

## Install

Releases will be available from Github in the future. For the initial development
period, the Go compilation will

Clone the repository:
```bash
~ $ mkdir -p ~/src/github.com/pior

~ $ cd ~/src/github.com/pior

~/src/github.com/pior $ git clone git@github.com:pior/dad.git
...
```

Install the go executable somewhere in your PATH:
```bash
~/src/github.com/pior $ make install
```

Install the shell integration (in `~/.bash_profile`):
```bash
type dad > /dev/null 2> /dev/null && [[ $- == *i* ]] && eval "$(dad --shell-init)"
```

## Usage

```bash
$ dad
Usage: dad [command] ...

Commands:

  new NAME          Create a new project

  cd NAME           Go to a project and activate environment

  up                Prepare your development environment

  test              Run the test suite (alias: t)

  server            Run the server
```
