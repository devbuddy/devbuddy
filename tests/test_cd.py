import pytest


def test_find_project(cmd, project_factory):
    project_1 = project_factory('devbuddy_tests', 'project')
    project_2 = project_factory('devbuddy_tests', 'project2')

    tests = {
        'devbuddy_tests/project': project_1.path,
        'devbuddy_tests/project2': project_2.path,

        'devbuddyproject': project_1.path,
        'devbuddyproject2': project_2.path,

        'proj': project_1.path,

        'pro2': project_2.path,

        'dtp': project_1.path,
        'dtp2': project_2.path,
    }

    for arg, path in tests.items():
        cmd.run('bud cd %s' % arg)

        output = cmd.run('pwd')
        found = output.strip()
        assert found == path, "'%s' should match project '%s', not '%s'" % (arg, path, found)


def test_ui(cmd, project_factory):
    project_factory('devbuddy_tests', 'repo')

    output = cmd.run('bud cd devbuddy_tests/repo')
    assert 'jumping to' in output
    assert 'devbuddy_tests/repo' in output
