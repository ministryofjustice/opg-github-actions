# #!/usr/bin/env bash
echo "MAJORS"
set -e

out=$(python ./next-tag.py \
    --test_file=./tests/majors.txt \
    --prerelease_suffix="moreactions" \
    --prerelease=true \
    --latest_tag="v1.5.0-moreactions.1" \
    --last_release="v1.4.0")

actual=$(echo "${out}" | sed -r -n 's/.*next_tag=(.*)$/\1/p' )
expected="2.0.0-moreactions.0"
echo "#1.1 [${expected}]=>[${actual}]"
if [ "${actual}" != "${expected}" ]; then
    echo "FAILED #1.1"
    echo "${out}"
    echo "==="
fi

out=$(python ./next-tag.py \
    --test_file=./tests/majors.txt \
    --prerelease_suffix="moreactions" \
    --prerelease=true \
    --latest_tag="v2.0.0-moreactions.0" \
    --last_release="v1.4.0")

actual=$(echo "${out}" | sed -r -n 's/.*next_tag=(.*)$/\1/p' )
expected="2.0.0-moreactions.1"
echo "#1.2 [${expected}]=>[${actual}]"
if [ "${actual}" != "${expected}" ]; then
    echo "FAILED #1.2"
    echo "${out}"
    echo "==="
fi


out=$(python ./next-tag.py \
    --test_file=./tests/majors.txt \
    --prerelease_suffix="moreactions" \
    --prerelease=true \
    --latest_tag="" \
    --last_release="v1.4.0")
actual=$(echo "${out}" | sed -r -n 's/.*next_tag=(.*)$/\1/p' )
expected="2.0.0-moreactions.0"
echo "#1.3 [${expected}]=>[${actual}]"
if [ "${actual}" != "${expected}" ]; then
    echo "FAILED #1.3"
    echo "${out}"
    echo "==="
fi


out=$(python ./next-tag.py \
    --test_file=./tests/majors.txt \
    --prerelease_suffix="moreactions" \
    --latest_tag="" \
    --last_release="v1.4.0")
actual=$(echo "${out}" | sed -r -n 's/.*next_tag=(.*)$/\1/p' )
expected="2.0.0"
echo "#1.4 [${expected}]=>[${actual}]"
if [ "${actual}" != "${expected}" ]; then
    echo "FAILED #1.4"
    echo "${out}"
    echo "==="
fi


out=$(python ./next-tag.py \
    --test_file=./tests/majors.txt \
    --prerelease_suffix="moreactions" \
    --latest_tag="" \
    --last_release="")
actual=$(echo "${out}" | sed -r -n 's/.*next_tag=(.*)$/\1/p' )
expected="1.0.0"
echo "#1.5 [${expected}]=>[${actual}]"
if [ "${actual}" != "${expected}" ]; then
    echo "FAILED #1.5"
    echo "${out}"
    echo "==="
fi

echo "MINORS"
out=$(python ./next-tag.py \
    --test_file=./tests/minors.txt \
    --prerelease_suffix="moreactions" \
    --prerelease=true \
    --latest_tag="v1.5.0-moreactions.0" \
    --last_release="v1.4.0")
actual=$(echo "${out}" | sed -r -n 's/.*next_tag=(.*)$/\1/p' )
expected="1.5.0-moreactions.1"
echo "#2.1 [${expected}]=>[${actual}]"
if [ "${actual}" != "${expected}" ]; then
    echo "FAILED #2.1"
    echo "${out}"
    echo "==="
fi

out=$(python ./next-tag.py \
    --test_file=./tests/minors.txt \
    --prerelease_suffix="moreactions" \
    --prerelease=true \
    --latest_tag="" \
    --last_release="v1.4.0")
actual=$(echo "${out}" | sed -r -n 's/.*next_tag=(.*)$/\1/p' )
expected="1.5.0-moreactions.0"
echo "#2.2 [${expected}]=>[${actual}]"
if [ "${actual}" != "${expected}" ]; then
    echo "FAILED #2.2"
    echo "${out}"
    echo "==="
fi


out=$(python ./next-tag.py \
    --test_file=./tests/minors.txt \
    --prerelease_suffix="moreactions" \
    --latest_tag="v1.5.0-moreactions.0" \
    --last_release="v1.4.0")
actual=$(echo "${out}" | sed -r -n 's/.*next_tag=(.*)$/\1/p' )
expected="1.5.0"
echo "#2.3 [${expected}]=>[${actual}]"
if [ "${actual}" != "${expected}" ]; then
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
echo "#2.4 [${expected}]=>[${actual}]"
if [ "${actual}" != "${expected}" ]; then
    echo "FAILED #2.4"
    echo "${out}"
    echo "==="
fi


echo "PATCH"
out=$(python ./next-tag.py \
    --test_file=./tests/patch.txt \
    --prerelease=true \
    --prerelease_suffix="moreactions" \
    --latest_tag="" \
    --last_release="v1.4.0")
actual=$(echo "${out}" | sed -r -n 's/.*next_tag=(.*)$/\1/p' )
expected="1.4.1-moreactions.0"
echo "#3.1 [${expected}]=>[${actual}]"
if [ "${actual}" != "${expected}" ]; then
    echo "FAILED #3.1"
    echo "${out}"
    echo "==="
fi

out=$(python ./next-tag.py \
    --test_file=./tests/patch.txt \
    --prerelease=true \
    --prerelease_suffix="moreactions" \
    --latest_tag="v1.4.1-moreactions.1" \
    --last_release="v1.4.0")
actual=$(echo "${out}" | sed -r -n 's/.*next_tag=(.*)$/\1/p' )
expected="1.4.1-moreactions.2"
echo "#3.2 [${expected}]=>[${actual}]"
if [ "${actual}" != "${expected}" ]; then
    echo "FAILED #3.2"
    echo "${out}"
    echo "==="
fi


out=$(python ./next-tag.py \
    --test_file=./tests/patch.txt \
    --prerelease_suffix="moreactions" \
    --latest_tag="v1.4.1-moreactions.1" \
    --last_release="v1.4.0")
actual=$(echo "${out}" | sed -r -n 's/.*next_tag=(.*)$/\1/p' )
expected="1.4.1"
echo "#3.3 [${expected}]=>[${actual}]"
if [ "${actual}" != "${expected}" ]; then
    echo "FAILED #3.3"
    echo "${out}"
    echo "==="
fi

out=$(python ./next-tag.py \
    --test_file=./tests/patch.txt \
    --prerelease_suffix="moreactions" \
    --latest_tag="" \
    --last_release="")    
actual=$(echo "${out}" | sed -r -n 's/.*next_tag=(.*)$/\1/p' )
expected="0.0.1"
echo "#3.4 [${expected}]=>[${actual}]"
if [ "${actual}" != "${expected}" ]; then
    echo "FAILED #3.4"
    echo "${out}"
    echo "==="
fi
