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
    project_path = os.path.abspath(os.path.join(os.path.dirname(__file__), os.pardir))

    proc = subprocess.run(
        'go build -o {}/bud'.format(binary_path),
        shell=True,
        cwd=project_path,
        stdout=subprocess.PIPE,
        stderr=subprocess.PIPE,
    )
    if proc.returncode:
        raise RuntimeError(
            "Failed to build binary to use in tests:\nstdout: %s\nstderr: %s" % (proc.stdout, proc.stderr)
        )


@pytest.fixture(scope='module')
def workdir(tmpdir_factory):
    return tmpdir_factory.mktemp("poipoidir")


def build_pexpect_bash(cwd):
    child = pexpect.spawn('bash', ['--norc', '--noprofile'], echo=False, encoding='utf-8', cwd=cwd)

    # If the user runs 'env', the value of PS1 will be in the output. To avoid
    # replwrap seeing that as the next prompt, we'll embed the marker characters
    # for invisible characters in the prompt; these show up when inspecting the
    # environment variable, but not when bash displays the prompt.
    ps1 = pexpect.replwrap.PEXPECT_PROMPT[:5] + u'\\[\\]' + pexpect.replwrap.PEXPECT_PROMPT[5:]
    ps2 = pexpect.replwrap.PEXPECT_CONTINUATION_PROMPT[:5] + u'\\[\\]' + pexpect.replwrap.PEXPECT_CONTINUATION_PROMPT[5:]
    prompt_change = u"PS1='{0}' PS2='{1}' PROMPT_COMMAND=''".format(ps1, ps2)

    return pexpect.replwrap.REPLWrapper(child, u'\\$', prompt_change, extra_init_cmd="export PAGER=cat")


def build_pexpect_zsh(cwd):
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
        extra_init_cmd="export PAGER=cat; export SHELL=/bin/zsh",
    )


PEXPECT_SHELLS = {'bash': build_pexpect_bash, 'zsh': build_pexpect_zsh}


@pytest.fixture(scope='module', params=['bash', 'zsh'])
def shell(workdir, binary_path, request):
    build_pexpect_shell = PEXPECT_SHELLS[request.param]
    shell_ = build_pexpect_shell(workdir)

    shell_.run_command('export PATH={}:$PATH'.format(binary_path))

    output = shell_.run_command('which bud')
    assert str(binary_path) in output

    shell_.run_command('eval "$(bud --shell-init)"')

    output = shell_.run_command('type bud')
    assert 'bud is a shell function' in output or 'bud is a function' in output

    output = shell_.run_command('bud --version')
    assert 'bud version devel' in output

    return shell_


@pytest.fixture
def make_test_repo(request):
    def func(name):
        path = os.path.expanduser('~/src/github.com/%s' % name)

        if os.path.exists(path):
            shutil.rmtree(path)

        os.makedirs(path)

        def cleanup():
            shutil.rmtree(path)
        request.addfinalizer(cleanup)

        return path

    return func
