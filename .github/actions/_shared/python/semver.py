
import string

def has_v(s: str) -> bool:
    return s.startswith('v')
# removes v from start of string - crude fix for semver tags having a v prefix
def trim_v(s: str):
    return (s[1:] if has_v(s) else s)    
