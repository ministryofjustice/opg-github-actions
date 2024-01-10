#!/usr/bin/env python3
import pytest
from pyactions.common import strhelper as sth

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
