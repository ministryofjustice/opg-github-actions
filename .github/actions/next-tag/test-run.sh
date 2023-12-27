#!/usr/bin/env bash
set -e
export RUN_AS_TEST="true"
err=""

echo "## Test Information"
echo "| \# | Expected | Actual | "
echo "| --- | --- | --- |"

out=$(python ./next-tag.py \
    --test_file=./tests/majors.txt \
    --prerelease_suffix="moreactions" \
    --prerelease=true \
    --latest_tag="v1.5.0-moreactions.1" \
    --default_bump="patch" \
    --last_release="v1.4.0")

actual=$(echo "${out}" | sed -r -n 's/.*next_tag=(.*)$/\1/p' )
expected="2.0.0-moreactions.0"
echo "| 1.1 | ${expected} | ${actual} |"
if [ "${actual}" != "${expected}" ]; then
    err="1"
    echo "FAILED #1.1"
    echo "${out}"
    echo "==="
fi

out=$(python ./next-tag.py \
    --test_file=./tests/majors.txt \
    --prerelease_suffix="moreactions" \
    --prerelease=true \
    --latest_tag="v2.0.0-moreactions.0" \
    --default_bump="patch" \
    --last_release="v1.4.0")

actual=$(echo "${out}" | sed -r -n 's/.*next_tag=(.*)$/\1/p' )
expected="2.0.0-moreactions.1"
echo "| 1.2 | ${expected} | ${actual} |"
if [ "${actual}" != "${expected}" ]; then
    err="1"
    echo "FAILED #1.2"
    echo "${out}"
    echo "==="
fi


out=$(python ./next-tag.py \
    --test_file=./tests/majors.txt \
    --prerelease_suffix="moreactions" \
    --prerelease=true \
    --latest_tag="" \
    --default_bump="patch" \
    --last_release="v1.4.0")
actual=$(echo "${out}" | sed -r -n 's/.*next_tag=(.*)$/\1/p' )
expected="2.0.0-moreactions.0"
echo "| 1.3 | ${expected} | ${actual} |"
if [ "${actual}" != "${expected}" ]; then
    err="1"
    echo "FAILED #1.3"
    echo "${out}"
    echo "==="
fi


out=$(python ./next-tag.py \
    --test_file=./tests/majors.txt \
    --prerelease_suffix="moreactions" \
    --latest_tag="" \
    --default_bump="patch" \
    --last_release="v1.4.0")
actual=$(echo "${out}" | sed -r -n 's/.*next_tag=(.*)$/\1/p' )
expected="2.0.0"
echo "| 1.4 | ${expected} | ${actual} |"
if [ "${actual}" != "${expected}" ]; then
    err="2"
    echo "FAILED #1.4"
    echo "${out}"
    echo "==="
fi


out=$(python ./next-tag.py \
    --test_file=./tests/majors.txt \
    --prerelease_suffix="moreactions" \
    --default_bump="patch" \
    --latest_tag="" \
    --last_release="")
actual=$(echo "${out}" | sed -r -n 's/.*next_tag=(.*)$/\1/p' )
expected="1.0.0"
echo "| 1.5 | ${expected} | ${actual} |"
if [ "${actual}" != "${expected}" ]; then
    err="3"
    echo "FAILED #1.5"
    echo "${out}"
    echo "==="
fi


out=$(python ./next-tag.py \
    --test_file=./tests/minors.txt \
    --prerelease_suffix="moreactions" \
    --prerelease=true \
    --default_bump="patch" \
    --latest_tag="v1.5.0-moreactions.0" \
    --last_release="v1.4.0")
actual=$(echo "${out}" | sed -r -n 's/.*next_tag=(.*)$/\1/p' )
expected="1.5.0-moreactions.1"
echo "| 2.1 | ${expected} | ${actual} |"
if [ "${actual}" != "${expected}" ]; then
    err="4"
    echo "FAILED #2.1"
    echo "${out}"
    echo "==="
fi

out=$(python ./next-tag.py \
    --test_file=./tests/minors.txt \
    --prerelease_suffix="moreactions" \
    --prerelease=true \
    --default_bump="patch" \
    --latest_tag="" \
    --last_release="v1.4.0")
