set -u

oneTimeSetUp() {
    eval "$(bud --shell-init)"
}

setUp() {
    cd $SHUNIT_TMPDIR
}

testInvalidUpSection() {
    cat > dev.yml <<YAML
up: blabla
YAML

    output=$(bud up)
    rc=$?
    assertEquals "command should have failed" 1 $rc
}

testUnknownTask() {
    cat > dev.yml <<YAML
up:
  - nopenopenope
YAML

    output=$(bud up)
    rc=$?
    assertEquals "command should have succeed" 0 $rc
}

testInvalidTask() {
    cat > dev.yml <<YAML
up:
  - true
YAML

    output=$(bud up)
    rc=$?
    assertEquals "command should have succeed" 0 $rc
}

SHUNIT_PARENT=$0
. ./shunit2.sh
