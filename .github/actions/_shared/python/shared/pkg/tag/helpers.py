import os
import importlib.util
from natsort import natsorted
from semver.version import Version
from git import Repo, Git

parent_dir_name = os.path.dirname(os.path.dirname(os.path.realpath(__file__)))
# load semver helpers
semver_mod = importlib.util.spec_from_file_location("semverh", parent_dir_name + '/semver/helpers.py')
svh = importlib.util.module_from_spec(semver_mod)  
semver_mod.loader.exec_module(svh)
# load rand
rand_mod = importlib.util.spec_from_file_location("randh", parent_dir_name + '/rand/helpers.py')
rnd = importlib.util.module_from_spec(rand_mod)  
rand_mod.loader.exec_module(rnd)
# load github
gh_mod = importlib.util.spec_from_file_location("githubh", parent_dir_name + '/github/helpers.py')
gh = importlib.util.module_from_spec(gh_mod)  
gh_mod.loader.exec_module(gh)


def generate_tag_to_create(tag_name: str, all_tags: list) -> str:
    """
    Return a tag name that can be created in the repository. 
    For semver formatted tag_names we use a truncated version
    of the branch name, the tag may already exist due to matching substrings 
    (update-somethings-please would try and create the same as update-somethings)
    so when that happens create a new branch name.
    When its a none-semver branch that exists, attach a new random segment on the
    end
    """

    rand_length = 3
    original_tag = tag_name
    with_v = svh.has_v(original_tag)
    valid_semver = Version.parse(svh.trim_v(original_tag))
    # if this is semver, then parse and update it
    if valid_semver:
        parsed_tag = Version.parse(svh.trim_v(tag_name) if with_v else tag_name)
        while tag_name in all_tags:
            # if this is a pre-release, then we can adjust that
            if parsed_tag.prerelease is not None:
                parsed_tag = parsed_tag.replace(prerelease=f"{rnd.rand(rand_length)}.0")
                tag_name = f"v{parsed_tag}" if with_v else f"{parsed_tag}"
            # otherwise, bump version as this should be release
            else:
                parsed_tag = parsed_tag.bump_major()
                tag_name = f"v{parsed_tag}" if with_v else f"{parsed_tag}"
    # if its not a semver then tag on a random suffix
    else:
        while tag_name in all_tags:
            tag_name = f"{original_tag}.{rnd.rand(rand_length)}"
    return tag_name


def repo_tags(path:str, params:str) -> list:
    """
    Create the repo object and call out to the 
    github helper to run the command
    """
    repo = Repo(path)
    return gh.tags(repo, params)


def create_tag(path:str, commitish:str, tag_to_create:str, push:bool):
    """
    Create this tag locally and if push is true, the push to remote
    """
    repo = Repo(path)
    repo.git.tag(tag_to_create, commitish)
    if push:
        print(f"Pushing tag [{tag_to_create}] to remote")
        repo.git.push('origin', tag_to_create)