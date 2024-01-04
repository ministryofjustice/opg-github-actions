#!/usr/bin/env python3
from semver.version import Version
import os
import importlib.util
from git import Repo, Git
import pytest
import shutil

# CUSTOM PATH LOADING
dir_name = os.path.dirname(os.path.realpath(__file__))
# load semver helpers
str_mod = importlib.util.spec_from_file_location("strhelper", dir_name + '/strhelper.py')
sth = importlib.util.module_from_spec(str_mod)
str_mod.loader.exec_module(sth)


@pytest.mark.parametrize(
    "expected,source,default",
    [
        (True, "true", False),
        (False, "w", False),
        (True, "yes", False),
        (True, True, False),
        (False, False, False)
    ]
)

def test_str_to_bool(expected:bool, source:bool|str, default:bool) -> None:
    """
    Test the str_to_bool conversion
    """
    actual:bool = sth.str_to_bool(source, default)

    assert (expected == actual) == True
