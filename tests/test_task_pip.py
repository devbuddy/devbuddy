
def test_requirements(cmd, project):
    project.write_devyml("""
        up:
        - python: 3.6.3
        - pip: [requirements.txt]
    """)
    project.write_file("requirements.txt", "test==2.3.4.5\n")

    cmd.run("bud up")

    output = cmd.run("pip freeze")
    assert 'test==2.3.4.5' in output
