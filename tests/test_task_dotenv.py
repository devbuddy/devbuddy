
def test_env(cmd, project):
    project.write_devyml("""
        up:
        - dotenv
    """)
    project.write_file(".env", "ENVVAR1=VALUE1\nENVVAR2=VALUE2")

    cmd.run("bud up")

    output = cmd.run("echo __${ENVVAR1}__")
    assert output == "__VALUE1__"

    output = cmd.run("echo __${ENVVAR2}__")
    assert output == "__VALUE2__"
