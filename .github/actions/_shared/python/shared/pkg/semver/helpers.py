
import string

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
