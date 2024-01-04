#!/usr/bin/env python3
from semver.version import Version
import re

class SemverHelper:
    """
    Handles parsing and processing of semver formatted strings
    Typically dealing with tags from git
    Works with prefixes on the semver version string
    """
    prefix:str = 'v'
    original:str = None
    _tag:str = None
    _parsed:Version = None

    def __init__(self, tag:str):
        self.original = tag
        self._tag = tag
        self._parsed = self.parse()

    def __str__(self):
        return self.tag()

    def tag(self) -> str:
        """Return the string version of the parsed version tag"""
        parsed:Version|None = self.parse()
        # if there is a parsed version, and it came with a prefix
        # then respect that prefix and return with it
        # otherwise return without a prefix
        if parsed is not None and self.has_prefix():
            return f"{self.prefix}{parsed}"
        elif parsed is not None:
            return f"{parsed}"
        return self._tag

    def has_prefix(self) -> bool:
        """
        Determine if the string (s) passed starts with a v prefix for semver parsing.
        """
        return self._tag.startswith(self.prefix)

    def without_prefix(self) -> str|None:
        """Trim a prefix from the start of tag string."""
        if self._tag is None:
            return None
        else:
            return (self._tag[1:] if self.has_prefix() else self._tag)

    def valid(self) -> bool:
        """Determine if tag is valid semver. Handles trimming of prefix"""
        if self._tag is None:
            return False
        else:
            return Version.is_valid(self.without_prefix())

    def parse(self, default:str = None) -> Version|None:
        """If the tag passed is a valid semver tag then return a version, otherwise return None"""
        if self.valid():
            self._parsed = Version.parse(self.without_prefix())
            return self._parsed
        elif default is not None:
            self._parsed = SemverHelper.default(default)
            return self._parsed
        return None

    def parsed(self) -> Version|None:
        """Return the active version thats been processed"""
        return self._parsed

    def update(self, tag:str|Version, with_prefix:bool = False) -> None:
        """
        Use the new tag thats passed in to update the current values.
        This is to allow increasing parts of the tag
        """
        new_tag:str = f"v{tag}" if with_prefix else f"{tag}"
        self._tag = new_tag
        self._parsed = self.parse()

    ## STATICS
    @staticmethod
    def default(tag:str = "0.0.0") -> Version:
        """Generate a default Version instance, using tag (should not contain a prefix)."""
        return Version.parse(tag)

    @staticmethod
    def is_prerelease(branch_name:str, release_branches:str|list, stated_prerelease_state:bool) -> bool:
        """
        Look at possibly overwriting if this is a prerelease or not - as in branch is main, but
        accidently set prerelease as false.

        Check if the branch_name is within the set of release_branches, if so return true.
        Otherwise, return the currently set value - likely from original command input
        """
        # convert to list from string
        release_branches:list = release_branches.split(',') if type(release_branches) is str else release_branches
        # if the branch is in the release branch set, return True, otherwise return current
        if branch_name in release_branches:
            return False
        return stated_prerelease_state


    @staticmethod
    def prereleases_filtered(tag_list:list, filter:str) -> dict:
        """
        Find all semver prereleases that match the filter pattern passed in.
        Designed to find specific ones relating to a branch (ie mybranch.0, mybranch.1 etc).
        """
        tags:dict = SemverHelper.list_to_dict(tag_list)
        pattern = re.compile(filter)
        prereleases:dict = {k:v for k,v in tags.items() if pattern.match(f"{v.prerelease}")}
        return prereleases

    @staticmethod
    def prereleases(tag_list:list) -> dict:
        """
        Find all semver prereleases from a set of strings (likely tags from git) and
        return a dict whose key is the tag and value is a Version instance
        """
        tags:dict = SemverHelper.list_to_dict(tag_list)
        prereleases:dict = {k:v for k,v in tags.items() if v.prerelease is not None}
        return prereleases

    @staticmethod
    def releases(tag_list:list) -> dict:
        """
        Find all semver releases from a set of strings (likely tags from git) and
        return a dict whose key is the tag and value is a Version instance
        """
        tags:dict = SemverHelper.list_to_dict(tag_list)
        releases:dict = {k:v for k,v in tags.items() if v.prerelease is None}
        return releases

    @staticmethod
    def max(versions:dict|list) -> Version|None:
        """
        Take a dict of Version instances (tag->Version) and return the max (natural ordering)
        version

        If there is a release and pre-release that match, release is used
        >>> SemverHelper.max(['v1.0.0-beta.9', 'v1.0.0-beta.10', '100.0.0', '100.0.0-test.0'])
        Version('100.0.0')
        """
        if type(versions) is list:
            versions:dict = SemverHelper.list_to_dict(versions)

        if len(versions) > 0:
            max_value = max(versions.values())
            max_items = [k for k,v in versions.items() if f"{v}" == max_value]
            if len(max_items) > 0:
                return max_items.pop()
        return None



    @staticmethod
    def to_dict(tag:str) -> dict:
        """
        Returns a dict containing a semver parsed version of the string tag passed.
        If the tag is not a valid semver, then the dict.tag will be None
        Uses the tag string passed as the key in the dict
        """
        s = SemverHelper(tag)
        return {tag: s.parse() }

    @staticmethod
    def list_to_dict(tags:list) -> dict:
        """
        Take a list of string tags, convert (when valid) each to a dict
        with key being its original tag and value being the semver Version
        of it and append
        - tag = '1.2.3-test.0' => {'1.2.3-test.0' => Version() }
        """
        as_dict:dict = {}
        for tag in tags:
            d = SemverHelper.to_dict(tag)
            if d and d[tag] is not None:
                as_dict.update(d)
        return as_dict
