#!/usr/bin/env python3
import os
import sys
from git import Repo
from natsort import natsorted
import importlib.util
import xmltodict
import re

# CUSTOM PATH LOADING
dir_name = os.path.dirname(os.path.realpath(__file__))
# load semver helpers
semver_mod = importlib.util.spec_from_file_location("semverhelper", dir_name + '/semverhelper.py')
svh = importlib.util.module_from_spec(semver_mod)
semver_mod.loader.exec_module(svh)
# load rand helpers
rand_mod = importlib.util.spec_from_file_location("randhelper", dir_name + '/randhelper.py')
rnd = importlib.util.module_from_spec(rand_mod)
rand_mod.loader.exec_module(rnd)



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
            self.repository.git.tag(tag, commitish)
        except Exception as err:
            print(f"Fatal error: could not create [{tag}] locally", file=sys.stderr)
            raise Exception(f"Fatal error: could not create [{tag}] locally")

        if push == True:
            print(f"Pushing tag [{tag}] to remote")
            self.repository.push('origin', tag)
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

    def commits(self, commitish_a:str, commitish_b:str) -> list:
        """
        Fetch all commits between commitish_a and commitish_b
        in an xml like format for easier parsing and handling
        of sepcial chars in commit messages (like quotes and
        slashes)
        Convert log into a list
        """
        logs = []
        # checkout between the locations to ensure we have logs
        try:
            self.repository.git.checkout(commitish_a)
            self.repository.git.checkout(commitish_b)
        except Exception:
            print("Failed to checkout")
            raise Exception("Failed to checkout to a commit")

        # get data from the log in an almost json format
        xmlish = "<commit><hash>%H</hash><subject>%s</subject><notes>%s</notes><body>%s</body></commit>"
        param = f"--pretty=format:{xmlish}"
        range = f"{commitish_a}...{commitish_b}"
        log_data = self.repository.git.log(param, range)
        # wrap log data in container tag for parsing
        log_data = f"<commits>\n{log_data.strip()}</commits>"
        logs = xmltodict.parse(log_data)
        # grab the list
        logs = logs['commits']['commit']
        return logs

    @staticmethod
    def find_bumps_from_commits(commits:list, default_bump:str) -> tuple:
        """
        Scan all fields in the commits passed looking for triggers of each type.
        Return counter of each.
        The count for default_bump starts at 1 instead of 0 to ensure something is
        always increased.
        """
        majors=1 if default_bump == "major" else 0
        minors=1 if default_bump == "minor" else 0
        patches=1 if default_bump == "patch" else 0
        for c in commits:
            # check each field in the dict
            for k in ['subject', 'notes', 'body']:
                majors = majors + 1 if "#major" in c[k] else majors
                minors = minors + 1 if "#minor" in c[k] else minors
                patches = patches + 1 if "#patch" in c[k] else patches

        return majors, minors, patches
