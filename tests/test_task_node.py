import pytest


def test_env(cmd, project):
    project.write_devyml("""
        up:
        - node: '10.15.0'
    """)

    cmd.run("bud up")

    output = cmd.run("node -v")
    assert "v10.15.0" == output


def test_npm_install(cmd, project):
    project.write_devyml("""
        up:
        - node: '10.15.0'
    """)

    project.write_file("package.json", """
        {
            "dependencies": {
                "testpackage": "1.0.0"
            }
        }
    """)

    cmd.run("bud up")

    output = cmd.run("npm list")
    assert "testpackage@1.0.0" in output

