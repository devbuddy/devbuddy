set -u

oneTimeSetUp() {
    # This runs in a subshell (subprocess).
    # Even if DevBuddy was enabled in the parent shell, we need to enable it here.
    eval "$(bud --shell-init)"
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
    bud create dadorg/repo

    assertTrue "directory was created" "[ -d ~/src/github.com/dadorg/repo ]"
    assertEquals "cd to project directory" ~/src/github.com/dadorg/repo "$PWD"
}

testCreateExisting() {
    bud create dadorg/repo2

    assertEquals "cd to project directory" ~/src/github.com/dadorg/repo2 "$PWD"
}

SHUNIT_PARENT=$0
. ./shunit2.sh
