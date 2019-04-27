import yaml


def test_create(cmd, tmpdir):
    cmd.run(f"cd {tmpdir}")

    output = cmd.run("bud init", check_activation_failure=False)
    assert "Open dev.yml to adjust for your needs." in output

    manifest = tmpdir.join("dev.yml")
    assert manifest.check(file=1)

    validate_manifest(manifest.open().read())


def test_create_existing(cmd, tmpdir):
    cmd.run(f"cd {tmpdir}")
    cmd.run("touch dev.yml")

    output = cmd.run("bud init", expect_exit_code=1)
    assert "creation failed" in output
    assert "file exists" in output


def validate_manifest(as_string: str) -> None:
    content = yaml.load(as_string)

    assert "up" in content
    assert isinstance(content["up"], list)
    for task in content["up"]:
        assert isinstance(task, (str, dict))

    if "commands" in content:
        assert isinstance(content["commands"], dict)

    if "open" in content:
        assert isinstance(content["open"], dict)
