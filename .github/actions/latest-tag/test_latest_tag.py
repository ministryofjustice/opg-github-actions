import os
import importlib.util
from git import Repo, Git
import pytest
import shutil

### local imports
app_root_dir = os.path.dirname(
    os.path.dirname(
        os.path.dirname(
            os.path.dirname(os.path.realpath(__file__))
        )
    )
)
dir_name = os.path.dirname(os.path.realpath(__file__))

# load cmd helper
mod = importlib.util.spec_from_file_location("latest-tag", dir_name + '/latest-tag.py')
cmd = importlib.util.module_from_spec(mod)
mod.loader.exec_module(cmd)
# load output helper
ohmod = importlib.util.spec_from_file_location("gh", app_root_dir + '/app/python/outputhelper.py')
oh = importlib.util.module_from_spec(ohmod)
ohmod.loader.exec_module(oh)

### logic
fh = open("./results.md", "a+")
o = oh.OutputHelper(False)
o.header(fh)
fh.close()


@pytest.fixture()
def setup(request) -> tuple:
    fh = open("./results.md", "a+")
    print("\nSetting up resources...")
    # clone the repo
    repo_root = "./repo-test/"

    # create all the test tags locally
    test_tags = [
        'v0.1.0',
        'v0.2.0-dependabotgithub.0',
        'v0.2.0-dependabotgithub.1',
        'v0.2.0-dependabotgithub.2',
        'v0.2.0-migrate.0',
        'v0.2.0-migrate.1',
        'v1.0.0',
        'v1.0.0-migrate.0',
        'v1.0.0-migrate.1',
        'v1.0.0-migrate.2',
        'v1.0.0-migrate.3',
        'v1.0.0-migrate.4',
        'v1.0.0-migrate.5',
        'v1.0.0-migrate.6',
        'v1.0.0-migrate.7',
        'v1.1.0',
        'v1.1.0-usev2workflows.0',
        'v1.2.0',
        'v1.2.0-cleanuptfversion.0',
        'v1.2.0-cleanuptfversion.1',
        'v1.3.0',
        'v1.3.0-workspacemanager.0',
        'v1.3.0-workspacemanager.1',
        'v1.3.0-workspacemanager.10',
        'v1.3.0-workspacemanager.11',
        'v1.3.0-workspacemanager.12',
        'v1.3.0-workspacemanager.13',
        'v1.3.0-workspacemanager.14',
        'v1.3.0-workspacemanager.2',
        'v1.3.0-workspacemanager.3',
        'v1.3.0-workspacemanager.4',
        'v1.3.0-workspacemanager.5',
        'v1.3.0-workspacemanager.6',
        'v1.3.0-workspacemanager.7',
        'v1.3.0-workspacemanager.8',
        'v1.3.0-workspacemanager.9',
        'v1.4.0',
        'v1.4.0-updatetfversionl.0',
        'v1.5.0-moreactions.0'
    ]

    yield repo_root, test_tags, fh

    print("\nPerforming teardown...")
    try:
        fh.close()
    except OSError as e:
        print("Error: %s - %s." % (e.filename, e.strerror))


def test_moreactions_found(setup) -> None:
    """
    Test moreactions is found
    """
    repo_root, test_tags, fh = setup
    branch = "more-actions"
    release_branches = "main,master"
    expected = "v1.5.0-moreactions.0"
    outputs = cmd.run(test_tags, branch, release_branches, "true", "moreactions")
    t1 = (outputs['latest'] == expected)
    assert True == t1
    o.result(expected, "==", outputs['latest'], t1 == True, fh)

def test_123rand123_not_found_latest_empty(setup) -> None:
    """
    Test 123rand123 is not found and latest is empty
    """
    repo_root, test_tags, fh = setup
    branch = "more-actions"
    release_branches = "main,master"
    expected = ""
    outputs = cmd.run(test_tags, branch, release_branches, "true", "123rand123")
    t1 = (outputs['latest'] == expected)
    assert True == t1
    o.result(expected, "==", outputs['latest'], t1 == True, fh)

def test_last_release_version_with_empties(setup) -> None:
    """
    Test last release version is 1.4.0
    """
    repo_root, test_tags, fh = setup
    branch = "beta"
    release_branches = "main,master"
    expected = "v1.4.0"
    outputs = cmd.run(test_tags, branch, release_branches, "true", "")
    t1 = (outputs['last_release'] == expected)
    assert True == t1
    o.result(expected, "==", outputs['last_release'], t1 == True, fh)

def test_branch_matches_release_branch_forces_a_release(setup) -> None:
    """
    Test branch that matches a release branch forces a release version
    """
    repo_root, test_tags, fh = setup
    branch = "main"
    release_branches = "main,master"
    expected = "v1.4.0"
    outputs = cmd.run(test_tags, branch, release_branches, "true", "123rand123")
    t1 = (outputs['last_release'] == expected)
    assert True == t1
    o.result(expected, "==", outputs['last_release'], t1 == True, fh)
