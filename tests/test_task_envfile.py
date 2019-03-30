import textwrap


def test_load_envfile_in_process(cmd, project):
    project.write_devyml("""
        up:
        - envfile
        - custom:
            name: succeed if TESTVAR is set
            met?: test "${TESTVAR}" == "1"
            meet: echo "TESTVAR is not set"; false
    """)

    project.write_file(".env", textwrap.dedent("""
        TESTVAR=1
    """))

    cmd.run("bud up")


def test_load_envfile_in_shell(cmd, project):
    project.write_devyml("""
        up:
        - envfile
    """)

    project.write_file(".env", textwrap.dedent("""
        TESTVAR=onetwo
    """))

    cmd.run("true")  # run the shell hook

    output = cmd.run("echo ${TESTVAR}")
    assert "onetwo" in output
