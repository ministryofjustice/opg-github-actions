#!/usr/bin/env bash
set -e

header(){
    echo "## Test Information"
    echo "| \# | &nbsp; | &nbsp; | &nbsp; | &nbsp; |"
    echo "| --- | A | condition | B | Pass / Fail |"
}


fail() {
    echo "| ${1} | ${2} | ${3} | ${4} | ❌ |"
}

pass() {
    echo "| ${1} | ${2} | ${3} | ${4} | ✅ |"
}