#!/usr/bin/env python3
import pytest
from actions.common import strhelper as sth

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


@pytest.mark.parametrize(
    "test_data,expected",
    [
        ("acleanstring", "acleanstring"),
        ("a string with spaces", "astringwithspaces"),
        ("other-characters1-*/?-$", "othercharacters1"),
        ("MixedStr/6/':#", "mixedstr6")
    ]
)
def test_safe(test_data:str, expected:str) -> None:
    """
    Test the safe conversion
    """
    actual:bool = sth.safe(test_data)

    assert (expected == actual) == True
