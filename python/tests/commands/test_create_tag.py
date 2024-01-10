#!/usr/bin/env python3
import pytest
from git import Repo
import shutil
from actions.commands import create_tag as cmd
from actions.common import outputhelper as oh

### logic
fh = open("./create_tag_results.md", "a+")
o = oh.OutputHelper(False)
o.header(fh, "Create Tag Test Results")
fh.close()

@pytest.fixture()
def setup_tags(request) -> tuple:
    fh = open("./create_tag_results.md", "a+")
    print("\nSetting up resources...")
    # clone the repo
    repo_root = "./create-tag-repo-test/"
    commitish = "dab1b63"
    url = "https://github.com/ministryofjustice/opg-github-actions.git"

    Repo.clone_from(url, repo_root)
    repo = Repo(repo_root)
    # create all the test tags locally
    test_tags = [
        'test_tag',
        'test_tag_wont_be_latest',
        'v2.4.9',
        'v2.4.15',
        'v1.0.1',
        'v20.10.10',
        'v20.10.9',
        'v21.0.1',
        'v22.0.0',
        'v23.0.0',
        'v1.5.0-clash.0',
        'v1.5.0-clash.1',
        'v9999.1.0'
    ]
    for t in test_tags:
        repo.git.tag(t, commitish)

    yield repo_root, commitish, fh

    print("\nPerforming teardown...")
    try:
        shutil.rmtree(repo_root)
        fh.close()
    except OSError as e:
        print("Error: %s - %s." % (e.filename, e.strerror))


def test_clashing_prerelease_tag_generates_new_tag(setup_tags) -> None:
    """
    We created v1.5.0-clash.0 v1.5.0-clash.1 already, so this
    should return a different tag as a prerelease the new tag
    should contain most of the original tag with the
    prerelease segment changing.
    """
    repo_root, commitish, fh = setup_tags
    pre="v1.5.0-"
    tag = f"{pre}clash.1"
    outputs = cmd.run( repo_root, commitish, tag, True )
    # the created tag should not match the request tag, as that already exists
    t1 = (tag != outputs['created_tag'])
    assert t1 == True
    # should have created a new tag based off the existing version
    t2 = (pre in outputs['created_tag'])
    assert t2 == True
    # output to gh
    o.result(tag, "!=", outputs['created_tag'], t1 == True, fh)
    o.result(pre, "in", outputs['created_tag'], t2 == True, fh)


def test_brand_new_tag_matches(setup_tags) -> None:
    """
    Brand new tag being created, so the returned value should
    be the same as requested
    """
    repo_root, commitish, fh = setup_tags
    tag = "v999.0.1"
    outputs = cmd.run( repo_root, commitish, tag, True )
    # the created tag should not match the request tag, as that already exists
    t1 = (tag == outputs['created_tag'])
    assert t1 == True
    o.result(tag, "==", outputs['created_tag'], t1 == True, fh)


def test_release_version_already_exists(setup_tags) -> None:
    """
    This release version already exists, so it should bump along
    the major number and reset the minor / patch.
    """
    repo_root, commitish, fh = setup_tags
    tag = "v9999.1.0"
    expected = "v10000.0.0"
    outputs = cmd.run( repo_root, commitish, tag, True )
    t1 = (expected == outputs['created_tag'])
    assert t1 == True
    o.result(expected, "==", outputs['created_tag'], t1 == True, fh)


def test_release_version_multiple_increaments(setup_tags) -> None:
    """
    This release version already exists, as do several others
    in this space, so the tag should jump ahead to an
    unused number.
    """
    repo_root, commitish, fh = setup_tags
    tag = "v21.0.1"
    expected = "v24.0.0"
    outputs = cmd.run( repo_root, commitish, tag, True )
    t1 = (expected == outputs['created_tag'])
    assert t1 == True
    o.result(expected, "==", outputs['created_tag'], t1 == True, fh)


def test_non_semver_tag_match(setup_tags) -> None:
    """
    This tag exists, but is not a semver tag, so the output
    should contain the original tag plus a random string
    """
    repo_root, commitish, fh = setup_tags
    tag = "test_tag"
    match = f"{tag}."
    outputs = cmd.run( repo_root, commitish, tag, True )
    t1 = (match in outputs['created_tag'])
    assert t1 == True
    o.result(match, "in", outputs['created_tag'], t1 == True, fh)
