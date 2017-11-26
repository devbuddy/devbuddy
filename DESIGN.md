# Design

## Commands

### dad clone GIT-REMOTE-URL

- support for github.com git remote urls
- derive org name and project name from url

1. validate url
2. check if project is already cloned
3. create parent dir if needed
4. `git clone ...`
5. finalizer: `cd:/full/path`

### dad cd PROJECT-NAME

1. find existing project
2. finalizer: `cd:/full/path`

### dad up

1. find project root
2. find project manifest
3. read *setup tasks*
4. run tasks

### dad TASK

1. find project root
2. find project manifest
3. read *tasks*
4. run the task command

### dad new PROJECT-NAME

1. ask for project parent folder `~/src/github.com/pior` (with default on project created/cloned)
2. create directory
3. finalizer: `cd:/full/path`


## Prompt hook

Features:
- detect when entering/leaving a project path
- activate/deactivate virtualenv
- notify about project manifest changes ("you didn't `dad up` for a long time")
- notify about things defined by *setup tasks* (like updates of `requirements.txt`)


## Config

Defaults:
```bash
export DAD_SOURCE_DIR=~/src
```


## Installation

Initialize the shell:

```bash
# DAD_DEBUG=1
eval "$(dad --init)"
```
