#!/usr/bin/env python3
import sys
from git import Repo
from natsort import natsorted
from . import semverhelper as svh, randhelper as rnd

class GitHelper:
    """
    Class to provider some commonly used calls to the gitpython lib
    """
    repository_path = ''
    repository = None

    def __init__(self, repository_path:str):
        self.repository_path = repository_path
        self.repository = Repo(self.repository_path)

    def tags(self, argument:str) -> list:
        """
        Use the repo object to fetch the tag data, split that
        by new lines and return a list
        - `param` to allow options (such as --list or --points-at)
        """
        all = list( self.repository.git.tag(argument).split("\n") )
        return natsorted(all)

    def create_tag(self, tag:str, commitish:str, push:bool = False) -> str:
        """
        Create tag locally at the `commistish` git reference and when
        `push` is True, send to remote.
        Presumes the tag has already been checked and corrected for use.
        """
        try:
            print(f"Creating tag [{tag}] locally")
            self.repository.git.tag(tag, commitish)
        except Exception as err:
            print(f"Fatal error: could not create [{tag}] locally", file=sys.stderr)
            raise Exception(f"Fatal error: could not create [{tag}] locally")

        if push == True:
            print(f"Pushing tag [{tag}] to remote")
            self.repository.git.push('origin', tag)
        return tag

    def tag_to_create(self, tag:str, all_tags:list) -> str:
        """
        Return a tag name that can be created in the repository.
        For semver formatted tag_names we use a truncated version
        of the branch name, the tag may already exist due to matching substrings
        (update-somethings-please would try and create the same as update-somethings)
        so when that happens create a new branch name.
        When its a non-semver branch that exists, attach a new random segment on the
        end
        """
        sv = svh.SemverHelper(tag)
        prefix = sv.has_prefix()
        random_length = 3
        original_tag = tag
        # If this is semver tag, then update it
        if sv.valid():
            prerelease = (sv.parsed().prerelease != None)
            while tag in all_tags:
                # if this is a prerelease, then swap the prerelease segment to be random
                if prerelease:
                    sv.update( sv.parsed().replace(prerelease=f"{rnd.rand(random_length)}.0"), prefix )
                # otherwise, if this is a release, then bump the major version
                else:
                    sv.update( sv.parsed().bump_major(), prefix )
                # refresh tag for loop
                tag = sv.tag()
        else:
            while tag in all_tags:
                tag = f"{original_tag}.{rnd.rand(random_length)}"

        return tag


    def commit_data(self, commitish_a, commitish_b) -> list:
        """
        Use the iteration method to find commits between the points
        and return list of them with a structured dict
        """
        parsed_commits:list = []
        for commit in self.repository.iter_commits(f"{commitish_a}...{commitish_b}"):
            parsed_commits.append({
                'hash': f"{commit}",
                'subject': commit.summary,
                'body': commit.message,
                'notes': ''
            })
        return parsed_commits


    def commits(self, commitish_a:str, commitish_b:str) -> list:
        """
        Fetch all commits between commitish_a and commitish_b
        in an xml like format for easier parsing and handling
        of sepcial chars in commit messages (like quotes and
        slashes)
        Convert log into a list
        """
        logs:list = []
        # checkout between the locations to ensure we have logs
        try:
            self.repository.git.checkout(commitish_a, "--")
            self.repository.git.checkout(commitish_b, "--")
        except Exception:
            print(f"Failed to checkout [{commitish_a} or {commitish_b}]")
            raise Exception(f"Failed to checkout to a commit [{commitish_a} or {commitish_b}]")

        log:list = self.commit_data(commitish_a, commitish_b)
        return log


## NON CLASS FUNCTIONS
def find_bumps_from_commits(commits:list, default_bump:str) -> tuple:
    """
    Scan all fields in the commits passed looking for triggers of each type.
    Return counter of each.
    """
    majors=0
    minors=0
    patches=0
    for c in commits:
        # check each field in the dict
        for k in ['subject', 'notes', 'body']:
            majors = majors + 1 if "#major" in c[k] else majors
            minors = minors + 1 if "#minor" in c[k] else minors
            patches = patches + 1 if "#patch" in c[k] else patches
    # if nothing has been found, use the default
    if majors == 0 and minors == 0 and patches == 0:
        if default_bump == "major":
            majors = 1
        elif default_bump == "minor":
            minors = 1
        elif default_bump == "patch":
            patches = 1

    return majors, minors, patches


def github_branch_data(event_name:str, event_data:dict) -> tuple:
    """
    Use the event name and data from github env to workout
    branch name to use and commitish references for comparison
    later on
    """
    source_commitish = ''
    destination_commitish = ''
    branch_name = ''
    # for pull requests we do this
    if event_name == "pull_request":
        # this is branch that started the pull request (test-123)
        source_commitish = event_data['pull_request']['head']['ref']
        # this should be something like main / master
        destination_commitish = event_data['pull_request']['base']['ref']
        # active branch is the same as source_branch on a pr
        branch_name = source_commitish
    elif event_name == "push":
        branch_name = event_data['ref']
        # use the before and after properties
        source_commitish = event_data['before']
        destination_commitish = event_data['after']

    branch_name = branch_name.replace('refs/head/', '').replace('refs/heads/', '')
    return (branch_name, source_commitish, destination_commitish)
