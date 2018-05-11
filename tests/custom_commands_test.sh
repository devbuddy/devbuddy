set -u

oneTimeSetUp() {
    eval "$(dad --shell-init)"
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
    dad mycmd

    assertEquals "somefile was created" "TESTTEST" "$(cat somefile)"
}

testArguments() {
    output=$(dad echo ARG1 ARG2)

    assertEquals "command called with arguments" "PREFIX ARG1 ARG2" "$output"
}

testSuccess() {
    dad success

    assertEquals "dad return with the right exit code" 0 $?
}

testFailure() {
    dad failure

    assertEquals "dad return with the right exit code" 1 $?
}

testRunInProjectRoot() {
    dad isprojectroot
    assertEquals "dad isprojectroot command succeed in project root" 0 $?

    mkdir -p subdir
    cd subdir
    dad isprojectroot
    assertEquals "dad custom commands run in project root" 0 $?
}

SHUNIT_PARENT=$0
. ./shunit2.sh
