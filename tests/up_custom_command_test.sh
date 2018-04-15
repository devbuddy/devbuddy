#!/bin/bash
set -u

oneTimeSetUp() {
    eval "$(dad --shell-init)"
}

setUp() {
    cd $SHUNIT_TMPDIR
}

testNotMetThenMet() {
    cat > dev.yml <<YAML
up:
  - custom:
      met?: 'test -e testfile'
      meet: 'touch testfile'
YAML

    output=$(dad up)
    rc=$?
    assertEquals "command failed" 0 $rc

    assertTrue "the test file was not created" '[ -e testfile ]'
}

testAlreadyMet() {
    cat > dev.yml <<YAML
up:
  - custom:
      met?: 'true'
      meet: 'touch sentinel'
YAML

    output=$(dad up)
    rc=$?
    assertEquals "command failed" 0 $rc

    assertFalse "the custom task should not have ran" '[ -e sentinel ]'
}

testNotMetThenFail() {
    cat > dev.yml <<YAML
up:
  - custom:
      met?: 'false'
      meet: 'false'
YAML

    output=$(dad up)
    rc=$?
    assertEquals "command did not failed" 1 $rc
}

SHUNIT_COLOR='none'  # Not macos compatible?
. shunit2/shunit2
