#!/usr/bin/env python3
import os
import importlib.util
import pytest

### local imports
app_root_dir = os.path.dirname(
    os.path.dirname(
        os.path.dirname(
            os.path.dirname(os.path.realpath(__file__))
        )
    )
)
dir_name = os.path.dirname(os.path.realpath(__file__))

# load cmd helper
mod = importlib.util.spec_from_file_location("latest-tag", dir_name + '/branch-name.py')
cmd = importlib.util.module_from_spec(mod)
mod.loader.exec_module(cmd)
# load output helper
ohmod = importlib.util.spec_from_file_location("gh", app_root_dir + '/app/python/outputhelper.py')
oh = importlib.util.module_from_spec(ohmod)
ohmod.loader.exec_module(oh)

### RESULT FILE
fh = open("./results.md", "a+")
o = oh.OutputHelper(False)
o.header(fh)
fh.close()



### SETUP TEST CONFIG DATA
testconfig = [
    {
        "expected": "test-branch",
        "event_name": "pull_request",
        "event_data": {
            'pull_request': {
                'head': {'ref': "test-branch"},
                'base': {'ref': "main"}
            }
        }
    },
    {
        "expected": "main",
        "event_name": "push",
        "event_data": {
            'ref': 'refs/head/main',
            'before': 'a12313',
            'after': 'dfg1231243'
        }
    },

]
# generate fields string from the keys
fields = ','.join(testconfig[0].keys())
# generate test tuple from config items
tests = [(v.values()) for v in testconfig]

@pytest.mark.parametrize(fields, tests)
def test_branch_name(
    expected:str,
    event_name:str,
    event_data:dict) -> None:
    """
    Check the branch data return matches expected values
    """
    outputs = cmd.run(
        event_name=event_name,
        event_data=event_data,
    )
    # dump data for debugging
    print(f"Expected {expected} Actual {outputs['branch_name']}")
    print(outputs, sep="\n")

    t1 = (outputs['branch_name'] == expected)
    fh = open("./results.md", "a+")
    o = oh.OutputHelper(False)
    o.result(expected, "==", outputs['branch_name'], t1 == True, fh)
    fh.close()

    assert t1 == True
