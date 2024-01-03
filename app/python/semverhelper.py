#!/usr/bin/env python3
from semver.version import Version

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
        parsed = self.parse()
        # if there is a parsed version, and it came with a prefix
        # then respect that prefix and return with it
        # otherwise return without a prefix
        if parsed is not None and self.has_prefix():
            return f"{self.prefix}{parsed}"
        elif parsed is not None:
            return f"{parsed}"
        return self._tag

    def has_prefix(self) -> bool:
        """Determine if the string (s) passed starts with a v prefix for semver parsing."""
        return self._tag.startswith(self.prefix)

    def without_prefix(self) -> str:
        """Trim a prefix from the start of tag string."""        
        return (self._tag[1:] if self.has_prefix() else self._tag)
    
    def valid(self) -> bool:
        """Determine if tag is valid semver. Handles trimming of prefix"""
        return Version.is_valid(self.without_prefix())
    
    def parse(self, default:str = None) -> Version|None:
        """If the tag passed is a valid semver tag then return a version, otherwise return None"""
        if self.valid():
            self._parsed = Version.parse(self.without_prefix())
            return self._parsed
        elif default is not None:
            self._parsed = Version.parse(default)
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
        new_tag = f"v{tag}" if with_prefix else f"{tag}"
        self._tag = new_tag
        self._parsed = self.parse()

    ## STATICS
    @staticmethod
    def default(tag:str = "0.0.0") -> Version:
        return Version.parse(tag)
    
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
        as_dict = {}
        for tag in tags:
            d = SemverHelper.to_dict(tag)
            if d and d[tag] is not None:
                as_dict.update(d)
        return as_dict