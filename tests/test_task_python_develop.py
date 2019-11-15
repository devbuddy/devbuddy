import os
import textwrap


def make_setuppy(version=1):
    return textwrap.dedent("""
        from setuptools import setup, find_packages
        setup(name='devbuddy-test-pkg', version='%s', extras_require={'test': ['pyreleaser==0.5.2']})

        open("sentinel", "w").write("")
    """) % version


def test_with_modification(cmd, project):
    project.write_devyml("""
        up:
        - python: 3.6.5
        - python_develop
    """)

    project.write_file("setup.py", make_setuppy(version=42))

    cmd.run("bud up")

    output = cmd.run("pip show devbuddy-test-pkg")
    assert "Version: 42" in output

    project.write_file("setup.py", make_setuppy(version=84))

    cmd.run("bud up")

    output = cmd.run("pip show devbuddy-test-pkg")
    assert "Version: 84" in output


def test_without_modification(cmd, project):
    project.write_devyml("""
        up:
        - python: 3.6.5
        - python_develop
    """)

    sentinel_path = os.path.join(project.path, "sentinel")

    project.write_file("setup.py", make_setuppy(version=42))

    cmd.run("bud up")

    assert os.path.exists(sentinel_path)

    os.unlink(sentinel_path)

    cmd.run("bud up")

    assert not os.path.exists(sentinel_path)


def test_without_extra(cmd, project):
    project.write_devyml("""
        up:
        - python: 3.6.5
        - python_develop
    """)
    project.write_file("setup.py", make_setuppy())

    cmd.run("bud up")

    output = cmd.run("pip freeze")
    assert 'pyreleaser==0.5.2' not in output.splitlines(False)


def test_with_extra(cmd, project):
    project.write_devyml("""
        up:
        - python: 3.6.5
        - python_develop:
            extras: [test]
    """)
    project.write_file("setup.py", make_setuppy())

    cmd.run("bud up")

    output = cmd.run("pip freeze")
    assert 'pyreleaser==0.5.2' in output.splitlines(False)


def test_with_unknown_extra(cmd, project):
    # Unknown extra are ignored: https://github.com/devbuddy/devbuddy/issues/229

    project.write_devyml("""
        up:
        - python: 3.6.5
        - python_develop:
            extras: [nope]
    """)
    project.write_file("setup.py", make_setuppy())

    cmd.run("bud up")
