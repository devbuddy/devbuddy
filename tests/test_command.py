
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


def test_debug_info(cmd, project):
    debug_info = cmd.run('bud --debug-info')
    assert '**DevBuddy version**' in debug_info
    assert 'SHELL="/bin/bash"' in debug_info
    assert 'Project not found' in debug_info

    project.write_devyml("")
    debug_info = cmd.run('bud --debug-info')
    assert 'Project path: ' in debug_info
