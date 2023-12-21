from git import Repo, Git
import os
import re
import semver
import argparse
from semver.version import Version


def arg_parser() -> argparse.ArgumentParser:
    parser = argparse.ArgumentParser("parse-tags")
    parser.add_argument('--test_file', default="", help="trigger the use of a test file for list of tags")
    
    parser.add_argument('--repository_root', default="./", help="Path to root of repository")    
    parser.add_argument('--default_branch', default="main", help="Base branch to compare against. (Default: main)")    
    
    parser.add_argument('--prerelease', default="", help="If set, then this is a pre-release")    
    parser.add_argument("--prerelease_suffix", default="beta", help="Prerelease naming")
    parser.add_argument("--latest_tag", default="", help="Last tag")
    parser.add_argument("--last_release", default="", help="Last release var tag")
    parser.add_argument("--with_v", default="false", help="apply prefix to the new tag")
    return parser

def trim_v(str):
    return (str[1:] if str.startswith('v') else str)    

def split_commits_from_lines(lines):
    split_lines = [[str(i) for i in line.strip().split(" ", 1)] for line in lines]
    return split_lines

def compare_to(latest_tag):
    compare=latest_tag
    if len(latest_tag) == 0:
        compare = "HEAD"
    return compare

def starting_tag(latest_tag, initial_version, prerelease, prerelease_suffix):
    if len(latest_tag) == 0:        
        tag = Version.parse(trim_v(initial_version))
        if prerelease:
            tag = tag.replace(prerelease=f"{prerelease_suffix}.0")
    else:
        tag = Version.parse(trim_v(latest_tag))

    return tag

def get_commits(test, test_file, default_branch, latest_tag):
     #use test data
    if test is not None and len(test) > 0 and len(test_file) > 0 :
        print("Using test data")
        commits = split_commits_from_lines( open(test_file).readlines() )
    else:
        commits = g.log("--oneline", f"{default_branch}...{latest_tag}")
        print(f"Getting commits between [{default_branch}]...[{latest_tag}]")
        commits = split_commits_from_lines( commits.split("\n") )
    return commits

def main():
    # get the args
    args = arg_parser().parse_args()
    test_file = args.test_file
    
    repo_root = args.repository_root
    default_branch = args.default_branch

    prerelease = (len(args.prerelease) > 0)    
    prerelease_suffix = args.prerelease_suffix
    last_release = args.last_release
    latest_tag = args.latest_tag
    # set the intial version to be last_release or 0.0.1 if thats empty
    initial_version = last_release if len(last_release) > 0 else "0.0.0"
    
    # test info setup
    test = os.getenv("RUN_AS_TEST")
    is_test = (test is not None and len(test) > 0 and len(test_file) > 0)

    tag = starting_tag(latest_tag, initial_version, prerelease, prerelease_suffix)
    compare = compare_to(latest_tag)
    g = Git(repo_root)
    
    # get the commits between shas
    commits = get_commits(test, test_file, default_branch, latest_tag)
    # look for #major, #minor #patch in commits
    major=0
    minor=0
    patch=0
    for c in commits:
        major = major + 1 if "#major" in c[1] else major
        minor = minor + 1 if "#minor" in c[1] else minor
        patch = patch + 1 if "#patch" in c[1] else patch

    print(f"Majors: [{major}] Minors: [{minor}] Patches: [{patch}]")

    # get the last release
    last_release = Version.parse(trim_v(last_release)) if Version.is_valid(trim_v(last_release)) else Version.parse(initial_version)
    new_tag = tag

    if major > 0 and new_tag.major <= last_release.major:
        print ("-> major bump")
        # if its prerelease, then re-add the suffix as a 0
        new_tag = new_tag.bump_major().replace(prerelease=f"{prerelease_suffix}.0") if prerelease else new_tag.bump_major()
    elif major == 0 and minor > 0 and new_tag.minor <= last_release.minor:
        print ("-> minor bump")
        new_tag = new_tag.bump_minor().replace(prerelease=f"{prerelease_suffix}.0") if prerelease else new_tag.bump_minor()
    elif major == 0 and minor ==0 and patch > 0 and new_tag.patch <= last_release.patch:
        print ("-> minor bump")
        new_tag = new_tag.bump_patch().replace(prerelease=f"{prerelease_suffix}.0") if prerelease else new_tag.bump_patch()
    elif prerelease:
        print ("-> prerelease")
        new_tag = new_tag.bump_prerelease()

    if not prerelease:
        new_tag = new_tag.replace(prerelease=None)
    
    new_tag_str = f"{new_tag}"
    
    if len(args.with_v) > 0 and args.with_v == "true":
        new_tag_str = f"v{new_tag_str}"

    print(f"prerelease={args.prerelease}")
    print(f"prerelease_processed={prerelease}")
    print(f"last_release={args.last_release}")
    print(f"last_release_processed={last_release}")
    print(f"initial_version={initial_version}")
    print(f"lastest_tag={args.latest_tag}")
    print(f"starting_tag={tag}")
    print(f"next_tag={new_tag_str}")    

    if 'GITHUB_OUTPUT' in os.environ:
        print("Pushing to GitHub Output")
        with open(os.environ['GITHUB_OUTPUT'], 'a') as fh:
            print(f"prerelease={args.prerelease}", file=fh)
            print(f"prerelease_processed={prerelease}", file=fh)
            print(f"last_release={args.last_release}", file=fh)
            print(f"last_release_processed={last_release}", file=fh)
            print(f"initial_version={initial_version}", file=fh)
            print(f"lastest_tag={args.latest_tag}", file=fh)
            print(f"starting_tag={tag}", file=fh)
            print(f"next_tag={new_tag_str}", file=fh)  
            
    


if __name__ == "__main__":
    main()