import os
import shutil
import subprocess

import pexpect.replwrap
import pytest


@pytest.fixture(scope='session')
def binary_path(tmpdir_factory):
    return tmpdir_factory.mktemp("bin")


@pytest.fixture(scope='session', autouse=True)
def binary(binary_path):
    subprocess.run('go build -o {}/bud'.format(binary_path), shell=True, check=True)


@pytest.fixture(scope='module')
def workdir(tmpdir_factory):
    return tmpdir_factory.mktemp("poipoidir")


def pexpect_bash(cwd):
    env = os.environ.copy()
    child = pexpect.spawn('bash', ['--norc', '--noprofile'], echo=False, encoding='utf-8', env=env, cwd=cwd)
    p = pexpect.replwrap.REPLWrapper(child, u'', None, extra_init_cmd="export PAGER=cat")
    return p


def pexpect_zsh(cwd):
    child = pexpect.spawn(
        'zsh', ['--no-globalrcs', '--no-rcs', '--no-zle', '--no-promptcr'],
        echo=False,
        encoding='utf-8',
        env={'PROMPT': 'ps1'},
        cwd=cwd,
    )

    return pexpect.replwrap.REPLWrapper(
        child,
        u'ps1',
        prompt_change=u"PROMPT='{}'",
        extra_init_cmd="unset zle_bracketed_paste; export PAGER=cat",
    )


@pytest.fixture(scope='module')
def shell(workdir, binary_path):
    name = 'zsh'

    if name == 'bash':
        p = pexpect_bash(workdir)
    elif name == 'zsh':
        p = pexpect_zsh(workdir)
        p.run_command('export SHELL=/bin/zsh')
    else:
        raise ValueError('unknown shell: %s' % name)

    p.run_command('export PATH={}:$PATH'.format(binary_path))

    output = p.run_command('which bud')
    assert str(binary_path) in output

    p.run_command('eval "$(bud --shell-init)"')

    output = p.run_command('type bud')
    assert 'bud is a shell function' in output

    output = p.run_command('bud --version')
    assert 'bud version devel' in output

    return p


@pytest.fixture
def make_test_repo(request):
    def func(name):
        path = os.path.expanduser('~/src/github.com/devbuddy_integration_tests/%s' % name)

        if os.path.exists(path):
            shutil.rmtree(path)

        os.makedirs(path)

        def cleanup():
            shutil.rmtree(path)
        request.addfinalizer(cleanup)

        return path

    return func
