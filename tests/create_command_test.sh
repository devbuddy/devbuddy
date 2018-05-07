set -uo pipefail

oneTimeSetUp() {
    # This runs in a subshell (subprocess).
    # Even if Dad was enabled in the parent shell, we need to enable it here.
    eval "$(dad --shell-init)"
}

setUp() {
    mkdir -p ~/src/github.com/dadorg/repo2
    cd
}

tearDown() {
    cd  # Get out of the repos being deleted
    rm -rf ~/src/github.com/dadorg/repo
    rm -rf ~/src/github.com/dadorg/repo2
}

testCreate() {
    dad create dadorg/repo

    assertTrue "directory was created" "[ -d ~/src/github.com/dadorg/repo ]"
    assertEquals "cd to project directory" ~/src/github.com/dadorg/repo "$PWD"
}

testCreateExisting() {
    dad create dadorg/repo2

    assertEquals "cd to project directory" ~/src/github.com/dadorg/repo2 "$PWD"
}

SHUNIT_PARENT=$0
. ./shunit2.sh
