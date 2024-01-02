#!/usr/bin/env bash
set -e
export RUN_AS_TEST="true"
# test commit to use
COMMITISH1="dab1b63"
# test counter
COUNT="0"
# for local
REPO_ROOT="./repo-test/"
#error flag
TEST_ERR=""

# generate these tags for tests
declare -a dummy_tags=(
    test_tag
    test_tag_wont_be_latest 
    v2.4.9 
    v2.4.15 
    v1.0.1 
    v20.10.10 
    v20.10.9 
    v21.0.1
    v22.0.0
    v23.0.0
    v1.5.0-clash.0 
    v1.5.0-clash.1
    v9999.1.0
)


setUp() {
    if [ -d "${REPO_ROOT}" ]; then
        rm -Rf "${REPO_ROOT}"        
    fi
    git clone https://github.com/ministryofjustice/opg-github-actions.git ${REPO_ROOT} &>/dev/null
    
    # create some dummy tags
    cd ${REPO_ROOT} 
    for tag in "${dummy_tags[@]}"; do
        git tag -d "${tag}" &>/dev/null || true 
        git tag "${tag}" "${COMMITISH1}"
    done
    cd - &>/dev/null
}

tearDown() {
    # remote dummy tags
    cd ${REPO_ROOT}
    for tag in "${dummy_tags[@]}"; do
        git tag -d "${tag}" &>/dev/null
    done
    
    cd - &>/dev/null
    if [ -d "${REPO_ROOT}" ]; then
        rm -Rf "${REPO_ROOT}"
    fi
}

removeTag() {
    cd ${REPO_ROOT}
    git tag -d "${1}" &>/dev/null || true 
    cd - &>/dev/null
}