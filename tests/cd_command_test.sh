set -u

oneTimeSetUp() {
    # This runs in a subshell (subprocess).
    # Even if DevBuddy was enabled in the parent shell, we need to enable it here.
    eval "$(bud --shell-init)"
}

setUp() {
    rm -rf ~/src/github.com/dadorgtest/repotest
    rm -rf ~/src/github.com/dadorgtest/repotest2
    mkdir -p ~/src/github.com/dadorgtest/repotest
    mkdir -p ~/src/github.com/dadorgtest/repotest2

    cd  # reset the cwd since this is changed in the tests
}

testCdExactMatch() {
    bud cd dadorgtest/repotest
    assertEquals "cd to project directory" ~/src/github.com/dadorgtest/repotest "$PWD"

    bud cd dadorgtest/repotest2
    assertEquals "cd to second project directory" ~/src/github.com/dadorgtest/repotest2 "$PWD"
}

# The search will stop at the first match

testCdByExactRepo() {
    bud cd repotest2
    assertEquals "cd to project directory" ~/src/github.com/dadorgtest/repotest2 "$PWD"
}

testCdByExactOrg() {
    bud cd dadorgtest
    assertEquals "cd to project directory" ~/src/github.com/dadorgtest/repotest "$PWD"
}

testCdByPrefixRepo() {
    bud cd repotes
    assertEquals "cd to project directory" ~/src/github.com/dadorgtest/repotest "$PWD"
}

testCdByPrefixOrg() {
    bud cd dadorgtes
    assertEquals "cd to project directory" ~/src/github.com/dadorgtest/repotest "$PWD"
}

testCdBySubstringRepo() {
    bud cd epotest2
    assertEquals "cd to project directory" ~/src/github.com/dadorgtest/repotest2 "$PWD"
}

testCdBySubstringOrg() {
    bud cd adorgtest
    assertEquals "cd to project directory" ~/src/github.com/dadorgtest/repotest "$PWD"
}

SHUNIT_PARENT=$0
. ./shunit2.sh
