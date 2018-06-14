import pytest


def test_find_project(shell, make_test_repo):
    path_1 = make_test_repo('devbuddy_tests/project')
    path_2 = make_test_repo('devbuddy_tests/project2')

    tests = {
        'devbuddy_tests/project': path_1,
        'devbuddy_tests/project2': path_2,

        'devbuddyproject': path_1,
        'devbuddyproject2': path_2,

        'proj': path_1,

        'pro2': path_2,

        'dtp': path_1,
        'dtp2': path_2,
    }

    for arg, path in tests.items():
        shell.run_command('bud cd %s' % arg)

        output = shell.run_command('pwd')
        found = output.strip()
        assert found == path, "'%s' should match project '%s', not '%s'" % (arg, path, found)


def test_ui(shell, make_test_repo):
    name = 'devbuddy_tests/repo'
    make_test_repo(name)

    output = shell.run_command('bud cd %s' % name)
    assert 'jumping to' in output
    assert name in output
