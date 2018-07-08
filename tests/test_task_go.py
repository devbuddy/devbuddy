import pytest


@pytest.fixture
def gopath(tmpdir_factory):
    return tmpdir_factory.mktemp("gopath")


def test_env(cmd, project, gopath):
    cmd.run(f"export GOPATH={gopath}")

    project.write_devyml("""
        up:
        - go: '1.5'
    """)

    cmd.run("bud up")
    cmd.assert_succeed()

    output = cmd.run("go version")
    assert "go version go1.5" in output


def test_warn_gopath_missing(cmd, project, gopath):
    cmd.run("unset GOPATH")

    project.write_devyml("""
        up:
        - go: '1.5'
    """)

    output = cmd.run("bud up")
    cmd.assert_failed()

    assert "The GOPATH environment variable should be set" in output
