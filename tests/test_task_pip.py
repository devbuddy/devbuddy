
def test_requirements(cmd, project):
    project.write_devyml("""
        up:
        - python: 3.6.5
        - pip: [requirements.txt]
    """)
    project.write_file("requirements.txt", "pyreleaser==0.5.2\n")

    cmd.run("bud up")

    output = cmd.run("pip freeze")
    assert 'pyreleaser==0.5.2' in output
