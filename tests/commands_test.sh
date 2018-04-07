#!/bin/bash

oneTimeSetUp() {
	# This runs in a subshell (subprocess).
	# Even if Dad was enabled in the parent shell, we need to enable it here.
	eval "$(dad --shell-init)"
}

setUp() {
	rm -rf ~/src/github.com/dadorgtest/repotest
	rm -rf ~/src/github.com/dadorgtest/repotest2
	mkdir -p ~/src/github.com/dadorgtest/repotest
	mkdir -p ~/src/github.com/dadorgtest/repotest2

	cd  # reset the cwd since this is changed in the tests
}

testCreate() {
	rm -rf ~/src/github.com/dadorgtest/repotest

	dad create dadorgtest/repotest

	assertTrue "directory was created" "[ -d ~/src/github.com/dadorgtest/repotest ]"
	assertEquals "cd to project directory" ~/src/github.com/dadorgtest/repotest "$PWD"
}

testCdExactMatch() {
	dad cd dadorgtest/repotest
	assertEquals "cd to project directory" ~/src/github.com/dadorgtest/repotest "$PWD"

	dad cd dadorgtest/repotest2
	assertEquals "cd to second project directory" ~/src/github.com/dadorgtest/repotest2 "$PWD"
}

# The search will stop at the first match

testCdByExactRepo() {
	dad cd repotest2
	assertEquals "cd to project directory" ~/src/github.com/dadorgtest/repotest2 "$PWD"
}

testCdByExactOrg() {
	dad cd dadorgtest
	assertEquals "cd to project directory" ~/src/github.com/dadorgtest/repotest "$PWD"
}

testCdByPrefixRepo() {
	dad cd repotes
	assertEquals "cd to project directory" ~/src/github.com/dadorgtest/repotest "$PWD"
}

testCdByPrefixOrg() {
	dad cd dadorgtes
	assertEquals "cd to project directory" ~/src/github.com/dadorgtest/repotest "$PWD"
}

testCdBySubstringRepo() {
	dad cd epotest2
	assertEquals "cd to project directory" ~/src/github.com/dadorgtest/repotest2 "$PWD"
}

testCdBySubstringOrg() {
	dad cd adorgtest
	assertEquals "cd to project directory" ~/src/github.com/dadorgtest/repotest "$PWD"
}

SHUNIT_COLOR='none'  # Not macos compatible?
. shunit2/shunit2
