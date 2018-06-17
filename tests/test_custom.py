import pytest


@pytest.fixture
def project(cmd, project_factory):
    p = project_factory('devbuddy_tests', 'poipoi')
    cmd.run(f"cd {p.path}")
    return p


def test_command(cmd, project):
    project.write_devyml("""
        commands:
          mycmd:
            run: touch somefile
    """)
    cmd.run("bud mycmd")
    project.assert_file('somefile')


def test_run_in_project_dir(cmd, project):
    project.write_devyml("""
        commands:
          mycmd:
            run: touch somefile
    """)
    cmd.run('mkdir subdir')
    cmd.run('cd subdir')
    cmd.run("bud mycmd")
    project.assert_file('somefile')


def test_exit_code(cmd, project):
    project.write_devyml("""
        commands:
          success:
            run: 'true'
          failure:
            run: 'false'
    """)

    cmd.run("bud success")
    cmd.assert_succeed()

    cmd.run("bud failure")
    cmd.assert_failed()


def test_with_arguments(cmd, project):
    project.write_devyml("""
        commands:
          mycmd:
            run: echo PREFIX
    """)
    output = cmd.run("bud mycmd ARG1 ARG2")
    lines = [l for l in output.splitlines() if not l.startswith("🐼")]
    assert ["PREFIX ARG1 ARG2"] == lines


