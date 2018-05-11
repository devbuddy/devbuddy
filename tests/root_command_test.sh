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
    output=$(dad | head -n1)
    rc=$?
    assertEquals "dad command returns zero" 0 $rc

    assertEquals "dad command output the usage message" "Usage:" "${output}"
}

testDadVersion() {
    output=$(dad --version)
    rc=$?

    assertEquals "dad command returns zero" 0 $rc
    assertEquals "dad command output the usage message" "dad version devel" "${output}"
}

SHUNIT_PARENT=$0
. ./shunit2.sh
