#!/usr/bin/env bash
actionRepo="${GH_ACTION_REPOSITORY}"
actionRef="${GH_ACTION_REF}"
basePath="${GH_WORKSPACE}/../opg-gha"
localBuildPath="${GH_WORKSPACE}/../opg-gha-build"
artifactPath="${basePath}/releases"
tarball="release.tar.gz"
os=$(uname | tr '[:upper:]' '[:lower:]')
arch=$(uname -m)
hostBuild="${os}_${arch}"
ok=">ok<"

####
# These vars will be passed along to determine next action
RELEASE=""
SELF_BUILD=""
TARGET_BUILD="${hostBuild}"
####

# Look for the reference in the existing releases
# - if we find that, we'll use that directly
echo -n "Trying direct release download using [${actionRef}] [${actionRepo}]..."
releases=$(gh release list --exclude-drafts=false --exclude-pre-releases=false -R "${actionRepo}")
listed=$( echo "${releases}" | grep "^${actionRef}" && echo "${ok}")
found=$(echo "${listed}" | tail -n1)

set -e
set -o pipefail
# If ref is in the release list, then we can download
# and then move the artifact
if [ "${found}" == "${ok}" ]; then
    echo " ✅"    
    mkdir -p ${basePath}
    mkdir -p ${artifactPath}
    cd ${basePath}
    echo -e "Downloading existing release [${actionRef}]..."
    # download the release tar ball
    gh release download "${actionRef}" -R "${actionRepo}" --clobber
    # move them
    echo -e "Moving tarball to [${artifactPath}]"
    mv *.tar.gz ${artifactPath}
    # expand the tarball
    cd ${artifactPath}
    echo -e "Expanding tarball [${artifactPath}/${tarball}]"
    tar -xzvf ${tarball}
    # look for this arch
    echo -n "Looking for binary for this runner [${artifactPath}/${hostBuild}]..."
    if [ -x "${hostBuild}" ]; then
        echo " ✅"
        RELEASE="${artifactPath}/${hostBuild}"
        echo -e "Set release: [${RELEASE}]"
    else
        echo " ❌"
        echo -e "Failed to find binary for this runner, will trigger the self build... "
        listed=""
    fi
else
    echo " ❌"
    echo -e "Releases: "
    echo -e "${releases}"
fi
# If we failed to download the artifact using the action_ref directly
# then its likely someone has used a prerelease or a git hash ref 
# to pass to us 
# In that case, we need to download the repo to a new location so
# that it can be built
if [ "${found}" != "${ok}" ]; then
    # build from local 
    echo -e "Cloning action repostitory [${actionRepo}] to [${localBuildPath}] ..."

    if [ -d "${localBuildPath}" ]; then
        rm -Rf ${localBuildPath}
    fi
    gh repo clone ${actionRepo} ${localBuildPath} -- -q
    
    cd ${localBuildPath}
    # checkout to the git ref
    checkout=$(git checkout -q -f ${actionRef} -- 2> /dev/null && echo "${ok}")
    if [ "${checkout}" == "${ok}" ]; then
        echo -e "Checked out action repo to [${actionRef}] [${localBuildPath}] ✅"
        echo -e "-- Commit --"
        git log -n1 --format="oneline"
        echo -e "------------"
        SELF_BUILD="${localBuildPath}"
    else
        err="ERROR: failed to checkout [${actionRef}]"
        echo -e "${err}"
        echo -e "${err}" >&2
        exit 1
    fi

  
fi

export RELEASE=${RELEASE}
export SELF_BUILD=${SELF_BUILD}
export TARGET_BUILD=${TARGET_BUILD}

echo "SELF_BUILD=${SELF_BUILD}" >> $GITHUB_OUTPUT
echo "RELEASE=${RELEASE}" >> $GITHUB_OUTPUT
echo "TARGET_BUILD=${TARGET_BUILD}" >> $GITHUB_OUTPUT