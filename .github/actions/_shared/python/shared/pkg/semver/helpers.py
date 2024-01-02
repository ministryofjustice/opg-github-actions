
import string
from semver.version import Version

def has_v(s: str) -> bool:
    """
    Determine if the string (s) passed starts with a v prefix for semver parsing.
    """
    return s.startswith('v')
# removes v from start of string - crude fix for semver tags having a v prefix
def trim_v(s: str):
    """
    Trim a v prefix from the start of a semver formatted string.
    Currently very crude and just checks the first char, removing it
    """
    return (s[1:] if has_v(s) else s)    


def to_valid_dict(tag) -> dict|None:
    """Converts a tag that may contain a v prefix into a valid semver version."""
    t = trim_v(tag)
    if Version.is_valid(t):            
        return {tag: Version.parse(t)}
    return None


def semver_list(tags:list) -> dict:
    """Takes a list of tags (as strings) and fetches just the semver versions of them as a dict."""
    semver_tags = {}
    for tag in tags:
        d = to_valid_dict(f"{tag}")
        if d is not None:
            semver_tags.update(d)
    return semver_tags
