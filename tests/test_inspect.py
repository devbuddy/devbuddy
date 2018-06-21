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
    cmd.assert_succeed()

    assert 'some-condition-command' in output


def test_without_manifest(cmd, project):
    output = cmd.run('bud inspect')
    cmd.assert_failed()
