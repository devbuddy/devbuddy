# Config: `dev.yml`

## Example

The file `dev.yml` must be placed at the root directory of your project.
It should be commited to your repository and shared with everyone working on the project.

The `up` section describes the tasks that the `dev up` command will run.

The `commands` section describes the custom commands like `dev test`.

Example:
```yaml
up:
  - go: 1.10.1
  - golang_dep
  - python: 3.6.5
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
```

## Tasks for Python

### `python`

This task will install the Python version (with PyEnv, which must be installed), create
a virtualenv and activate it in your shell.

```yaml
up:
  - python: 3.6.5
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

## Tasks for Go

## `go`

This task will download the Go distribution from `dl.google.com/go` and activate it
in your shell (with `GOROOT`).

```yaml
up:
  - go: 1.10.1
```

## `golang_dep`

This task will run [Go Dep](https://github.com/golang/dep) if needed.

A Go environment must be selected before.

```yaml
up:
  - go: 1.10.1
  - golang_dep
```

## Custom task

This task will run a command if a condition is not met.
The condition is expressed as a command.

```yaml
up:
  - custom:
      desc: Install shellcheck with Brew
      met?: test -e /usr/local/Cellar/shellcheck
      meet: brew install shellcheck
```
