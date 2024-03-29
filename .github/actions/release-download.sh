#!/usr/bin/env bash
actionRepo="${GH_ACTION_REPOSITORY}"
actionRef="${GH_ACTION_REF}"
basePath="${GH_WORKSPACE}/../opg-gha"
localBuildPath="${GH_WORKSPACE}/../opg-gha-build"
artifactPath="${basePath}/releases"
tarball="release.tar.gz"
os=$(uname | tr '[:upper:]' '[:lower:]')
arch=$(uname -m)

if [ -z "${hostBuild}" ]; then    
    hostBuild="${os}_${arch}"
fi

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
releases=$(gh release list --exclude-drafts=false --exclude-pre-releases=false -R "${actionRepo}" 2>/dev/null)
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
    gh release download "${actionRef}" -R "${actionRepo}" --clobber 2>/dev/null
    # move them
    echo -e "Moving tarball to [${artifactPath}]"
    mv *.tar.gz ${artifactPath}
    # expand the tarball
    cd ${artifactPath}
    echo -e "Expanding tarball [${artifactPath}/${tarball}]"
    tar -xzvf ${tarball}
    # look for this arch
    echo -n "Looking for binary for this runner [${artifactPath}/${hostBuild}]..."
    # make sure the binary is an executable
    if [ -x "${hostBuild}" ]; then
        echo " ✅"
        RELEASE="${artifactPath}/${hostBuild}"
        echo -e "Set release: [${RELEASE}]"
    # if we fail to get an executable bin then we reset listed to empty
    # so it will trigger the download and self build block below
    else
        echo " ❌"
        echo -e "Failed to find binary for this runner, will trigger the self build... "
        listed=""
    fi
else
    echo " ❌"
    echo -e "\nAvailable releases: "
    echo -e "${releases}\n"
fi
# If we failed to download the artifact using the action_ref directly
# then its likely someone has used a prerelease or a git hash ref 
# to pass to us 
# In that case, we need to download the repo to a new location so
# that it can be built
if [ "${found}" != "${ok}" ]; then
    # build from local 
    echo -n "Cloning action repostitory [${actionRepo}] to [${localBuildPath}] ..."
    # if the target directory exists, remove it before cloning
    if [ -d "${localBuildPath}" ]; then
        rm -Rf ${localBuildPath}
    fi
    # use the gh cli to clone so we dont have to work out the url path
    gh repo clone ${actionRepo} ${localBuildPath} -- -q 2>/dev/null
    echo " ✅"

    cd ${localBuildPath}

    # If this is a pr, then use the gh cli to checkout instead
    # of mapping refs 
    if [[ "${actionRef}" == "refs/pull"* && "${actionRef}" == *"/merge" ]]; then        
        echo -e "This seems to be a pull request: [${actionRef}]"
        prNumber=$(echo "${actionRef}" | tr -cd '[:digit:]')
        echo -e "pr number: ${prNumber}"
        pr=$(gh pr checkout "${prNumber}" && echo "${ok}" )
        checkout=$(echo "${pr}" | tail -n1)
    # otherwise checkout directly
    else     
        checkout=$(git checkout -q -f ${actionRef} -- 2> /dev/null && echo "${ok}")
    fi

    # if the checkout worked, then all ok, otherwise we exit with a code 1
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

# export the import variables and push them to the github output as well
export RELEASE=${RELEASE}
export SELF_BUILD=${SELF_BUILD}
export TARGET_BUILD=${TARGET_BUILD}

echo -e "SELF_BUILD=${SELF_BUILD}"
echo -e "RELEASE=${RELEASE}"
echo -e "TARGET_BUILD=${TARGET_BUILD}"

echo "SELF_BUILD=${SELF_BUILD}" >> $GITHUB_OUTPUT
echo "RELEASE=${RELEASE}" >> $GITHUB_OUTPUT
echo "TARGET_BUILD=${TARGET_BUILD}" >> $GITHUB_OUTPUT