actual=$(echo "${out}" | sed -r -n 's/.*next_tag=(.*)$/\1/p' )
expected="1.5.0-moreactions.0"
echo "| 2.2 | ${expected} | ${actual} |"
if [ "${actual}" != "${expected}" ]; then
    err="1"
    echo "FAILED #2.2"
    echo "${out}"
    echo "==="
fi


out=$(python ./next-tag.py \
    --test_file=./tests/minors.txt \
    --prerelease_suffix="moreactions" \
    --latest_tag="v1.5.0-moreactions.0" \
    --default_bump="patch" \
    --last_release="v1.4.0")
actual=$(echo "${out}" | sed -r -n 's/.*next_tag=(.*)$/\1/p' )
expected="1.5.0"
echo "| 2.3 | ${expected} | ${actual} |"
if [ "${actual}" != "${expected}" ]; then
    err="5"
    echo "FAILED #2.3"
    echo "${out}"
    echo "==="
fi

out=$(python ./next-tag.py \
    --test_file=./tests/minors.txt \
    --prerelease_suffix="moreactions" \
    --latest_tag="" \
    --last_release="")
actual=$(echo "${out}" | sed -r -n 's/.*next_tag=(.*)$/\1/p' )
expected="0.1.0"
echo "| 2.4 | ${expected} | ${actual} |"
if [ "${actual}" != "${expected}" ]; then
    err="6"
    echo "FAILED #2.4"
    echo "${out}"
    echo "==="
fi


out=$(python ./next-tag.py \
    --test_file=./tests/patch.txt \
    --prerelease=true \
    --prerelease_suffix="moreactions" \
    --latest_tag="" \
    --default_bump="patch" \
    --last_release="v1.4.0")
actual=$(echo "${out}" | sed -r -n 's/.*next_tag=(.*)$/\1/p' )
expected="1.4.1-moreactions.0"
echo "| 3.1 | ${expected} | ${actual} |"
if [ "${actual}" != "${expected}" ]; then
    err="7"
    echo "FAILED #3.1"
    echo "${out}"
    echo "==="
fi

out=$(python ./next-tag.py \
    --test_file=./tests/patch.txt \
    --prerelease=true \
    --default_bump="patch" \
    --prerelease_suffix="moreactions" \
    --latest_tag="v1.4.1-moreactions.1" \
    --last_release="v1.4.0")
actual=$(echo "${out}" | sed -r -n 's/.*next_tag=(.*)$/\1/p' )
expected="1.4.1-moreactions.2"
echo "| 3.2 | ${expected} | ${actual} |"
if [ "${actual}" != "${expected}" ]; then
    err="8"
    echo "FAILED #3.2"
    echo "${out}"
    echo "==="
fi


out=$(python ./next-tag.py \
    --test_file=./tests/patch.txt \
    --default_bump="patch" \
    --prerelease_suffix="moreactions" \
    --latest_tag="v1.4.1-moreactions.1" \
    --last_release="v1.4.0")
actual=$(echo "${out}" | sed -r -n 's/.*next_tag=(.*)$/\1/p' )
expected="1.4.1"
echo "| 3.3 | ${expected} | ${actual} |"
if [ "${actual}" != "${expected}" ]; then
    err="9"
    echo "FAILED #3.3"
    echo "${out}"
    echo "==="
fi

out=$(python ./next-tag.py \
    --test_file=./tests/patch.txt \
    --prerelease_suffix="moreactions" \
    --default_bump="patch" \
    --latest_tag="" \
    --last_release="")    
actual=$(echo "${out}" | sed -r -n 's/.*next_tag=(.*)$/\1/p' )
expected="0.0.1"
echo "| 3.4 | ${expected} | ${actual} |"
if [ "${actual}" != "${expected}" ]; then
    err="10"
    echo "FAILED #3.4"
    echo "${out}"
    echo "==="
fi

# this is a prerelease with an existing tag, so just bump the prerelease
out=$(python ./next-tag.py \
    --test_file=./tests/none.txt \
    --prerelease=true \
    --default_bump="patch" \
    --prerelease_suffix="moreactions" \
    --latest_tag="v1.0.1-moreactions.0" \
    --last_release="v1.0.0")    
actual=$(echo "${out}" | sed -r -n 's/.*next_tag=(.*)$/\1/p' )
expected="1.0.1-moreactions.1"
echo "| 4.1 | ${expected} | ${actual} |"
if [ "${actual}" != "${expected}" ]; then
    err="11"
    echo "FAILED #4.1"
    echo "${out}"
    echo "==="
