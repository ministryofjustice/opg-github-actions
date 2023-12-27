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
    parser.add_argument("--default_bump", default="patch", help="If there are no triggers in commits, bump by this")
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

def initial_semver_tag(tag):   
    return Version.parse(trim_v(tag))

def get_commits(test, test_file, repo_root, default_branch, latest_tag):
     #use test data
    if test is not None and len(test) > 0 and len(test_file) > 0 :
        print("Using test data")
        commits = split_commits_from_lines( open(test_file).readlines() )
    else:
        g = Git(repo_root)
        commits = g.log("--oneline", f"{default_branch}...{latest_tag}")
        print(f"Getting commits between [{default_branch}]...[{latest_tag}]")
        commits = split_commits_from_lines( commits.split("\n") )
    return commits

def main():
    # get the args
    args = arg_parser().parse_args()
    test_file = args.test_file
    
    is_prerelease = (len(args.prerelease) > 0)    
    # set the intial version to be last_release or 0.0.1 if thats empty
    base = "0.0.0"
    last_release = Version.parse(trim_v(args.last_release)) if Version.is_valid(trim_v(args.last_release)) else Version.parse(base)
    latest_tag = Version.parse(trim_v(args.latest_tag)) if len(args.latest_tag) > 0 and Version.is_valid(trim_v(args.latest_tag)) else None

    # test info setup
    test = os.getenv("RUN_AS_TEST")
    is_test = (test is not None and len(test) > 0 and len(test_file) > 0)

    starting_tag = last_release
    compare = compare_to(args.latest_tag)
    # get the commits between shas
    commits = get_commits(test, test_file, args.repository_root, args.default_branch, args.latest_tag)
    # look for #major, #minor #patch in commits
    # - use the default_bump to always increase one
    major=1 if args.default_bump == "major" else 0
    minor=1 if args.default_bump == "minor" else 0
    patch=1 if args.default_bump == "patch" else 0
    for c in commits:
        major = major + 1 if "#major" in c[1] else major
        minor = minor + 1 if "#minor" in c[1] else minor
        patch = patch + 1 if "#patch" in c[1] else patch

   
    print(f"Majors: [{major}] Minors: [{minor}] Patches: [{patch}]")
    
    new_tag = starting_tag

    # Bump the tag along based on what was found
    if major > 0:
        print ("-> major bump")
        new_tag = new_tag.bump_major()
    elif minor > 0:
        print ("-> minor bump")
        new_tag = new_tag.bump_minor()
    elif patch > 0:
        print ("-> patch bump")
        new_tag = new_tag.bump_patch()

    # If this is a prerelease, and there is a pre-existing tag we should copy over
    # the prerelease information to the new tag
    # existing_tag = v2.0.0-moreactions.0 would become v2.0.0-moreactions.1
    if is_prerelease and latest_tag is not None:
        print ("-> prerelease bump with tag")
        new_tag = new_tag.replace(prerelease=latest_tag.prerelease).bump_prerelease()
    # If this prerelease is the first of its kind the setup the prerelease segment
    # to use the suffix
    elif is_prerelease and latest_tag is None:
         print ("-> prerelease bump without tag")
         new_tag = new_tag.replace(prerelease=f"{args.prerelease_suffix}.0")

    # if this is the first version of the new major (so latest_tag is v1.5.0-moreactions.1)
    # then reset the prerelease counter
    if major > 0 and latest_tag is not None and latest_tag.major < new_tag.major:
        print ("-> reset prerelease")
        new_tag = new_tag.replace(prerelease=f"{args.prerelease_suffix}.0") 

    # generate the string version for output
    new_tag_str = f"{new_tag}"
    # prepend the v if enabled
    if len(args.with_v) > 0 and args.with_v == "true":
        new_tag_str = f"v{new_tag_str}"

    print("NEXT TAG DATA")
    print(f"repository_root={args.repository_root}")
    print(f"default_branch={args.default_branch}")
    print(f"prerelease={args.prerelease}")
    print(f"prerelease_processed={is_prerelease}")
    print(f"default_bump={args.default_bump}")
    print(f"last_release={args.last_release}")
    print(f"last_release_processed={last_release}")
    print(f"lastest_tag={args.latest_tag}")
    print(f"starting_tag={starting_tag}")
    # needs to be last for bash tests
    print(f"next_tag={new_tag_str}")    

    if 'GITHUB_OUTPUT' in os.environ:
        print("Pushing to GitHub Output")
        with open(os.environ['GITHUB_OUTPUT'], 'a') as fh:
            print(f"repository_root={args.repository_root}", file=fh)
            print(f"default_branch={args.default_branch}", file=fh)
            print(f"prerelease={args.prerelease}", file=fh)
            print(f"prerelease_processed={is_prerelease}", file=fh)
            print(f"default_bump={args.default_bump}", file=fh)
            print(f"last_release={args.last_release}", file=fh)
            print(f"last_release_processed={last_release}", file=fh)
            print(f"lastest_tag={args.latest_tag}", file=fh)
            print(f"starting_tag={starting_tag}", file=fh)
            print(f"next_tag={new_tag_str}", file=fh)  
            
    


if __name__ == "__main__":
    main()