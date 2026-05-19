# Config: `dev.yml`

## Example

The file `dev.yml` must be placed at the root directory of your project.
It should be committed to your repository and shared with everyone working on the project.

The `up` section describes the tasks that the `bud up` command will run.
Some tasks prepare your environment to use a language like `python` or `go`.
Other tasks ensure your environment is up to date, like `pip` or `golang_dep`.
The special `custom` task lets you handle specific cases needed for your project. See [Tasks](#tasks).

The `commands` section describes the project commands like `bud test`. See [Project Commands](#project-commands).

The `open` section describes the project links available through `bud open <pattern>`. See [Open Command](#open-command).

**`dev.yml`**:
```yaml
up:
  - go:
      version: '1.12'
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

## Environment variables

Environment variables can be set by using the `env` key:

```yaml
env:
  ENV: development
  DATABASE_URL: mysql://localhost:3306/dev
  MEMCACHE_URL: localhost:11211
```

The environment variables will be set as soon as you enter the project. They can be used in `custom` tasks, in
the `commands` section and in your shell.

## Tasks

Some language tasks compile runtimes from source. On macOS, DevBuddy checks for
the Xcode command-line tools before starting those compilers. If they're missing,
DevBuddy runs `xcode-select --install`, then asks you to complete Apple's installer
dialog and re-run `bud up`.

DevBuddy also checks Xcode first-launch setup with `xcodebuild -checkFirstLaunchStatus`.
If setup is incomplete and `bud up` is running in an interactive terminal, DevBuddy runs
`sudo xcodebuild -runFirstLaunch`. That command installs required Xcode components and
accepts the Xcode/SDK license. In non-interactive sessions, DevBuddy prints the command
for you to run manually.

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

This task will install the Python version with pyenv, create a virtualenv, and
activate it in your shell. If pyenv is missing, DevBuddy installs it with Homebrew.
On macOS, DevBuddy checks for the Xcode command-line tools before invoking
`pyenv install`.

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

A Python environment must be selected first.

Currently this task cannot detect whether it should run. PR welcome!

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

A Python environment must be selected first.

Currently this task cannot detect whether it should run. PR welcome!

```yaml
up:
  - python: 3.6.5
  - pipfile
```

### `uv`

This task installs `uv` in the selected Python virtualenv if needed, then runs
`uv sync --active --inexact` so `uv` syncs the project into the DevBuddy-managed
virtualenv instead of creating a separate `.venv`.

A Python environment must be selected first. The project must have a
`pyproject.toml`; `uv.lock` is tracked when present.

```yaml
up:
  - python: 3.13.7
  - uv
```

Dependency groups, extras, and lockfile-only installs can be requested:

```yaml
up:
  - python: 3.13.7
  - uv:
      groups:
        - dev
      extras:
        - postgres
      frozen: true
```

Set `exact: true` to omit `--inexact` and let `uv sync` remove packages that are
not declared by the project.

### `go`

This task will download the Go distribution from `dl.google.com/go` and activate it
in your shell (with `GOROOT`).

```yaml
up:
  - go: 1.10.1
```

Go task with explicit version:
```yaml
up:
  - go:
      version: '1.12'
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
Optionally, it can also install the Node.js dependencies.

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

### `ruby`

This task will install the Ruby version (with rbenv, which is installed automatically via
Homebrew if missing) and activate it in your shell. If a `Gemfile` is present, gems will
be installed with `bundle install`. Both `Gemfile` and `Gemfile.lock` are watched, so
either changing will re-trigger the install.

The Ruby version can be set explicitly in `dev.yml`, or omitted to read it from a
`.ruby-version` file at the project root (an optional `ruby-` engine prefix is stripped):

```yaml
up:
  - ruby: 3.3.0
```

```yaml
up:
  - ruby  # reads .ruby-version
```

If both an explicit version in `dev.yml` and a `.ruby-version` file are present and
they disagree, DevBuddy warns on `bud up` and on shell activation, then proceeds with
the `dev.yml` version. Remove one to silence the warning.

#### Tradeoffs and limitations

- **macOS-first install path.** Bootstrapping rbenv is done via `brew install rbenv`.
  On Linux, install rbenv yourself before running `bud up`; the task will detect it
  and proceed.
- **Native build dependencies are partly guided.** `rbenv install` compiles Ruby from
  source and needs system headers (`libssl-dev`, `libyaml-dev`, `libffi-dev`, etc. on
  Linux; the Xcode command-line tools on macOS). DevBuddy handles the macOS check
  described above; Linux system packages still need to be handled through `apt:` or
  installed manually.
- **No shims, no `GEM_HOME`/`RUBYOPT`.** The autoenv feature simply prepends the
  selected version's `bin/` directory to `PATH`, which is enough for `ruby`, `gem`,
  `bundle`, and gem-installed executables. Tools that rely on rbenv shims or that
  read `GEM_HOME`/`GEM_PATH` directly will not see a version-specific environment.
- **Bundler is assumed to ship with Ruby.** Ruby 2.6 and later bundle Bundler in the
  stdlib, so the task does not install it explicitly. Older Ruby versions are not
  supported.
- **Gems install into the rbenv version directory** rather than a per-project
  `vendor/bundle`. If you want isolation, run `bundle config set --local path vendor/bundle`
  in the project before `bud up`.

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

The project can define custom commands in `dev.yml` that can be called with: `bud <command>`. Additional arguments are
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

`bud test` is not much shorter than calling `script/test`, for example.
The idea is to introduce an indirection that will be easy to document and remember by being consistent across projects
regardless of the programming language used (`rails test`? `pytest -v`? `npm test`? `go test ./...`?).

```bash
$ bud lint
🐼  running script/lint
pkg/project/current.go:14:2:warning: unused variable or constant someVariable declared but not used (varcheck)
```

## Open Command

The command `bud open <pattern>` opens a project link with the OS default handler (using `open`/`xdg-open`).

### Custom links

They are defined in `dev.yml`:
```yaml
open:
  staging: https://staging.myapp.com
  doc: https://godoc.org/github.com/org/myapp
```

Tip: `bud open` is enough if there is only one link.

Matching for custom links is fuzzy, so `bud open stg` can match `staging`.

### GitHub links:
- `github`/`gh`: open the GitHub source code page for your checked out branch
- `pullrequest`/`pr`: open the GitHub pull-request page for your checked out branch

Built-in GitHub link names are exact aliases only (no fuzzy matching).
