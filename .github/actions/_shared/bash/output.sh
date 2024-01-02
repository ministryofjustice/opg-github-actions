#!/usr/bin/env bash
set -e

header(){
    echo "## Test Information"    
    echo "| \# | A | condition | B | Pass |"
    echo "| --- | --- | --- | --- | --- |"
}


fail() {
    echo "| ${1} | ${2} | ${3} | ${4} | ❌ |"
    TEST_ERR="true"
}

pass() {
    echo "| ${1} | ${2} | ${3} | ${4} | ✅ |"
}