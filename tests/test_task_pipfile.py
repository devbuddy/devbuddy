
def test_simple(cmd, project):
    project.write_devyml("""
        up:
        - python: 3.6.5
        - pipfile
    """)
    project.write_file("Pipfile", """[packages]\n"pyreleaser" = "==0.5.2"\n""")

    cmd.run("bud up")

    output = cmd.run("pip freeze")
    assert 'pyreleaser==0.5.2' in output
