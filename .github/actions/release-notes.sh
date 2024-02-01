#!/usr/bin/env bash
target="${GH_REPO}"
commit="${GH_COMMIT}"
previous="${LAST_TAG}"
tag="${TAGNAME}"
debug="${DEBUG}"
maxLength=124000

accept="Accept: application/vnd.github+json"
ver="X-GitHub-Api-Version: 2022-11-28"
endpoint="/repos/${target}/releases/generate-notes"
tagParam="-f tag_name=${tag} -f previous_tag_name=${previous} "
commitParam="-f target_commitish=${commit} "

echo -n "Generating release notes..."
generatedNotes=$(gh api --method POST -H "${accept}" -H "${ver}" ${endpoint} ${tagParam} ${commitParam} 2>/dev/null )
genLen=${#generatedNotes}

if [ "${genLen}" -le "0" ]; then
    echo " ❌"
    echo -e "Generate notes api call failed"
    exit 1
fi

echo " ✅"
body=$(echo ${generatedNotes} | jq ".body")
len=${#body}
echo -n "Generated release note body is [${len}] characters"

if [ "${len}" -ge "${maxLength}" ]; then
    echo " ❌"
    echo -e "Truncating..."    
    body=${body:0:${maxLength}}
    len=${#body}
    echo -e "New length [${len}] ✅"
else 
    echo " ✅"
fi

if [ "${debug}" == "1" ]; then
    echo -e "==="
    echo -e "${body}"
    echo -e "==="
fi
# export variables back for use in workflow
export RELEASE_BODY=${body}
echo "RELEASE_BODY=${RELEASE_NOTE_BODY}" >> $GITHUB_OUTPUT
