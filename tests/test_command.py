
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
