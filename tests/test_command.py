

def test_option_version(shell):

    output = shell.run_command('which bud')
    assert '/bud' in output

    output = shell.run_command('bud --version')
    assert 'bud version devel' in output


def test_usage(shell):
    output = shell.run_command('bud')
    assert output.startswith('Usage:')


def test_shell_init(shell):
    init = shell.run_command('bud --shell-init')
    assert 'bud() {' in init
    assert '__bud_prompt_command' in init


def test_command_cd(shell, make_test_repo):
    path = make_test_repo('testrepo')

    output = shell.run_command('bud cd testrepo')
    assert 'jumping to' in output
    assert 'devbuddy_integration_tests/testrepo' in output

    output = shell.run_command('pwd')
    assert output.strip() == path

    path2 = make_test_repo('testrepo2')
    shell.run_command('bud cd testrepo2')
    output = shell.run_command('pwd')
    assert output.strip() == path2
