#!/usr/bin/env bash
set -e

header(){
    echo "## Test Information"
    echo "| \# | &nbsp; | &nbsp; | &nbsp; | &nbsp; |"
    echo "| --- | --- | --- | --- | --- |"
}


fail() {
    echo "| ${1} | ${2} | ${3} | ${4} | ❌ |"
}

pass() {
    echo "| ${1} | ${2} | ${3} | ${4} | ✅ |"
}