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

def to_valid_dict(str) -> dict|None:
    t = (str[1:] if str.startswith('v') else s)    
    if Version.is_valid(t):            
        return {"raw": str, "tag": Version.parse(t)}
    return None

def tags_from_file(file):
    lines=[[str(i) for i in line.strip().split(" ", 1)] for line in open(file).readlines()]
    tags=[]    
    for line in lines:
        d = to_valid_dict(line[0])
        if d is not None:
            tags.append(d)        
    return tags

def semver_list(tags):
    semver_tags = []
    for tag in tags:
        d = to_valid_dict(f"{tag}")
        if d is not None:
            semver_tags.append(d)
    return semver_tags

def main():
    args = arg_parser().parse_args()
    test_file = args.test_file
    repo_root = args.repository_root
    prerelease = (len(args.prerelease) > 0)
    prerelease_suffix = args.prerelease_suffix

    # use test content
    test = os.getenv("RUN_AS_TEST")
    is_test = False
    if len(test) > 0 and len(test_file) > 0:
        is_test = True
        tags = tags_from_file(test_file)
    else:
        repo = Repo(repo_root)    
        tags = semver_list(repo.tags)
        
    # if pre release, find set matching that pattern
    matching = []
    if prerelease:
        pattern = re.compile(f"{prerelease_suffix}.[0-9]+$")
        matching = list( filter(lambda t:( pattern.match(f"{t['tag'].prerelease}") ), tags ) )
    else:
        matching = list( filter(lambda t:( t['tag'].prerelease is None ), tags ) )

    last = ""
    if len(matching) > 0:
        # use the raw tag
        last = max(matching, key=lambda x: x['raw'] )
        last = last.get("raw")        

    # summary for shell
    print(f"test={is_test}")
    print(f"prerelease={prerelease}")
    print(f"prerelease_suffix={prerelease_suffix}")
    print(f"latest={last}")

    if 'GITHUB_OUTPUT' in os.environ:
        print("Pushing to GitHub Output")
        with open(os.environ['GITHUB_OUTPUT'], 'a') as fh:
            print(f'test={is_test}', file=fh)
            print(f'prerelease={prerelease}', file=fh)
            print(f'prerelease_suffix={prerelease_suffix}', file=fh)
            print(f'latest={last}', file=fh)
    

if __name__ == "__main__":
    main()