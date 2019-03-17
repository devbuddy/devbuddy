import os
import random
import shutil
import subprocess
import textwrap

import pexpect.replwrap
import pytest


def pytest_addoption(parser):
    parser.addoption("--shell", default="bash", help="Shell to use to test DevBuddy", choices=['bash', 'zsh'])


@pytest.fixture(scope='session')
def workdir(tmpdir_factory):
    return tmpdir_factory.mktemp("workdir")


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


def build_pexpect_bash(workdir):
    env = {
        'PATH': os.getenv('PATH'),
        'LANG': 'C.UTF-8',
        'LC_ALL': 'C.UTF-8',
    }

    child = pexpect.spawn(
        'bash', ['--norc', '--noprofile'],
        echo=False,
        encoding='utf-8',
        env=env,
        cwd=str(workdir),
    )
    child.timeout = 180

    # If the user runs 'env', the value of PS1 will be in the output. To avoid
    # replwrap seeing that as the next prompt, we'll embed the marker characters
    # for invisible characters in the prompt; these show up when inspecting the
    # environment variable, but not when bash displays the prompt.
    ps1 = pexpect.replwrap.PEXPECT_PROMPT[:5] + u'\\[\\]' + pexpect.replwrap.PEXPECT_PROMPT[5:]
    ps2 = pexpect.replwrap.PEXPECT_CONTINUATION_PROMPT[:5] + u'\\[\\]' + pexpect.replwrap.PEXPECT_CONTINUATION_PROMPT[5:]
    prompt_change = u"PS1='{0}' PS2='{1}' PROMPT_COMMAND=''".format(ps1, ps2)

    return pexpect.replwrap.REPLWrapper(child, u'\\$', prompt_change, extra_init_cmd="export PAGER=cat")


def build_pexpect_zsh(workdir):
    env = {
        'PROMPT': 'ps1',
        'PATH': os.getenv('PATH'),
        'LANG': 'C.UTF-8',
        'LC_ALL': 'C.UTF-8',
    }

    child = pexpect.spawn(
        'zsh', ['--no-globalrcs', '--no-rcs', '--no-zle', '--no-promptcr'],
        echo=False,
        encoding='utf-8',
        env=env,
        cwd=str(workdir),
    )
    child.timeout = 180

    return pexpect.replwrap.REPLWrapper(
        child,
        u'ps1',
        prompt_change=u"PROMPT='{}'",
        extra_init_cmd="export PAGER=cat; export SHELL=/bin/zsh",
    )


PEXPECT_SHELLS = {'bash': build_pexpect_bash, 'zsh': build_pexpect_zsh}


class CommandTestHelper:
    def __init__(self, pexpect_wrapper):
        self._pexpect_wrapper = pexpect_wrapper

    def _run(self, command):
        return self._pexpect_wrapper.run_command(command).strip()

    def run(self, command, expect_exit_code=0, check_activation_failure=True):
        output = self._run(command)

        error_lines = [l for l in output.splitlines() if 'failed to activate ' in l]
        if check_activation_failure and error_lines:
            pytest.fail(f"Failed to activate features:\n{error_lines}")

        code = self.get_exit_code()
        if code != expect_exit_code:
            pytest.fail(f"Command failed with code {code}. Output:\n{output}")

        return output

    def get_exit_code(self):
        lines = self._run("echo $?").splitlines()
        return int(lines[0])  # Ignore lines produced by prompt hook


@pytest.fixture(scope='session')
def cmd(binary_path, workdir, request):
    shell_name = request.config.getoption("--shell")

    build_pexpect_shell = PEXPECT_SHELLS[shell_name]
    pexpect_wrapper = build_pexpect_shell(workdir)

    pexpect_wrapper.run_command('export PATH={}:$PATH'.format(binary_path))

    output = pexpect_wrapper.run_command('which bud')
    assert str(binary_path) in output

    pexpect_wrapper.run_command('eval "$(bud --shell-init)"')

    output = pexpect_wrapper.run_command('type bud')
    assert 'bud is a shell function' in output or 'bud is a function' in output

    output = pexpect_wrapper.run_command('bud --version')
    assert 'bud version devel' in output

    return CommandTestHelper(pexpect_wrapper)


class ProjectTestHelper:
    def __init__(self, org, name, path):
        self.org = org
        self.name = name
        self.path = os.path.expanduser('~/src/github.com/%s/%s' % (org, name))

    def write_devyml(self, body):
        self.write_file('dev.yml', textwrap.dedent(body))

    def write_file(self, local_path, data):
        with open(os.path.join(self.path, local_path), 'w') as fp:
            fp.write(data)

    def write_file_dedent(self, local_path, data):
        self.write_file(local_path, textwrap.dedent(data).lstrip())

    def assert_file(self, local_path, expect_content=None, present=True):
        path = os.path.join(self.path, local_path)

        exists = os.path.exists(path)

        if present:
            assert exists, f"file \"{local_path}\" should exist"
        else:
            assert not exists, f"file \"{local_path}\" should not exist"

        if expect_content is not None:
            with open(path, 'r') as fp:
                content = fp.read()
            assert content == expect_content, f"file content not as expected for \"{local_path}\""


@pytest.fixture
def srcdir():
    return os.path.expanduser('~')


@pytest.fixture
def project_factory(request, srcdir):
    def func(org, name):
        path = srcdir.join('src').join('github.com').join(org).join(name)
        project = ProjectTestHelper(org, name, path)

        if os.path.exists(project.path):
            shutil.rmtree(project.path)
        os.makedirs(project.path)

        def cleanup():
            shutil.rmtree(project.path)
        request.addfinalizer(cleanup)

        return project

    return func


@pytest.fixture
def project(cmd, project_factory):
    p = project_factory('devbuddy_tests', 'test_%s' % random.randrange(10e6))
    cmd.run(f"cd {p.path}")
    return p
