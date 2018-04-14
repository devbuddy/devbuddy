#!/bin/bash
set -u

oneTimeSetUp() {
    eval "$(dad --shell-init)"

    # Check the test requirements
    pyenv versions | grep -q 3.6.3
}

setUp() {
    cd $SHUNIT_TMPDIR

    cat > dev.yml <<YAML
up:
  - python: 3.6.3
  - pip: [requirements.txt, requirements-2.txt]
YAML

    cat > requirements.txt <<EOF
pyreleaser==0.3
EOF
    cat > requirements-2.txt <<EOF
dataclasses==0.4
EOF
}

testEnv() {
    output=$(dad up)
    rc=$?
    assertEquals "dad command returns zero" 0 $rc

    eval "$(command dad --shell-hook)"  # Simulate a prompt hook

    version=$(python --version)
    assertEquals "python version" "Python 3.6.3" "$version"

    version=$(pip show pyreleaser | grep Version)
    assertEquals "package version" "Version: 0.3" "$version"

    version=$(pip show dataclasses | grep Version)
    assertEquals "package version" "Version: 0.4" "$version"
}

testFailure() {
    rm requirements-2.txt

    output=$(dad up)
    rc=$?
    assertEquals "dad command returns non-zero" 1 $rc
}

SHUNIT_COLOR='none'  # Not macos compatible?
. shunit2/shunit2
