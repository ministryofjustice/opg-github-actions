#!/usr/bin/env python3
from semver.version import Version
import os
import importlib.util
from git import Repo, Git
import pytest
import shutil

# CUSTOM PATH LOADING
dir_name = os.path.dirname(os.path.realpath(__file__))
# load semver helpers
semver_mod = importlib.util.spec_from_file_location("semverhelper", dir_name + '/semverhelper.py')
svh = importlib.util.module_from_spec(semver_mod)  
semver_mod.loader.exec_module(svh)


@pytest.fixture()
def setup_with_without_v(request) -> tuple:    
    print("\nSetting up resources...")   
    with_prefix = ["v1.2.3-beta.0", "v1.2.3-beta.0+b1", "v1.0.0"] 
    without_prefix = ["1.2.3-beta.0", "1.2.3-beta.0+b1", "1.0.0"]
    
    yield with_prefix, without_prefix

    print("\nPerforming teardown...")

@pytest.fixture()
def setup_valid_invalid(request) -> tuple:    
    print("\nSetting up resources...")   
    valid = ["v1.2.3-beta.0", "v1.2.3-beta.0+b1", "v1.0.0"] 
    invalid = ["1.2beta.0", "plain_text_tag", "v1beta-0.2.1+b1", "more-plain-test"]
    
    yield valid, invalid
    print("\nPerforming teardown...")


def test_has_prefix(setup_with_without_v) -> None:
    """
    Test how the prefix flag is being handled with various
    tags
    """
    with_v, without_v = setup_with_without_v
    # these should be true
    for t in with_v:              
        assert svh.SemverHelper(t).has_prefix() == True
    # should be false
    for t in without_v:        
        assert svh.SemverHelper(t).has_prefix() == False
    # will currently return true, but shouldn't really
    # FIX
    test_tags = ["verylongprefix1.2.3"]
    for t in test_tags:        
        assert svh.SemverHelper(t).has_prefix() == True

def test_without_prefix(setup_with_without_v) -> None:
    """
    Test how the prefix is being trimmed from the original tag
    """
    with_v, without_v = setup_with_without_v
    # returned value should match a trimmed version of the original    
    for t in with_v:              
        assert (svh.SemverHelper(t).without_prefix() == t[1:]) == True
    # returned value should be an exact match
    for t in without_v:              
        assert (svh.SemverHelper(t).without_prefix() == t) == True

def test_valid(setup_valid_invalid) -> None:
    """
    Check that known good and bad tags returned the 
    expected result
    """
    valid, invalid = setup_valid_invalid
    # should all valid
    for t in valid:              
        assert (svh.SemverHelper(t).valid() == True) == True
    # should all be invalid
    for t in invalid:              
        assert (svh.SemverHelper(t).valid() == False) == True

def test_parse(setup_valid_invalid) -> None:
    """
    Check that valid tags return a Version class and 
    invalid ones return none type
    """
    valid, invalid = setup_valid_invalid
    for t in valid:              
        assert (isinstance(svh.SemverHelper(t).parse(), Version) == True) == True
    # check invalid doesnt match Version, but is None
    for t in invalid:              
        assert (isinstance(svh.SemverHelper(t).parse(), Version) == False) == True
        assert (svh.SemverHelper(t).parse() == None) == True

def test_parse_uses_default() -> None:    
    """
    Check that call to parse with a bad tag and a default
    returns the default
    """
    bad = "test-tag-not-real"
    # use the __str__ method to check
    assert ( f"{svh.SemverHelper(bad).parse('0.0.0') }" == "0.0.0") == True
    # check against the return
    tester = svh.SemverHelper(bad).parse('0.0.0')
    dummy = Version.parse('0.0.0')
    assert (tester == dummy) == True

def test_parsed_changes_with_update(setup_valid_invalid) -> None:
    """
    Test how the updating of _parsed class attr is 
    correct
    """
    valid, invalid = setup_valid_invalid
    for t in valid:
        s1 = svh.SemverHelper(t)
        p1 = s1.parsed()
        assert (type(p1) is Version) == True        
        # now lets bump versions along
        p2 = p1.bump_minor()
        s1.update( p2 )
        p3 = s1.parsed()
        assert (type(p3) is Version) == True
        assert (p2 == p3) == True
        assert (p3 != p1) == True

def test_update(setup_valid_invalid) -> None:
    """
    Test to check the update handles version and
    string params and ends up with correct result
    via a parsed()
    """
    valid, invalid = setup_valid_invalid
    for t in valid:
        # test version bumping as a version and string
        ms = svh.SemverHelper(t)        
        mo = ms.parsed()        
        major = mo.bump_major()
        ms.update(major)
        assert (ms.parsed() == major) == True
        ms.update(f"{major}")
        assert (ms.parsed() == major) == True


def test_tag_matches_prefix(setup_with_without_v) -> None:
    """
    Check that tag returns string that matches the original -
    as in if it had a prefix its returned with one
    """
    with_v, without_v = setup_with_without_v
    for t in with_v:
        assert (svh.SemverHelper(t).tag() == t) == True
    for t in without_v:
        assert (svh.SemverHelper(t).tag() == t) == True


def test_to_dict(setup_valid_invalid) -> None:
    """
    Test that valid tags generate a correct dict
    and invalid ones return None
    """
    valid, invalid = setup_valid_invalid
    for t in valid:
        d = svh.SemverHelper.to_dict(t)
        assert (d[t] is not None) == True
    for t in invalid:
        d = svh.SemverHelper.to_dict(t)
        assert (d[t] is None) == True


def test_list_to_dict(setup_valid_invalid) -> None:
    """
    Test that all valid tags are returned within the 
    dict and invalid ones are discounted
    """
    valid, invalid = setup_valid_invalid
    # check that all valid items are returned, but invalid are not
    all = svh.SemverHelper.list_to_dict(valid + invalid)
    real = {k:v for k, v in all.items() if v is not None}
    assert (len(real) == len(valid) ) == True

    all = svh.SemverHelper.list_to_dict(invalid)
    real = {k:v for k, v in all.items() if v is not None}
    assert (len(real) == 0 ) == True
