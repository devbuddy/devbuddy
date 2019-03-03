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

    output = cmd.run("go version")
    assert "go version go1.5" in output


def test_warn_gopath_missing(cmd, project, gopath):
    cmd.run("unset GOPATH")

    project.write_devyml("""
        up:
        - go: '1.5'
    """)

    output = cmd.run("bud up", expect_exit_code=1)
    assert "The GOPATH environment variable should be set" in output


def test_with_modules(cmd, project, gopath):
    # We want to support pre-modules and modules projects in the same environment
    # so we set a GOPATH as it would be for pre-modules setup
    # Devbuddy will set GO111MODULES=on to force-enable Go modules even if we are in the GOPATH
    cmd.run(f"export GOPATH={gopath}")

    project.write_devyml("""
        up:
        - go:
            version: '1.12'
            modules: true
    """)

    output = cmd.run("bud up")

    project.write_file_dedent("main.go", """
        package main

        import (
            "fmt"
            "github.com/spf13/pflag"
        )

        func main() {
            pflag.Parse()
            fmt.Println(pflag.Arg(0))
        }
    """)

    project.write_file_dedent("go.mod", """
        module poipoi

        require github.com/spf13/pflag v1.0.3
    """)

    project.write_file_dedent("go.sum", """
        github.com/spf13/pflag v1.0.3 h1:zPAT6CGy6wXeQ7NtTnaTerfKOsV6V6F8agHXFiazDkg=
        github.com/spf13/pflag v1.0.3/go.mod h1:DYY7MBk1bdzusC3SYhjObp+wFpr4gzcvqqNjLnInEg4=
    """)

    cmd.run("go mod tidy")
    cmd.run("go mod download")

    output = cmd.run("go run main.go Test1234")
    assert output == "Test1234"
