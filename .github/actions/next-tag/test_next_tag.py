#!/usr/bin/env python3
from semver.version import Version
import os
import importlib.util
from git import Repo, Git
import pytest
import shutil

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
mod = importlib.util.spec_from_file_location("latest-tag", dir_name + '/next-tag.py')
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


majors = [
    {"commit": "b5df8f1", "subject": "test1", "notes": "", "body": "nothing"},
    {"commit": "b5df8f2", "subject": "test2", "notes": "", "body": "#major please"},
    {"commit": "b5df8f3", "subject": "test 3 #major", "notes": "", "body": "nothing here"},
    {"commit": "b5df8f4", "subject": "test", "notes": "", "body": "nothing here"},
    {"commit": "b5df8f5", "subject": "test", "notes": "", "body": "#minor fix"},
]
minors = [
    {"commit": "b5df8f1", "subject": "test1", "notes": "", "body": "nothing"},
    {"commit": "b5df8f2", "subject": "test2", "notes": "", "body": "#minor please"},
    {"commit": "b5df8f3", "subject": "test", "notes": "", "body": "nothing here"},
    {"commit": "b5df8f4", "subject": "test", "notes": "", "body": "nothing here"},
    {"commit": "b5df8f5", "subject": "test", "notes": "", "body": "fix"},
]
patches = [
    {"commit": "b5df8f1", "subject": "test1", "notes": "", "body": "nothing"},
    {"commit": "b5df8f2", "subject": "test2", "notes": "", "body": "#what please"},
    {"commit": "b5df8f3", "subject": "test", "notes": "", "body": "nothing here"},
    {"commit": "b5df8f4", "subject": "test", "notes": "", "body": "nothing here"},
    {"commit": "b5df8f5", "subject": "test", "notes": "", "body": "#patch"},
]
none = [
    {"commit": "b5df8f1", "subject": "test1", "notes": "", "body": "nothing"},
    {"commit": "b5df8f2", "subject": "test2", "notes": "", "body": "#what please"},
    {"commit": "b5df8f3", "subject": "test", "notes": "", "body": "nothing here"},
    {"commit": "b5df8f4", "subject": "test", "notes": "", "body": "nothing here"},
    {"commit": "b5df8f5", "subject": "test", "notes": "", "body": "patch me"},
]

