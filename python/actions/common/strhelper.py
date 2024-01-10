#!/usr/bin/env python3

def str_to_bool(source:str|bool, default:bool=False) -> bool:
    """
    source might be a string or a bool (input from cli as well as func calls)
    so convert to a bool

    Provide a default option to allow for params needing to be true
    """
    # if its already a bool, return
    if type(source) is bool:
        return source
    elif type(source) is str and source.lower() in ["true", "yes", "y"]:
        return True
    return default

def safe(source:str) -> str:
    """
    Convert string to lowercase alphanumeric only versio
    of itself
    """
    source = source.lower()
    source = ''.join([ c if c.isalnum() else '' for c in source ])
    return source

def int_or_none(source:str|int|None) -> int|None:
    """
    """
    if type(source) is int:
        return source
    elif type(source) is str and len(source) > 0:
        return int(source)

    return None
