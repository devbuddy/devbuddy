
def test_set_env(cmd, project):
    project.write_devyml("""
        env:
          TESTVAR: TESTVALUE
        up: []
    """)

    cmd.run("true")

    output = cmd.run("echo $TESTVAR")
    assert "TESTVALUE" in output
