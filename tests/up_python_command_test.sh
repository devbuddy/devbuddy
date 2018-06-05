set -u

oneTimeSetUp() {
    eval "$(bud --shell-init)"

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
six==1.10.0
EOF
    cat > requirements-2.txt <<EOF
dataclasses==0.4
EOF
}

testEnv() {
    output=$(bud up)
    rc=$?
    assertEquals "bud command returns zero" 0 $rc

    eval "$(command bud --shell-hook)"  # Simulate a prompt hook

    version=$(python --version)
    assertEquals "python version" "Python 3.6.3" "$version"

    version=$(pip show six | grep Version)
    assertEquals "package version" "Version: 1.10.0" "$version"

    version=$(pip show dataclasses | grep Version)
    assertEquals "package version" "Version: 0.4" "$version"
}

testFailure() {
    rm requirements-2.txt

    output=$(bud up)
    rc=$?
    assertEquals "bud command returns non-zero" 1 $rc
}

SHUNIT_PARENT=$0
. ./shunit2.sh
