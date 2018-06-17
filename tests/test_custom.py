import pytest


@pytest.fixture
def project(cmd, project_factory):
    p = project_factory('devbuddy_tests', 'poipoi')
    cmd.run(f"cd {p.path}")
    return p


def test_should_run(cmd, project):
    project.write_devyml("""
        up:
        - custom:
            name: TestCustom
            met?: test -e sentinel
            meet: touch sentinel
    """)

    cmd.run('bud up')
    project.assert_file("sentinel")


def test_should_not_run(cmd, project):
    project.write_devyml("""
        up:
        - custom:
            name: TestCustom
            met?: 'true'
            meet: touch sentinel
    """)

    cmd.run("bud up")
    project.assert_file("sentinel", present=False)


def test_should_run_and_fail(cmd, project):
    project.write_devyml("""
        up:
        - custom:
            name: TestCustom
            met?: 'false'
            meet: 'false'
    """)

    output = cmd.run("bud up")
    cmd.assert_failed()

    assert 'command failed' in output



