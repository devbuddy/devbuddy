set -u

oneTimeSetUp() {
    eval "$(bud --shell-init)"
}

setUp() {
    cd $SHUNIT_TMPDIR
    rm -f dev.yml
}

testWithManifest() {
    cat > dev.yml <<YAML
up:
  - python: 3.6.3
  - pip: [requirements.txt]
YAML

    output=$(bud inspect)
    rc=$?
    assertEquals "command failed" 0 $rc
}

testWithouthManifest() {
    output=$(bud inspect)
    rc=$?
    assertEquals "command didn't fail" 1 $rc
}


SHUNIT_PARENT=$0
. ./shunit2.sh
