set -u

oneTimeSetUp() {
    eval "$(bud --shell-init)"
}

setUp() {
    cd $SHUNIT_TMPDIR

    cat > dev.yml <<YAML
commands:
  mycmd:
    run: echo TESTTEST > somefile
  echo:
    run: echo PREFIX
  isprojectroot:
    run: test -e dev.yml
  success:
    run: true
  failure:
    run: false
YAML
}

testSimple() {
    bud mycmd

    assertEquals "somefile was created" "TESTTEST" "$(cat somefile)"
}

testArguments() {
    output=$(bud echo ARG1 ARG2)

    assertEquals "command called with arguments" "PREFIX ARG1 ARG2" "$output"
}

testSuccess() {
    bud success

    assertEquals "bud return with the right exit code" 0 $?
}

testFailure() {
    bud failure

    assertEquals "bud return with the right exit code" 1 $?
}

testRunInProjectRoot() {
    bud isprojectroot
    assertEquals "bud isprojectroot command succeed in project root" 0 $?

    mkdir -p subdir
    cd subdir
    bud isprojectroot
    assertEquals "bud custom commands run in project root" 0 $?
}

SHUNIT_PARENT=$0
. ./shunit2.sh
