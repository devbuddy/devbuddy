#!/bin/bash
set -u

oneTimeSetUp() {
    eval "$(dad --shell-init)"
}

setUp() {
    cd $SHUNIT_TMPDIR

    cat > dev.yml <<YAML
up: []
commands: {}
YAML
}

testDadUsage() {
    output=$(dad)
    rc=$?

    lines=(${output[@]})

    assertEquals "dad command returns zero" 0 $rc
    assertEquals "dad command output the usage message" "Usage:" "${lines[0]}"
}

testDadVersion() {
    output=$(dad --version)
    rc=$?

    assertEquals "dad command returns zero" 0 $rc
    assertEquals "dad command output the usage message" "dad version devel" "${output}"
}

SHUNIT_COLOR='none'  # Not macos compatible?
. shunit2/shunit2
