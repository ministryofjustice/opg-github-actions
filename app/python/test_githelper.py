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
git_mod = importlib.util.spec_from_file_location("githelper", dir_name + '/githelper.py')
ghm = importlib.util.module_from_spec(git_mod)  
git_mod.loader.exec_module(ghm)


@pytest.fixture()
def setup_repo(request) -> tuple:    
    print("\nSetting up resources...")   
    # download repo
    repo_root = "./githelper-test/"
    url = "https://github.com/ministryofjustice/opg-github-actions.git"
    Repo.clone_from(url, repo_root)
    # add dummy tags for tests
    commitish = "dab1b63"
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
        'v9999.1.0',
        '5.0.1-test.0',
        '5.0.1'
    ]
    for t in test_tags:        
        repo.git.tag(t, commitish)

    yield repo_root, commitish, test_tags

    print("\nPerforming teardown...")
    try:
        shutil.rmtree(repo_root)
    except OSError as e:
        print("Error: %s - %s." % (e.filename, e.strerror))

@pytest.fixture()
def setup_tags(request) -> tuple:
    print("\nSetting up resources...")   
    
    semver_non_clashing_tags = {
        'v1.6.0-': 'v1.6.0-beta.0+rb1',
        'v10.0.1': 'v10.0.1',
        '11.0.1-': '11.0.1-beta.1'
    }
    semver_clashing_tags = {
        'v1.5.0-': 'v1.5.0-clash.0',
        '5.0.1': '5.0.1-test.0'
    }
    non_clashing_tags = {
        'just_a_tag': 'just_a_tag',
        'a-tag-to-create': 'a-tag-to-create'
    }
    clashing_tags = {
        'test_tag':'test_tag'
    }
    
    yield semver_non_clashing_tags, semver_clashing_tags, non_clashing_tags, clashing_tags
    print("\nPerforming teardown...")    


def test_tags_at_point(setup_repo) -> None:
    """
    Test tag listing is being generated correctly and all test tests
    are found within the list of tags at a set point
    """
    path, commitish, test_tags = setup_repo
    r = ghm.GitHelper(path)

    tags = r.tags(f"--points-at={commitish}")    
    # find only the test tags
    found_tests = list (filter(lambda x: (x in test_tags), tags))
    assert (len(found_tests) > 0) == True
    assert (len(found_tests) == len(test_tags)) == True

def test_tags_full(setup_repo) -> None:
    """
    Test tag listing is being generated correctly and all test tests
    are found within the full list of tags
    """
    path, commitish, test_tags = setup_repo
    r = ghm.GitHelper(path)

    tags = r.tags("--list")    
    # find only the test tags
    found_tests = list (filter(lambda x: (x in test_tags), tags))
    assert (len(found_tests) > 0) == True
    assert (len(found_tests) == len(test_tags)) == True

def test_tag_to_create(setup_repo, setup_tags) -> None:
    """
    Test the logical loops in the create tag to make sure clean ones 
    work and a clash is dealt with correctly by increamenting
    """
    path, commitish, test_tags = setup_repo
    semver_non_clashing_tags, semver_clashing_tags, non_clashing_tags, clashing_tags = setup_tags
    r = ghm.GitHelper(path)

    for t in (semver_non_clashing_tags | non_clashing_tags).values():
        tag = r.tag_to_create(t, test_tags)
        assert (tag == t) == True

    # clashing tags should not match, but should contain a segment
    # of the original and maintain prefix
    for k,t in (semver_clashing_tags | clashing_tags).items():
        tag = r.tag_to_create(t, test_tags)
        assert (tag != t) == True
        assert (k in tag) == True


def test_tag_creation_works(setup_repo) -> None:
    """
    Test that a created tag is in the list of tags
    """
    path, commitish, test_tags = setup_repo
    r = ghm.GitHelper(path)

    # create the tag
    tag = "v10000.0.0"
    r.create_tag(tag, commitish, False)
    # this should now exist in the list of tags at the location
    created = (tag in r.tags(f"--points-at={commitish}"))
    assert created == True

def test_tag_creation_fails(setup_repo) -> None:
    """
    Test that create tag fails as the tag already exists
    and an exception is thrown
    """
    path, commitish, test_tags = setup_repo
    r = ghm.GitHelper(path)
    # trying to create the tag should trigger exceptions
    tag = "test_tag"
    with pytest.raises(Exception) as err:
        r.create_tag(tag, commitish, False)
    try:
        r.create_tag(tag, commitish, False)
        assert False
    except Exception:
        assert True

    

