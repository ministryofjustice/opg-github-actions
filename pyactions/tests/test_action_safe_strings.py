#!/usr/bin/env python3
import pytest
from pyactions.common import outputhelper as oh
from pyactions.actions import safe_strings as cmd

### RESULT FILE
fh = open("./results.md", "a+")
o = oh.OutputHelper(False)
o.header(fh)
fh.close()



### SETUP TEST CONFIG DATA
testconfig = [
    {
        "expected": "botlikesslashesanddashes",
        "test_string": "bot/likes/slashes-and-dashes",
        "suffix": None,
        "length": None,
        "conditional_match": None,
        "conditional_value": None,
    },
    {
        "expected": "avery",
        "test_string": "a-very-very-very-very-long-string-that-should-not-be-used",
        "suffix": None,
        "length": "5",
        "conditional_match": None,
        "conditional_value": None,
    },
    {
        "expected": "branch1",
        "test_string": "branch",
        "suffix": "1",
        "length": None,
        "conditional_match": None,
        "conditional_value": None,
    },
    {
        "expected": "branch1",
        "test_string": "branchname-to-shortern",
        "suffix": "1",
        "length": 7,
        "conditional_match": None,
        "conditional_value": None,
    },
    {
        "expected": "production",
        "test_string": "main",
        "suffix": None,
        "length": None,
        "conditional_match": "main",
        "conditional_value": "production",
    },

]
# generate fields string from the keys
fields = ','.join(testconfig[0].keys())
# generate test tuple from config items
tests = [(v.values()) for v in testconfig]

@pytest.mark.parametrize(fields, tests)
def test_safe_strings(
    expected:str,
    test_string:str,
    suffix:str|None,
    length:str|int|None,
    conditional_match:str|None,
    conditional_value:str|None) -> None:
    """
    Check the safe string return matches expected values
    """
    outputs = cmd.run(
        original=test_string,
        suffix=suffix,
        length=length,
        conditional_match=conditional_match,
        conditional_value=conditional_value
    )
    t1 = (outputs['safe'] == expected)
    # dump data for debugging
    print(f"Expected {expected} Actual {outputs['safe']}")
    print(outputs, sep="\n")
    fh = open("./results.md", "a+")
    o = oh.OutputHelper(False)
    o.result(expected, "==", outputs['safe'], t1 == True, fh)
    fh.close()

    assert t1 == True
