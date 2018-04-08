#!/bin/bash

oneTimeSetUp() {
	eval "$(dad --shell-init)"
}

setUp() {
	rm -rf ~/src/github.com/dadorg/repo
	mkdir -p ~/src/github.com/dadorg/repo
	cd ~/src/github.com/dadorg/repo

	cat > dev.yml <<YAML
commands:
	mycmd:
		run: echo TESTTEST > somefile
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

testSuccess() {
	dad success

	assertEquals "dad return with the right exit code" 0 $?
}

testFailure() {
	dad failure

	assertEquals "dad return with the right exit code" 1 $?
}

SHUNIT_COLOR='none'  # Not macos compatible?
. shunit2/shunit2
