from git import Repo
import os
import re
import semver
import argparse
from semver.version import Version


def arg_parser() -> argparse.ArgumentParser:
    parser = argparse.ArgumentParser("parse-tags")
    parser.add_argument('--test_file', default="", help="trigger the use of a test file for list of tags. Requires ENV RUN_AS_TEST to be set as well.")
    
    parser.add_argument('--repository_root', default="./", help="Path to root of repository")
    
    parser.add_argument('--prerelease', default="", help="If set, then this is a pre-release. Can be overridden if branch_name matches a release_branches item.")
    parser.add_argument("--prerelease_suffix", default="beta", help="Prerelease naming")
    
    parser.add_argument("--branch_name", required=True, help="Current branch name. Used to double check if this is a release or not.")
    parser.add_argument('--release_branches', default="main,master", help="List of branches that are considered a release")
    return parser

def to_valid_dict(tag) -> dict|None:
    t = (tag[1:] if tag.startswith('v') else tag)    
    if Version.is_valid(t):            
        return {tag: Version.parse(t)}
    return None

def tags_from_file(file) -> dict:
    lines=[[str(i) for i in line.strip().split(" ", 1)] for line in open(file).readlines()]
    tags={}    
    for line in lines:
        d = to_valid_dict(line[0])
        if d is not None:
            tags.update(d)        
    return tags

def semver_list(tags) -> dict:
    semver_tags = {}
    for tag in tags:
        d = to_valid_dict(f"{tag}")
        if d is not None:
            semver_tags.update(d)
    return semver_tags

def is_prerelease(prerelease, branch_name, release_branches) -> bool:
    if branch_name in release_branches.split(","):
        return False
    return len(prerelease) > 0


def main():
    args = arg_parser().parse_args()
    test_file = args.test_file
    repo_root = args.repository_root
    prerelease_by_branch = is_prerelease(args.prerelease, args.branch_name, args.release_branches)
    prerelease_suffix = args.prerelease_suffix

    # use test content
    test = os.getenv("RUN_AS_TEST")
    is_test = False
    if test is not None and len(test) > 0 and len(test_file) > 0:
        print("Using test data")
        is_test = True
        tags = tags_from_file(test_file)
    else:
        repo = Repo(repo_root)    
        tags = semver_list(repo.tags)
        
    last_release = ""
    latest = ""

    # get the releases, and track last one in particular
    releases = {k:v for k,v in tags.items() if v.prerelease is None} 
    if len(releases) > 0:        
        max_release_val = max(releases.values())
        release_items = [k for k,v in releases.items() if f"{v}" == max_release_val]
        last_release = release_items.pop()
        
    # if pre release, find set matching that pattern
    matching = []
    if prerelease_by_branch:
        pattern = re.compile(f"{prerelease_suffix}.[0-9]+$")
        matching = {k:v for k,v in tags.items() if pattern.match(f"{v.prerelease}")} 
    else:
        matching = releases

    # print(*matching, sep="\n")
    if len(matching) > 0:        
        latest_val = max(matching.values())        
        latest_items = [k for k,v in matching.items() if f"{v}" == latest_val]
        latest = latest_items.pop()

    # summary for shell
    print(f"test={is_test}")
    print(f"prerelease_argument={args.prerelease}")
    print(f"prerelease_calculated={prerelease_by_branch}")
    print(f"prerelease_suffix={prerelease_suffix}")
    print(f"latest={latest}")
    print(f"last_release={last_release}")

    if 'GITHUB_OUTPUT' in os.environ:
        print("Pushing to GitHub Output")
        with open(os.environ['GITHUB_OUTPUT'], 'a') as fh:
            print(f'test={is_test}', file=fh)
            print(f'prerelease={prerelease_by_branch}', file=fh)
            print(f'prerelease_suffix={prerelease_suffix}', file=fh)
            print(f'latest={latest}', file=fh)
            print(f'last_release={last_release}', file=fh)
    

if __name__ == "__main__":
    main()