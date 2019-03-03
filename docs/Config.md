# Config: `dev.yml`

## Example

The file `dev.yml` must be placed at the root directory of your project.
It should be commited to your repository and shared with everyone working on the project.

The `up` section describes the tasks that the `bud up` command will run.
Some tasks prepare your environment to use a language like `python` or `go`.
Other tasks ensure your environment is up to date like `pip` or `golang_dep`.
The special `custom` task let you handle specific case needed for your project. See [Tasks](#tasks).

The `commands` section describes the project commands like `bud test`. See [Project Commands](#project-commands).

The `open` section describes the project links available through `bud open <name>`. See [Open Command](#open-command).

**`dev.yml`**:
```yaml
up:
  - go:
      version: '1.12'
      modules: true
  - golang_dep
  - python: 3.6.5
  - apt: [git, curl]
  - homebrew: [git, curl]
  - pip:
      - python/requirements-dev.txt
  - custom:
      name: Download GeoIP db
      met?: test -e GeoIP.dat
      meet: curl -L http://geolite.maxmind.com/download/geoip/database/GeoLiteCountry/GeoIP.dat.gz | gunzip > GeoIP.dat

commands:
  test:
    desc: Run tests for Python and Go
    run: pytest -v python/tests && go test $(go list ./...)

open:
  staging: https://staging.myapp.com
  doc: https://godoc.org/github.com/org/myapp
```

## Tasks

### `apt`

This task will install the Debian packages specified if DevBuddy is running on Debian.
Otherwise the task will be ignored.

```yaml
up:
  - apt:
    - python3-dev
```

### `homebrew`

This task will install the Homebrew recipes if DevBuddy is running on macOS.
Otherwise the task will be ignored.

```yaml
up:
  - homebrew:
    - cmake
```

### `python`

This task will install the Python version (with PyEnv, which must be installed), create
a virtualenv and activate it in your shell.

```yaml
up:
  - python: 3.6.5
```

### `python_develop`

This task will install the Python project in development mode (`pip install -e .`).

This task runs if `setup.py` has changed since the last `bud up`.


```yaml
up:
  - python: 3.6.5
  - python_develop
```

### `pip`

This task will install a pip requirements file.

A Python environment must be selected before.

Currently this task can't detect whether it should run or not. PR welcome!

```yaml
up:
  - python: 3.6.5
  - pip:
    - requirements.txt
    - requirements-dev.txt
```

### `pipenv`

This task will install a [Pipfile](https://github.com/pypa/pipfile) with
[Pipenv](https://github.com/pypa/pipenv).

A Python environment must be selected before.

Currently this task can't detect whether it should run or not. PR welcome!

```yaml
up:
  - python: 3.6.5
  - pipfile
```

### `go`

This task will download the Go distribution from `dl.google.com/go` and activate it
in your shell (with `GOROOT`).

```yaml
up:
  - go: 1.10.1
```

Force the usage of Go modules:
```yaml
up:
  - go:
      version: '1.12'
      modules: true
```

### `golang_dep`

This task will run [Go Dep](https://github.com/golang/dep) if needed.

A Go environment must be selected before.

```yaml
up:
  - go: 1.10.1
  - golang_dep
```

### `node`

This task will download the Node distribution from `nodejs.org` and activate it in your shell.
Optionally, it can also install the NodeJS dependencies.

```yaml
up:
  - node: 10.15.0
```

To install the dependencies with NPM:

```yaml
up:
  - node: 10.15.0
      npm: true
```

### `custom`

This task will run a command if a condition is not met.
The condition is expressed as a command.

```yaml
up:
  - custom:
      name: Install shellcheck with Brew
      met?: test -e /usr/local/Cellar/shellcheck
      meet: brew install shellcheck
```

## Project commands

The project can define custom command in `dev.yml` that can be called with: `bud <command>`. Additional arguments are
also passed to the command: `bud <command> <arg> <arg>...`.

```yaml
commands:
  test:
    desc: Run tests for Python and Go
    run: pytest -v python/tests && go test $(go list ./...)
  lint:
    desc: Run the linters
    run: script/run_all_linters
```

`bud test` is not much shorter than calling `script/test` for example.
The idea is to introduce an indirection that will be easy to document and remember by being consistent across projects
regardless of the programming language used (`rails test`? `pytest -v`? `npm test`? `go test ./...`?).

```bash
$ bud lint
🐼  running script/lint
pkg/project/current.go:14:2:warning: unused variable or constant someVariable declared but not used (varcheck)
```

## Open Command

The command `bud open <name>` will open a link about the project with the OS default handler (using `open`/`xdg-open`).

### Custom links

They are defined in `dev.yml`:
```yaml
open:
  staging: https://staging.myapp.com
  doc: https://godoc.org/github.com/org/myapp
```

Tip: `dev open` is enough if there is only one link.

### Github links:
- `github`/`gh`: open the Github source code page for your checked out branch
- `pullrequest`/`pr`: open the Github pull-request page for your checked out branch
