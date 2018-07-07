# Contibuting to DevBuddy

## Quickstart

Install DevBuddy following the steps in the [README](README.md#install).

Clone the project:
```shell
$ bud clone devbuddy/devbuddy
Cloning into '/Users/pior/src/github.com/devbuddy/devbuddy'...
🐼  jumping to github.com:devbuddy/devbuddy
🐼  golang activated. (version: 1.10.1)
🐼  failed to activate python. Try running 'bud up' first! (version: 3.6.5)

~/src/github.com/devbuddy/devbuddy (master) $ bud up
◼︎ Golang (1.10.1)
◼︎ Go Dep (dep ensure)
  ▪︎Run dep ensure
  Running: dep ensure
◼︎ Python (3.6.5)
  ▪︎create virtualenv
  Running: /Users/pior/.pyenv/versions/3.6.5/bin/virtualenv /Users/pior/.local/share/bud/virtualenvs/devbuddy-1911165133-3.6.5
  Using base prefix '/Users/pior/.pyenv/versions/3.6.5'
  New python executable in /Users/pior/.local/share/bud/virtualenvs/devbuddy-1911165133-3.6.5/bin/python3.6
  Also creating executable in /Users/pior/.local/share/bud/virtualenvs/devbuddy-1911165133-3.6.5/bin/python
  Installing setuptools, pip, wheel...done.
◼︎ Pip (tests/requirements.txt)
  ▪︎install tests/requirements.txt
  Running: pip install --require-virtualenv -r tests/requirements.txt
  Collecting pytest==3.6.1 (from -r tests/requirements.txt (line 1))
    Using cached https://files.pythonhosted.org/packages/d3/75/e79b66c9fe6166a90004bb8fb02bab06213c3348e93f3be41d7eaf625554/pytest-3.6.1-py2.py3-none-any.whl
[...]
  Collecting ptyprocess>=0.5 (from pexpect==4.6.0->-r tests/requirements.txt (line 2))
    Downloading https://files.pythonhosted.org/packages/d1/29/605c2cc68a9992d18dada28206eeada56ea4bd07a239669da41674648b6f/ptyprocess-0.6.0-py2.py3-none-any.whl
  Installing collected packages: six, more-itertools, attrs, atomicwrites, pluggy, py, pytest, ptyprocess, pexpect
  Successfully installed atomicwrites-1.1.5 attrs-18.1.0 more-itertools-4.2.0 pexpect-4.6.0 pluggy-0.6.0 ptyprocess-0.6.0 py-1.5.4 pytest-3.6.1 six-1.11.0
◼︎ Homebrew (curl)
◼︎ Custom (Install shellcheck)
◼︎ Custom (Install gometalinter)
🐼  python activated. (version: 3.6.5)
```

### Run the tests

```shell
$ bud test
🐼  running script/test
?   	github.com/devbuddy/devbuddy	[no test files]
?   	github.com/devbuddy/devbuddy/pkg/cmd	[no test files]
[...]
?   	github.com/devbuddy/devbuddy/pkg/test	[no test files]
ok  	github.com/devbuddy/devbuddy/pkg/utils	(cached)
```

### Run the linters

```shell
$ bud lint
🐼  running script/lint
```

### Run the integration tests

```shell
$ bud integration
🐼  running pytest -v tests
============================= test session starts ==============================
platform darwin -- Python 3.6.5, pytest-3.6.1, py-1.5.4, pluggy-0.6.0 -- /Users/pior/.local/share/bud/virtualenvs/devbuddy-1911165133-3.6.5/bin/python3.6
cachedir: .pytest_cache
rootdir: /Users/pior/src/github.com/devbuddy/devbuddy, inifile:
collecting ... collected 26 items

tests/test_cd.py::test_find_project PASSED                               [  3%]
tests/test_cd.py::test_ui PASSED                                         [  7%]
tests/test_command.py::test_option_version PASSED                        [ 11%]
[...]
tests/test_up.py::test_invalid_manifest_with_string PASSED               [ 92%]
tests/test_up.py::test_unknown_task PASSED                               [ 96%]
tests/test_up.py::test_invalid_task PASSED                               [100%]

========================== 26 passed in 28.53 seconds ==========================
```
