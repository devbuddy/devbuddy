
def test_simple(cmd, project):
    project.write_devyml("""
        up:
        - python: 3.6.5
        - pipfile
    """)
    project.write_file("Pipfile", """[packages]\n"test" = "==2.3.4.5"\n""")

    output = cmd.run("bud up")

    output = cmd.run("pip freeze")
    assert 'test==2.3.4.5' in output
