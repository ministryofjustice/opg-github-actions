#!/usr/bin/env bash
set -e
# remember to include output as well
# source $(dirname $0)/output.sh

# get the value from the output (var_name, output)
test_output_value() {
    val=$(echo "${2}" | sed -r -n "s/.*${1}=(.*)$/\1/p" )
    echo -e "${val}"
}


# if it DOES NOT match, its a failure
test_should_equal() {
    num="${1}"
    expected="${2}"
    actual="${3}"
    
    if [ "${actual}" != "${expected}" ]; then
        fail "${num}" "${expected}" "==" "${actual}"
    else
        pass "${num}" "${expected}" "==" "${actual}"
    fi            
}


# if the needle is in the haystack, then true
test_should_contain() {
    num="${1}"
    needle="${2}"
    haystack="${3}"
    
    if [ "${haystack}" == *"${needle}*" ]; then
        pass "${num}" "${haystack}" "contains" "${needle}"
    else
        pass "${num}" "${haystack}" "contains" "${needle}"
    fi            
}

# if it DOES match, its a failure
test_should_not_equal() {    
    num="${1}"
    expected="${2}"
    actual="${3}"
    
    if [ "${actual}" == "${expected}" ]; then
        fail "${num}" "${expected}" "!=" "${actual}"
    else
        pass "${num}" "${expected}" "!=" "${actual}"
    fi   

}