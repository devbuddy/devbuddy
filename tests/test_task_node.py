import pytest


def test_env(cmd, project):
    project.write_devyml("""
        up:
        - node: '10.15.0'
    """)

    cmd.run("bud up")

    output = cmd.run("node -v")
    assert "v10.15.0" == output

