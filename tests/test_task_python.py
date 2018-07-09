
def test_task(cmd, project):
    project.write_devyml("""
        up:
        - python: 3.6.5
    """)

    output = cmd.run("bud up")
    cmd.assert_succeed()

    output = cmd.run("python --version")
    assert output == "Python 3.6.5"

    output = cmd.run("python -c 'import sys; print(sys.prefix)'")
    assert "/.local/share/bud/virtualenvs/" in output, "the virtualenv is not properly activated"
