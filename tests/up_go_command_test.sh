set -u

oneTimeSetUp() {
    eval "$(dad --shell-init)"
}

setUp() {
    cd $SHUNIT_TMPDIR

    cat > dev.yml <<YAML
up:
  - go: '1.5'
YAML

    rm -rf gopath
    export GOPATH="$PWD/gopath"
}

testMissingGOPATH() {
    unset GOPATH

    output=$(dad up | grep 'Warning.*GOPATH')  # Expect a warning about GOPATH
    rc=$?
    assertEquals "dad up should warn about GOPATH not set" 0 $rc
}

testEnv() {
    output=$(dad up)
    rc=$?
    assertEquals "dad command returns zero" 0 $rc

    eval "$(command dad --shell-hook)"  # Simulate a prompt hook

    version=$(go version | cut -d ' ' -f 3)
    assertEquals "go version" "go1.5" "$version"
}

SHUNIT_PARENT=$0
. ./shunit2.sh
