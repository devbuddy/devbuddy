
def test_create(cmd):
    cmd.run("bud create devbuddy_tests/newproject")

    output = cmd.run("pwd")
    assert output.endswith("/src/github.com/devbuddy_tests/newproject")


def test_existing_project(cmd, project):
    output = cmd.run(f"bud create {project.org}/{project.name}")
    assert "project already exists locally" in output

    output = cmd.run("pwd")
    assert output.endswith(f"/src/github.com/{project.org}/{project.name}")