### SETUP ALL THE EXPECTED MATCHES FOR STANDARD TAG TESTS
testconfig = [
    {
        "expected": "2.0.0-moreactions.0",
        "prerelease": "true",
        "prerelease_suffix": "moreactions",
        "latest_tag": "v1.5.0-moreactions.1",
        "last_release": "v1.4.0",
        "default_bump": "patch",
        "with_v": "",
        "commits": majors,
    },
    {
        "expected": "2.0.0-moreactions.1",
        "prerelease": "true",
        "prerelease_suffix": "moreactions",
        "latest_tag": "v2.0.0-moreactions.0",
        "last_release": "v1.4.0",
        "default_bump": "patch",
        "with_v": "",
        "commits": majors,
    },
    {
        "expected": "2.0.0-moreactions.0",
        "prerelease": "true",
        "prerelease_suffix": "moreactions",
        "latest_tag": None,
        "last_release": "v1.4.0",
        "default_bump": "patch",
        "with_v": "",
        "commits": majors,
    },
    {
        "expected": "2.0.0",
        "prerelease": "",
        "prerelease_suffix": "moreactions",
        "latest_tag": None,
        "last_release": "v1.4.0",
        "default_bump": "patch",
        "with_v": "",
        "commits": majors,
    },
    {
        "expected": "1.0.0",
        "prerelease": "",
        "prerelease_suffix": "moreactions",
        "latest_tag": None,
        "last_release": None,
        "default_bump": "patch",
        "with_v": "",
        "commits": majors,
    },
    {
        "expected": "1.5.0-moreactions.1",
        "prerelease": "true",
        "prerelease_suffix": "moreactions",
        "latest_tag": "v1.5.0-moreactions.0",
        "last_release": "v1.4.0",
        "default_bump": "patch",
        "with_v": "",
        "commits": minors,
    },
    {
        "expected": "v1.5.0-moreactions.0",
        "prerelease": "true",
        "prerelease_suffix": "moreactions",
        "latest_tag": None,
        "last_release": "v1.4.0",
        "default_bump": "patch",
        "with_v": "true",
        "commits": minors,
    },
    {
        "expected": "v1.5.0",
        "prerelease": "",
        "prerelease_suffix": "moreactions",
        "latest_tag": "v1.5.0-moreactions.0",
        "last_release": "v1.4.0",
        "default_bump": "patch",
        "with_v": "true",
        "commits": minors,
    },
    {
        "expected": "0.1.0",
        "prerelease": "",
        "prerelease_suffix": "moreactions",
        "latest_tag": None,
        "last_release": None,
        "default_bump": "patch",
        "with_v": "",
        "commits": minors,
    },
    {
        "expected": "1.4.1-moreactions.0",
        "prerelease": "true",
        "prerelease_suffix": "moreactions",
        "latest_tag": None,
        "last_release": "v1.4.0",
        "default_bump": "patch",
        "with_v": "",
        "commits": patches,
    },
    {
        "expected": "1.4.1",
        "prerelease": None,
        "prerelease_suffix": "moreactions",
        "latest_tag": "1.4.1-moreactions.1",
        "last_release": "v1.4.0",
        "default_bump": "patch",
        "with_v": "",
        "commits": patches,
    },
    {
        "expected": "0.0.1",
        "prerelease": None,
        "prerelease_suffix": "moreactions",
        "latest_tag": None,
        "last_release": None,
        "default_bump": "patch",
        "with_v": "",
        "commits": patches,
    },
    {
        "expected": "1.0.1-moreactions.1",
        "prerelease": True,
        "prerelease_suffix": "moreactions",
        "latest_tag": "v1.0.1-moreactions.0",
        "last_release": Version.parse("1.0.0"),
        "default_bump": "patch",
        "with_v": "",
        "commits": none,
    },
    {
        "expected": "1.0.1-moreactions.0",
        "prerelease": True,
        "prerelease_suffix": "moreactions",
        "latest_tag": None,
        "last_release": Version.parse("1.0.0"),
        "default_bump": "patch",
        "with_v": "",
        "commits": none,
    },
    {
        "expected": "2.0.0-moreactions.0",
        "prerelease": True,
        "prerelease_suffix": "moreactions",
        "latest_tag": None,
        "last_release": Version.parse("1.0.0"),
        "default_bump": "major",
        "with_v": "",
        "commits": none,
    },
    {
        "expected": "1.1.0",
        "prerelease": False,
        "prerelease_suffix": "moreactions",
        "latest_tag": None,
        "last_release": Version.parse("1.0.0"),
        "default_bump": "minor",
        "with_v": "",
        "commits": none,
    },
    {
        "expected": "0.0.1",
        "prerelease": False,
        "prerelease_suffix": "moreactions",
        "latest_tag": None,
        "last_release": None,
        "default_bump": "patch",
        "with_v": "",
        "commits": none,
    },
    {
        "expected": "0.0.1-moreactions.0",
        "prerelease": True,
        "prerelease_suffix": "moreactions",
        "latest_tag": None,
        "last_release": None,
        "default_bump": "patch",
        "with_v": "",
        "commits": none,
    },
    {
        "expected": "v2.0.0",
        "prerelease": False,
        "prerelease_suffix": "moreactions",
        "latest_tag": None,
        "last_release": "v1.1.0",
        "default_bump": "major",
        "with_v": True,
        "commits": none,
    },
]
# generate fields string from the keys
fields = ','.join(testconfig[0].keys())
# generate test tuple from config items
tests = [(v.values()) for v in testconfig]

@pytest.mark.parametrize(fields, tests)
def test_next_tag_result_matches(
    expected:str,
    prerelease:str,
    prerelease_suffix:str,
    latest_tag:Version|None|str,
    last_release:Version|None|str,
    default_bump:str,
    with_v:str,
    commits:list) -> None:
    """
    Check that the config passed results in the
    next_tag matching the expected value.

    Use a parameterised test to reduce the 
    setup burden
    """
    outputs = cmd.run(
        test = True, 
        last_release = last_release, 
        latest_tag = latest_tag, 
        prerelease = prerelease, 
        prerelease_suffix = prerelease_suffix, 
        default_bump = default_bump, 
        with_v = with_v,
        commits = commits
    )    
    t1 = (outputs['next_tag'] == expected)    
    # dump data for debugging
    print(f"Expected {expected} Actual {outputs['next_tag']}")
    print(outputs, sep="\n")
    fh = open("./results.md", "a+")
    o = oh.OutputHelper(False)
    o.result(expected, "==", outputs['next_tag'], t1 == True, fh)  
    fh.close()

    assert t1 == True

