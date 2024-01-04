#!/usr/bin/env python3
import string

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
