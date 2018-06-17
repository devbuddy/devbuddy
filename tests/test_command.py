
def test_option_version(cmd):
    output = cmd.run('which bud')
    assert '/bud' in output

    output = cmd.run('bud --version')
    assert 'bud version devel' in output


def test_usage(cmd):
    output = cmd.run('bud')
    assert output.startswith('Usage:')


def test_cmd_init(cmd):
    init = cmd.run('bud --shell-init')
    assert 'bud() {' in init
    assert '__bud_prompt_command' in init
