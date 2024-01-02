from git import Repo
import os
import re
import semver
import argparse
from semver.version import Version
import importlib.util

# local imports
parent_dir_name = os.path.dirname(os.path.dirname(os.path.realpath(__file__)))
# load cli helper
cli_mod = importlib.util.spec_from_file_location("clih", parent_dir_name + '/_shared/python/shared/pkg/cli/helpers.py')
cli = importlib.util.module_from_spec(cli_mod)  
cli_mod.loader.exec_module(cli)
# load semver helpers
semver_mod = importlib.util.spec_from_file_location("semverh", parent_dir_name + '/_shared/python/shared/pkg/semver/helpers.py')
svh = importlib.util.module_from_spec(semver_mod)  
semver_mod.loader.exec_module(svh)



def arg_parser() -> argparse.ArgumentParser:
    parser = argparse.ArgumentParser("latest-tag")
    parser.add_argument('--test_file', default="", help="trigger the use of a test file for list of tags. Requires ENV RUN_AS_TEST to be set as well.")
    
    parser.add_argument('--repository_root', default="./", help="Path to root of repository")
    
    parser.add_argument('--prerelease', default="", help="If set, then this is a pre-release. Can be overridden if branch_name matches a release_branches item.")
    parser.add_argument("--prerelease_suffix", default="beta", help="Prerelease naming")
    
    parser.add_argument("--branch_name", required=True, help="Current branch name. Used to double check if this is a release or not.")
    parser.add_argument('--release_branches', default="main,master", help="List of branches that are considered a release")
    return parser


def tags_from_file(file) -> dict:
    lines=[[str(i) for i in line.strip().split(" ", 1)] for line in open(file).readlines()]
    tags={}    
    for line in lines:
        d = svh.to_valid_dict(line[0])
        if d is not None:
            tags.update(d)        
    return tags



def is_prerelease(prerelease, branch_name, release_branches) -> bool:
    if branch_name in release_branches.split(","):
        return False
    return len(prerelease) > 0


def run(
        test:bool, 
        test_file:str, 
        repo_root:str, 
        branch_name:str, 
        release_branches,
        prerelease, 
        prerelease_suffix:str) -> dict:
    
    prerelease_by_branch = is_prerelease(prerelease, branch_name, release_branches)
    # use test content
    is_test = False
    if test == True and len(test_file) > 0:
        print("Using test data")
        is_test = True
        tags = tags_from_file(test_file)
    else:
        repo = Repo(repo_root)    
        tags = svh.semver_list(repo.tags)
        
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

    if len(matching) > 0:        
        latest_val = max(matching.values())        
        latest_items = [k for k,v in matching.items() if f"{v}" == latest_val]
        latest = latest_items.pop()

    return {
        'test': is_test,
        'prerelease_argument': prerelease,
        'prerelease_calculated': prerelease_by_branch,
        'prerelease_suffix': prerelease_suffix,
        'latest': f"{latest}",
        'last_release': f"{last_release}"
    }


def main():
    args = arg_parser().parse_args()

    outputs = run(
        len(os.getenv("RUN_AS_TEST")) > 0,
        args.test_file,
        args.repository_root,
        args.branch_name,
        args.release_branches,
        args.prerelease,
        args.prerelease_suffix
    )

    print("LATEST TAG DATA")    
    cli.results(outputs, 'GITHUB_OUTPUT' in os.environ)


if __name__ == "__main__":
    main()