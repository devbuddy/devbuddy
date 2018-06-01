set -u

oneTimeSetUp() {
    eval "$(bud --shell-init)"
}

setUp() {
    cd $SHUNIT_TMPDIR

    cat > dev.yml <<YAML
up: []
commands: {}
YAML
}

testDadUsage() {
    output=$(bud | head -n1)
    rc=$?
    assertEquals "bud command returns zero" 0 $rc

    assertEquals "bud command output the usage message" "Usage:" "${output}"
}

testDadVersion() {
    output=$(bud --version)
    rc=$?

    assertEquals "bud command returns zero" 0 $rc
    assertEquals "bud command output the usage message" "bud version devel" "${output}"
}

SHUNIT_PARENT=$0
. ./shunit2.sh
