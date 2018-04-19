#!/bin/bash
set -u

oneTimeSetUp() {
    eval "$(dad --shell-init)"
}

setUp() {
    cd $SHUNIT_TMPDIR
}

testInvalidUpSection() {
    cat > dev.yml <<YAML
up: blabla
YAML

    output=$(dad up)
    rc=$?
    assertEquals "command should have failed" 1 $rc
}

testUnknownTask() {
    cat > dev.yml <<YAML
up:
  - nopenopenope
YAML

    output=$(dad up)
    rc=$?
    assertEquals "command should have succeed" 0 $rc
}

testInvalidTask() {
    cat > dev.yml <<YAML
up:
  - true
YAML

    output=$(dad up)
    rc=$?
    assertEquals "command should have succeed" 0 $rc
}

SHUNIT_COLOR='none'  # Not macos compatible?
. shunit2/shunit2
