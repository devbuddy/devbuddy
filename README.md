# Dad (prototype) - Open implementation of Shopify Dev

## Current state

Simple prototype which purpose is to help me define the minimal set of features
needed to build a useful tool.

## Install

Clone the repository:
```bash
$ git clone git@github.com:pior/dad-proto.git ~/.dad
```

Source `dad.sh` in your profile (`~/.bash_profile`):
```bash
[[ -f ~/.dad/dad.sh ]] && [[ $- == *i* ]] && source ~/.dad/dad.sh
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
