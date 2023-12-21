from git import Repo
import os
import re
import semver
import argparse
from semver.version import Version


def arg_parser() -> argparse.ArgumentParser:
    parser = argparse.ArgumentParser("parse-tags")
    parser.add_argument('--test_file', default="", help="trigger the use of a test file for list of tags")
    parser.add_argument('--repository_root', default="./", help="Path to root of repository")
    parser.add_argument('--prerelease', default="", help="If set, then this is a pre-release")
    parser.add_argument("--prerelease_suffix", default="beta", help="Prerelease naming")
    return parser


def tags_from_file(file):
    lines=[[str(i) for i in line.strip().split(" ", 1)] for line in open(file).readlines()]
    tags=[]    
    for line in lines:
        s = line[0]
        t = (s[1:] if s.startswith('v') else s)    
        if Version.is_valid(t):            
            tags.append(Version.parse(t))
    return tags

def semver_list(tags):
    semver_tags = []
    for tag in tags:
        s = f"{tag}"            
        t = (s[1:] if s.startswith('v') else s)    
        if Version.is_valid(t):        
            semver_tags.append(Version.parse(t))
    return semver_tags

def main():
    args = arg_parser().parse_args()
    test_file = args.test_file
    repo_root = args.repository_root
    prerelease = (len(args.prerelease) > 0)
    prerelease_suffix = args.prerelease_suffix

    # use test content
    if os.getenv("RUN_AS_TEST") and len(test_file) > 0:
        tags = tags_from_file(test_file)
    else:
        repo = Repo(repo_root)    
        tags = semver_list(repo.tags)
        
    # if pre release, find set matching that pattern
    matching = []
    if prerelease:
        pattern = re.compile(f"{prerelease_suffix}.[0-9]+$")
        matching = list( filter(lambda t:( pattern.match(f"{t.prerelease}") ), tags ) )
    else:
        matching = list( filter(lambda t:( t.prerelease is None ), tags ) )

    last = ""
    if len(matching) > 0:
        last = max(matching)

    # summary for shell
    print(f"prerelease={prerelease}")
    print(f"prerelease_suffix={prerelease_suffix}")
    print(f"latest={last}")

    if 'GITHUB_OUTPUT' in os.environ:
        print("Pushing to GitHub Output")
        with open(os.environ['GITHUB_OUTPUT'], 'a') as fh:
            print(f'prerelease={prerelease}', file=fh)
            print(f'prerelease_suffix={prerelease_suffix}', file=fh)
            print(f'latest={last}', file=fh)
    

if __name__ == "__main__":
    main()