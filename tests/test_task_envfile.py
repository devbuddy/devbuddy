import textwrap


def test_task(cmd, project):
    project.write_devyml("""
        up:
        - envfile
    """)

    project.write_file(".env", textwrap.dedent("""
        TESTVAR=onetwo
        PATH=nopenope
    """))

    cmd.run("true")

    output = cmd.run("echo ${TESTVAR}")
    assert "onetwo" in output
    output = cmd.run("echo ${PATH}")
    assert "nopenope" not in output
