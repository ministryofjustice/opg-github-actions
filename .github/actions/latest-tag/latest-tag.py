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
# string helper
st_mod = importlib.util.spec_from_file_location("strhelper", app_root_dir + '/app/python/strhelper.py')
st = importlib.util.module_from_spec(st_mod)
st_mod.loader.exec_module(st)


def arg_parser() -> argparse.ArgumentParser:
    parser = argparse.ArgumentParser("latest-tag")
    parser.add_argument('--repository_root', default="./", help="Path to root of repository")
    parser.add_argument('--prerelease', default="", help="If set, then this is a pre-release. Can be overridden if branch_name matches a release_branches item.")
    parser.add_argument("--prerelease_suffix", default="beta", help="Prerelease naming")
    parser.add_argument("--branch_name", required=True, help="Current branch name. Used to double check if this is a release or not.")
    parser.add_argument('--release_branches', default="main,master", help="List of branches that are considered a release")
    return parser


def run(
        tag_list:list,
        branch_name:str,
        release_branches:str,
        prerelease:bool|str,
        prerelease_suffix:str) -> dict:
    """
    Use the data passed to determine the following states:
        - test:bool
        - prerelease:bool
        - last_release: str
        - latest: str
    In order to do that we first confirm if the prerelease param passed is accurate
    (might be passsed true but on main branch).

    Convert the tag_list from a set of strings to a dict with its original value as key
    and the value being a Version instance, returning only semver valid ones.

    Fetch all releases from the set of tags and the last one (natural ordering)

    If this is a prerelease, then find all other prerelease that are for this branch and
    then determine the last one of those.

    Return the data.

    """
    # set these to an empty string by default
    last_release = ""
    latest = ""
    # convert prerelease to a bool
    prerelease:bool = st.str_to_bool(prerelease)
    # ensure prerelease_calculated true if this is a release branch
    prerelease_calculated:bool = sv.is_prerelease(branch_name, release_branches, prerelease)
    # convert list of strings to a dict of str->Version of semver releases only
    tags:dict = sv.list_to_dict(tag_list)
    # get the releases
    releases:dict = sv.releases(tag_list)
    # get last release
    if len(releases) > 0:
        last_release:Version|None = sv.max_version(releases)

    # if pre release, find all that match the semver pattern with this suffix
    matching:dict = {}
    if prerelease_calculated:
        matching = sv.prereleases_filtered(tags, f"{prerelease_suffix}.[0-9]+$")
    else:
        matching = releases

    # fetch the latest tag
    if len(matching) > 0:
        latest:Version|None = sv.max_version(matching)

    return {
        'prerelease_argument': prerelease,
        'prerelease_calculated': prerelease_calculated,
        'prerelease_suffix': prerelease_suffix,
        'latest': f"{latest}",
        'last_release': f"{last_release}"
    }


def main():
    args = arg_parser().parse_args()

    r = ghm.GitHelper(args.repository_root)
    all_tags = r.tags("--list")

    outputs = run(
        all_tags,
        args.branch_name,
        args.release_branches,
        args.prerelease,
        args.prerelease_suffix
    )

    print("# latest-tag outputs:")
    o = oh.OutputHelper(('GITHUB_OUTPUT' in os.environ))
    o.out(outputs)


if __name__ == "__main__":
    main()
