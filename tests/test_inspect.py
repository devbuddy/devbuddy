import pytest


def test_invalid_manifest_with_string(cmd, project):
    project.write_devyml("""
        up:
        - custom:
            name: Title
            met?: some-condition-command
            meet: some-action-command
        - pip: [requirements.txt]
    """)

    output = cmd.run('bud inspect')
    assert 'Title' in output
    assert 'some-condition-command' in output
    assert 'requirements.txt' in output


def test_without_manifest(cmd, project):
    cmd.run('bud inspect', expect_exit_code=1)
