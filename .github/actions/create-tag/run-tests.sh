#!/usr/bin/env bash
set -e
source $(dirname $0)/../_shared/bash/output.sh
source $(dirname $0)/../_shared/bash/testresults.sh
source $(dirname $0)/test-setup.sh


setUp
header

# we created v1.5.0-clash.0 v1.5.0-clash.1 already, so this 
# should return a different tag as a prerelease
# the new tag should contain most of the original tag
# with the prerelease segment changing
COUNT=$((COUNT+1))
out=$(python ./create-tag.py \
    --repository_root="${REPO_ROOT}" \
    --commitish="${COMMITISH1}" \
    --tag_name="v1.5.0-clash.1")
test_val=$(test_output_value "created_tag" "${out}")
test_should_not_equal "${COUNT}" "v1.5.0-clash.1" "${test_val}"
removeTag "${test_val}"
COUNT=$((COUNT+1))
test_should_contain "${COUNT}" "v1.5.0-" "${test_val}"


# this is a brand new tag, so should return the matching value
# should return a different tag
COUNT=$((COUNT+1))
out=$(python ./create-tag.py \
    --repository_root="${REPO_ROOT}" \
    --commitish="${COMMITISH1}" \
    --tag_name="v999.0.1")

test_val=$(test_output_value "created_tag" "${out}")
test_should_equal "${COUNT}" "v999.0.1" "${test_val}"
removeTag "${test_val}"

# this release version already exists, so it should bump along
# the major number and reset the minor / patch
COUNT=$((COUNT+1))
out=$(python ./create-tag.py \
    --repository_root="${REPO_ROOT}" \
    --commitish="${COMMITISH1}" \
    --tag_name="v9999.1.0")
test_val=$(test_output_value "created_tag" "${out}")
test_should_equal "${COUNT}" "v10000.0.0" "${test_val}"
removeTag "${test_val}"

# this release version already exists, as do several others 
# in this space, so the tag should jump ahead to an
# unused number
COUNT=$((COUNT+1))
out=$(python ./create-tag.py \
    --repository_root="${REPO_ROOT}" \
    --commitish="${COMMITISH1}" \
    --tag_name="v21.0.1")
test_val=$(test_output_value "created_tag" "${out}")
test_should_equal "${COUNT}" "v24.0.0" "${test_val}"
removeTag "${test_val}"

# Test non-semver tag
COUNT=$((COUNT+1))
out=$(python ./create-tag.py \
    --repository_root="${REPO_ROOT}" \
    --commitish="${COMMITISH1}" \
    --tag_name="test_tag")
test_val=$(test_output_value "created_tag" "${out}")
test_should_contain "${COUNT}" "test_tag." "${test_val}"


tearDown

if [ -n "${TEST_ERR}" ]; then
    echo "ERROR"
    exit 1
fi