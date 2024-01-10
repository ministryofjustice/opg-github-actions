#!/usr/bin/env python3
import pytest
from semver.version import Version
from actions.common import semverhelper as svh


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
        d = svh.to_dict(t)
        assert (d[t] is not None) == True
    for t in invalid:
        d = svh.to_dict(t)
        assert (d[t] is None) == True


def test_list_to_dict(setup_valid_invalid) -> None:
    """
    Test that all valid tags are returned within the
    dict and invalid ones are discounted
    """
    valid, invalid = setup_valid_invalid
    # check that all valid items are returned, but invalid are not
    all = svh.list_to_dict(valid + invalid)
    real = {k:v for k, v in all.items() if v is not None}
    assert (len(real) == len(valid) ) == True

    all = svh.list_to_dict(invalid)
    real = {k:v for k, v in all.items() if v is not None}
    assert (len(real) == 0 ) == True


@pytest.mark.parametrize(
    "expected,branch_name,release_branches,prerelease",
    [
        (True, "test", "main", True),
        (True, "v1.5.0-beta.0", "main,master", True),
        (False, "main", "main,master", True)
    ]
)
def test_is_prerelease(
    expected:bool,
    branch_name:str,
    release_branches:str|list,
    prerelease:bool|str,
) -> None:
    """
    Check the logic of is_prerelese matches expected
    """
    actual = svh.is_prerelease(
        branch_name=branch_name,
        release_branches=release_branches,
        stated_prerelease_state=prerelease
    )
    print(f"branch_name={branch_name}\nrelease_branches={release_branches}\nstated_prerelease_state={prerelease}")
    print(f"Expected [{expected}] Actual [{actual}]")
    assert (expected == actual) == True


@pytest.mark.parametrize(
    "expected,search,tag_list",
    [
        (1, "dummy", ["v1.5.0-dummy.0", "1.2.0-test.0", "test_tag", "1.0.0"]),
        (2, "beta", ["v1.5.0-beta.0", "1.2.0-beta.0", "test_tag", "1.0.0"]),
        (0, "beta", ["v1.5.0-dummy.0", "1.2.0-test.0", "test_tag", "1.0.0"]),
    ]
)
def test_prereleases_filtered(expected:int, search:str, tag_list:list) -> None:
    """
    Check various tag lists to make sure the prerelease style filter
    find the correct amount
    """
    filter = f"{search}.[0-9]+$"
    found = svh.prereleases_filtered(tag_list, filter)
    actual = len(found)

    assert (expected == actual) == True


@pytest.mark.parametrize(
    "expected,tag_list",
    [
        (2, ["v1.5.0-dummy.0", "1.2.0-test.0", "test_tag", "1.0.0"]),
        (1, ["v1.5.0-beta.0+b1", "1.2.0", "test_tag", "1.0.0"]),
        (0, ["v1.5.0", "1.2.0", "test_tag", "1.0.0", "weirdprefix1.2.0-test.0"]),
    ]
)
def test_prereleases(expected:int, tag_list:list) -> None:
    """
    Check prereleases finds correct tags out of a set
    """
    found = svh.prereleases(tag_list)
    actual = len(found)

    assert (expected == actual) == True

@pytest.mark.parametrize(
    "expected,tag_list",
    [
        (1, ["v1.5.0-dummy.0", "1.2.0-test.0", "test_tag", "1.0.0"]),
        (2, ["v1.5.0-beta.0+b1", "1.2.0", "test_tag", "1.0.0"]),
        (3, ["v1.5.0", "1.2.0", "test_tag", "1.0.0", "weirdprefix1.2.0-test.0"]),
        (0, ["v1.5.0-b", "1.2.0-beta.0", "test_tag", "not-a-relealse-1.0.0", "weirdprefix1.2.0-test.0"]),
    ]
)
def test_releases(expected:int, tag_list:list) -> None:
    """
    Check releases finds correct tags out of a set
    """
    found = svh.releases(tag_list)
    actual = len(found)

    assert (expected == actual) == True


@pytest.mark.parametrize(
    "expected,tag_list",
    [
        ("1.0.0-beta.10", ["1.0.0-beta.9", "1.0.0-beta.10"]),
        ("v1.0.0-beta.10", ["v1.0.0-beta.9", "v1.0.0-beta.10"]),
        # If there is a release and pre-release that match, release is used
        ("100.0.0", ["v1.0.0-beta.9", "v1.0.0-beta.10", "100.0.0", "100.0.0-test.0"]),
        ("100.5.0-test.0", ["v1.0.0-beta.9", "v1.0.0-beta.10", "100.1.0", "100.5.0-test.0"]),
    ]
)
def test_max(expected:str, tag_list:list) -> None:
    """
    Check releases finds correct tags out of a set
    """
    tags:dict = svh.list_to_dict(tag_list)
    actual = svh.max_version(tags)

    assert (expected == actual) == True

### SETUP NEXT TAG GENERATION TESTS
testconfig = [
    {
        "expected": "2.0.0",
        "major_bump": 1,
        "minor_bump": 1,
        "patch_bump": 1,
        "is_prerelease": False,
        "prerelease_suffix": "beta",
        "latest_tag": "2.0.0-beta.0",
        "last_release": "1.0.1"
    },
    {
        "expected": "1.1.0",
        "major_bump": 0,
        "minor_bump": 1,
        "patch_bump": 0,
        "is_prerelease": False,
        "prerelease_suffix": "beta",
        "latest_tag": "1.1.0-beta.0",
        "last_release": "1.0.0"
    },
    {
        "expected": "1.1.0-beta.1",
        "major_bump": 0,
        "minor_bump": 1,
        "patch_bump": 0,
        "is_prerelease": True,
        "prerelease_suffix": "beta",
        "latest_tag": "1.1.0-beta.0",
        "last_release": "1.0.0"
    },
    {
        "expected": "0.1.0-beta.0",
        "major_bump": 0,
        "minor_bump": 1,
        "patch_bump": 0,
        "is_prerelease": True,
        "prerelease_suffix": "beta",
        "latest_tag": None,
        "last_release": None
    },
    {
        "expected": "1.0.0",
        "major_bump": 1,
        "minor_bump": 0,
        "patch_bump": 0,
        "is_prerelease": False,
        "prerelease_suffix": "beta",
        "latest_tag": None,
        "last_release": None
    },
]
# generate fields string from the keys
fields = ','.join(testconfig[0].keys())
# generate test tuple from config items
tests = [(v.values()) for v in testconfig]
@pytest.mark.parametrize(fields, tests)
def test_next_tag(
    expected:str|Version,
    major_bump:int,
    minor_bump:int,
    patch_bump:int,
    is_prerelease:bool,
    prerelease_suffix:str,
    latest_tag:str,
    last_release:str
) -> None:
    """
    Test that the generated tag matches what we would expect
    """
    latest_tag:Version = Version.parse(latest_tag) if type(latest_tag) is str else latest_tag
    last_release:Version = Version.parse(last_release) if type(last_release) is str else last_release

    actual = svh.next_tag(
        major_bump=major_bump,
        minor_bump=minor_bump,
        patch_bump=patch_bump,
        is_prerelease=is_prerelease,
        prerelease_suffix=prerelease_suffix,
        latest_tag=latest_tag,
        last_release=last_release
    )
    assert (f"{expected}" == f"{actual}")
