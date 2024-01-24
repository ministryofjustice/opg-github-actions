#!/usr/bin/env bash

ACTION_REPO="${GH_ACTION_REPOSITORY}"
ACTION_REF="${GH_ACTION_REF}"
BASE_PATH="${GH_WORKSPACE}/../opg-gha"
ARTIFACT_PATH="${BASE_PATH}/releases"
TARBALL="release.tar.gz"

os=$(uname | tr '[:upper:]' '[:lower:]')
arch=$(uname -m)
HOST_BUILD="${os}_${arch}"

ok="ok!"
# If github action repo or the action ref are empty, then fail with error
if [ -z "${ACTION_REF}" ] || [ -z "${ACTION_REPO}" ]; then 
    err="ERROR: this composite action must be run via full path (eg ministryofjustice/opg-github-actions/.github/actions/terraform-version@v2.3.1)"
    echo -e "${err}" >&2
    exit 1
fi

mkdir -p ${ARTIFACT_PATH}
# Try to download the release artifact directly, presuming 
# action_ref is a release tag
echo -n "Trying direct release download using [${ACTION_REF}] [${ACTION_REPO}]..."
releases=$(gh release list --exclude-drafts=false --exclude-pre-releases=false -R "${ACTION_REPO}")
# swallow the error as we want to move to hash lookup
DIRECT=$( echo "${releases}" | grep "^${ACTION_REF}" && echo "${ok}")

set -e
set -o pipefail

# If direct worked, then move the downloaded artifact
if [ "${DIRECT}" == "${ok}" ]; then
    gh release download "${ACTION_REF}" -R "${ACTION_REPO}" --clobber
    echo " ✅"
    mv *.tar.gz ${ARTIFACT_PATH}
fi
# If we failed to download the artifact using the action_ref directly
# then its likely someone has used a git hash ref to pass to us
# In that case, try to find a semver tag that points to that hash
if [ "${DIRECT}" != "${ok}" ]; then
    echo " ❌"
    echo -e "Releases found:"
    echo -e "${releases}"
    echo -e "Will try to convert [${ACTION_REF}] to a known release tag."
    echo -e "Cloning action repostitory [${ACTION_REPO}] locally..."
    REF=""

    gh repo clone ${ACTION_REPO} ${BASE_PATH} -- --mirror
    cd ${ARTIFACT_PATH}

    # look for semver tags at this ref
    tags=$(git tag --points-at="${ACTION_REF}")
    semverish=$(echo ${tags} | grep -E '^v[0-9]{1,}.[0-9]{1,}.[0-9]{1,}' )
    semverishCount=$(echo ${semverish} | wc -l | tr -d ' ')

    if [ "${semverishCount}" -ge "1" ]; then
        tag=$(echo ${semverish} | tail -n1 )
        echo -e "Using [${tag}] as release tag."
        REF=$(gh release download ${tag} --clobber && echo "${ok}")
    fi

    if [ "${REF}" != "${ok}" ]; then
        err="ERROR: could not find release artifact for [${ACTION_REF}]"
        echo -e "${err}" >&2
        exit 1
    fi
fi

# make sure we have an artifact that matches this runner
cd ${ARTIFACT_PATH}
ls -la
pwd
if [ ! -r "${TARBALL}" ]; then
    err="ERROR: could not find a readable tar ball at [${ARTIFACT_PATH}/${TARBALL}]"
fi
echo -e "Expanding tar ball [${TARBALL}]"
tar -xzvf ${TARBALL}

if [ ! -x "${HOST_BUILD}" ]; then
    err="ERROR: unable to find executable at [${ARTIFACT_PATH}/${HOST_BUILD}]"
    echo -e "${err}" >&2
    exit 1
fi

echo -e "release=${ARTIFACT_PATH}/${HOST_BUILD}"
echo "release=${ARTIFACT_PATH}/${HOST_BUILD}" >> $GITHUB_OUTPUT
