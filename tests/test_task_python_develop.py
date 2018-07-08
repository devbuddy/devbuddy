import os
import textwrap


def make_setuppy(version):
    return textwrap.dedent(f"""
        from setuptools import setup, find_packages
        setup(name='devbuddy-test-pkg', version='{version}')

        open("sentinel", "w").write("")
    """)


def test_with_modification(cmd, project):
    project.write_devyml("""
        up:
        - python: 3.6.5
        - python_develop
    """)

    project.write_file("setup.py", make_setuppy(version=42))

    cmd.run("bud up")
    cmd.assert_succeed()

    output = cmd.run("pip show devbuddy-test-pkg")
    assert "Version: 42" in output

    project.write_file("setup.py", make_setuppy(version=84))

    cmd.run("bud up")
    cmd.assert_succeed()

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
    cmd.assert_succeed()

    assert os.path.exists(sentinel_path)

    os.unlink(sentinel_path)

    cmd.assert_succeed()
    cmd.run("bud up")

    assert not os.path.exists(sentinel_path)
