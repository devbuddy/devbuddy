
def test_create(cmd):
    cmd.run("bud create devbuddy_tests/newproject")

    output = cmd.run("pwd")
    assert output.endswith("/src/github.com/devbuddy_tests/newproject")


def test_existing_project(cmd, project):
    output = cmd.run(f"bud create {project.org}/{project.name}")
    assert "github.com:devbuddy_tests/poipoi" in output

    output = cmd.run("pwd")
    assert output.endswith("/src/github.com/devbuddy_tests/poipoi")