fi

# this is a prerelease without an existing tag, so should get a minor bump by default and
# a prerelease segment added
out=$(python ./next-tag.py \
    --test_file=./tests/none.txt \
    --prerelease=true \
    --default_bump="patch" \
    --prerelease_suffix="moreactions" \
    --last_release="v1.0.0")    
actual=$(echo "${out}" | sed -r -n 's/.*next_tag=(.*)$/\1/p' )
expected="1.0.1-moreactions.0"
echo "| 4.2 | ${expected} | ${actual} |"
if [ "${actual}" != "${expected}" ]; then
    err="12"
    echo "FAILED #4.2"
    echo "${out}"
    echo "==="
fi


# a prerelease, with a major bump
out=$(python ./next-tag.py \
    --test_file=./tests/none.txt \
    --default_bump="major" \
    --prerelease="true" \
    --prerelease_suffix="moreactions" \
    --last_release="v1.0.0")    
actual=$(echo "${out}" | sed -r -n 's/.*next_tag=(.*)$/\1/p' )
expected="2.0.0-moreactions.0"
echo "| 4.3 | ${expected} | ${actual} |"
if [ "${actual}" != "${expected}" ]; then
    err="13"
    echo "FAILED #4.3"
    echo "${out}"
    echo "==="
fi

# not a prerelease, so should bump the minor version segment
out=$(python ./next-tag.py \
    --test_file=./tests/none.txt \
    --default_bump="minor" \
    --prerelease_suffix="moreactions" \
    --last_release="v1.0.0")    
actual=$(echo "${out}" | sed -r -n 's/.*next_tag=(.*)$/\1/p' )
expected="1.1.0"
echo "| 4.4 | ${expected} | ${actual} |"
if [ "${actual}" != "${expected}" ]; then
    err="14"
    echo "FAILED #4.4"
    echo "${out}"
    echo "==="
fi

# no previous tags, nothing in commits to trigger a bump
# default as patch, not a prerelease
out=$(python ./next-tag.py \
    --test_file=./tests/none.txt \
    --prerelease_suffix="moreactions" \
    --latest_tag="" \
    --default_bump="patch" \
    --last_release="")
actual=$(echo "${out}" | sed -r -n 's/.*next_tag=(.*)$/\1/p' )
expected="0.0.1"
echo "| 4.5 | ${expected} | ${actual} |"
if [ "${actual}" != "${expected}" ]; then
    err="15"
    echo "FAILED #4.5"
    echo "${out}"
    echo "==="
fi

# no previous tags, nothing in commits to trigger a bump
# default as patch, a prerelease so should generate that segement
out=$(python ./next-tag.py \
    --test_file=./tests/none.txt \
    --prerelease="true" \
    --prerelease_suffix="moreactions" \
    --latest_tag="" \
    --default_bump="patch" \
    --last_release="")
actual=$(echo "${out}" | sed -r -n 's/.*next_tag=(.*)$/\1/p' )
expected="0.0.1-moreactions.0"
echo "| 4.6 | ${expected} | ${actual} |"
if [ "${actual}" != "${expected}" ]; then
    err="16"
    echo "FAILED #4.6"
    echo "${out}"
    echo "==="
fi
# previous release, not a prerelease, with_v and major bump
out=$(python ./next-tag.py \
    --test_file=./tests/none.txt \
    --prerelease_suffix="moreactions" \
    --last_release="v1.1.0" \
    --default_bump="major" \
    --with_v="true")
actual=$(echo "${out}" | sed -r -n 's/.*next_tag=(.*)$/\1/p' )
expected="v2.0.0"
echo "| 4.7 | ${expected} | ${actual} |"
if [ "${actual}" != "${expected}" ]; then
    err="17"
    echo "FAILED #4.7"
    echo "${out}"
    echo "==="
fi


# # local test using git rather than test files
# python ./next-tag.py \
#     --prerelease="true" \
#     --prerelease_suffix="moreactions" \
#     --last_release="v1.4.0" \
#     --latest_tag="v1.5.0-moreactions.14" \
#     --with_v="true"

if [ -n "${err}" ]; then
    echo "err: ${err}"
    exit 1
fi