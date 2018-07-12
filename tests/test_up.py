import pytest


def test_invalid_manifest_with_string(cmd, project):
    project.write_devyml("""
        up: somestring
    """)

    output = cmd.run('bud up', expect_exit_code=1)

    # This test exists to show how bad the output is. This should be improved.
    assert 'yaml: unmarshal errors' in output
    assert 'cannot unmarshal !!str `somestring` into []interface {}' in output


def test_unknown_task(cmd, project):
    project.write_devyml("""
        up:
          - notatask
    """)

    output = cmd.run('bud up')

    assert 'notatask' in output
    assert 'Unknown task' in output


def test_invalid_task(cmd, project):
    project.write_devyml("""
        up:
          - true
    """)

    output = cmd.run('bud up', expect_exit_code=1)
    assert 'invalid task: "true"' in output
