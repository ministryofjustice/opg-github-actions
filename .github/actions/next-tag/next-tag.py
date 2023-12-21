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
    initial_version = last_release if len(last_release) > 0 else "0.0.1"
    
    # test info setup
    test = os.getenv("RUN_AS_TEST")
    is_test = False

    tag=starting_tag(latest_tag, initial_version, prerelease, prerelease_suffix)
    compare=compare_to(latest_tag)
    g = Git(repo_root)
    
    # get the commits between shas
        #use test data
    if test is not None and len(test) > 0 and len(test_file) > 0 :
        print("Using test data")
        is_test = True
        commits = split_commits_from_lines( open(test_file).readlines() )
    else:
        commits = g.log("--oneline", f"{default_branch}...{latest_tag}")
        commits = split_commits_from_lines( commits.split("\n") )
    
    # look for #major, #minor #patch in commits
    major=0
    minor=0
    patch=0
    for c in commits:
        major = major + 1 if "#major" in c[1] else major
        minor = minor + 1 if "#minor" in c[1] else minor
        patch = patch + 1 if "#patch" in c[1] else patch

    print(*commits, sep="\n")
    print(f"Majors: [{major}] Minors: [{minor}] Patches: [{patch}]")

    # get the last release
    last_release = Version.parse(trim_v(last_release)) if Version.is_valid(trim_v(last_release)) else Version.parse(initial_version)

    print(f"last_release:{last_release}")
    print(f"tag:{tag}")
    new_tag = tag
    #v1.0 => v2.0
    #v1.4.0-beta.0 when commit contains #major => v2.0.0-beta.0     
    if major > 0 and tag.major <= last_release.major:
        bumped = last_release.bump_major()
        new_tag = Version.parse(f"{bumped.major}.0.0")
        if prerelease:
            new_tag = new_tag.replace(prerelease=f"{prerelease_suffix}.0")   
    #v1.5.0-beta.1 => v1.5.0 when released
    elif prerelease is False and minor > 0:
        new_tag = last_release.bump_minor()
    #v1.5.0
    elif prerelease is False and patch > 0:
        new_tag = last_release.bump_patch()
    

    # if this is a pre-release, not a mjor and has an existing tag then increase the counter
    if prerelease and len(latest_tag) > 0:
        new_tag = new_tag.bump_prerelease()    
    
    print(f"new_tag:{new_tag}")    
        
    


if __name__ == "__main__":
    main()