from git import Repo
import os
import re
import semver
import argparse
from semver.version import Version
import importlib.util

## LOCAL IMPORTS
# up 4 levels to root or repo
app_root_dir = os.path.dirname(
    os.path.dirname(
        os.path.dirname(
            os.path.dirname(os.path.realpath(__file__))
        )
    )
)
# git helper
git_mod = importlib.util.spec_from_file_location("githelper", app_root_dir + '/app/python/githelper.py')
ghm = importlib.util.module_from_spec(git_mod)
git_mod.loader.exec_module(ghm)
# output helper
out_mod = importlib.util.spec_from_file_location("outputhelper", app_root_dir + '/app/python/outputhelper.py')
oh = importlib.util.module_from_spec(out_mod)
out_mod.loader.exec_module(oh)
# semver helper
sv_mod = importlib.util.spec_from_file_location("semverhelper", app_root_dir + '/app/python/semverhelper.py')
sv = importlib.util.module_from_spec(sv_mod)
sv_mod.loader.exec_module(sv)


def arg_parser() -> argparse.ArgumentParser:
    parser = argparse.ArgumentParser("latest-tag")
    parser.add_argument('--repository_root', default="./", help="Path to root of repository")

    parser.add_argument('--prerelease', default="", help="If set, then this is a pre-release. Can be overridden if branch_name matches a release_branches item.")
    parser.add_argument("--prerelease_suffix", default="beta", help="Prerelease naming")

    parser.add_argument("--branch_name", required=True, help="Current branch name. Used to double check if this is a release or not.")
    parser.add_argument('--release_branches', default="main,master", help="List of branches that are considered a release")
    return parser




def is_prerelease(prerelease:bool, branch_name:str, release_branches:str) -> bool:
    """Check if the branch is a listed release branch, if so overwrite and flag as release."""
    if branch_name in release_branches.split(","):
        return False
    return prerelease


def run(
        test:bool,
        tags:list,
        branch_name:str,
        release_branches:str,
        prerelease:bool|str,
        prerelease_suffix:str) -> dict:

    # convert prerelease to a bool
    if type(prerelease) is str:
        if len(prerelease) > 0 and prerelease.lower() == "true":
            prerelease = True
        else:
            prerelease = False

    prerelease_by_branch = is_prerelease(prerelease, branch_name, release_branches)

    tags = sv.SemverHelper.list_to_dict(tags)

    last_release = ""
    latest = ""

    # get the releases, and track last one in particular
    releases = {k:v for k,v in tags.items() if v.prerelease is None}
    if len(releases) > 0:
        max_release_val = max(releases.values())
        release_items = [k for k,v in releases.items() if f"{v}" == max_release_val]
        last_release = release_items.pop()

    # if pre release, find all that match the semver pattern with this suffix
    matching = []
    if prerelease_by_branch:
        pattern = re.compile(f"{prerelease_suffix}.[0-9]+$")
        matching = {k:v for k,v in tags.items() if pattern.match(f"{v.prerelease}")}
    else:
        matching = releases

    # fetch the latest tag
    if len(matching) > 0:
        latest_val = max(matching.values())
        latest_items = [k for k,v in matching.items() if f"{v}" == latest_val]
        latest = latest_items.pop()

    return {
        'test': test,
        'prerelease_argument': prerelease,
        'prerelease_calculated': prerelease_by_branch,
        'prerelease_suffix': prerelease_suffix,
        'latest': f"{latest}",
        'last_release': f"{last_release}"
    }


def main():
    args = arg_parser().parse_args()

    r = ghm.GitHelper(args.repository_root)
    all_tags = r.tags("--list")

    outputs = run(
        len(os.getenv("RUN_AS_TEST")) > 0,
        all_tags,
        args.branch_name,
        args.release_branches,
        args.prerelease,
        args.prerelease_suffix
    )

    print("LATEST TAG DATA")
    o = oh.OutputHelper(('GITHUB_OUTPUT' in os.environ))
    o.out(outputs)


if __name__ == "__main__":
    main()
