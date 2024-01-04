from git import Repo, Git
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
# str helper
st_mod = importlib.util.spec_from_file_location("strhelper", app_root_dir + '/app/python/strhelper.py')
st = importlib.util.module_from_spec(st_mod)
st_mod.loader.exec_module(st)



def arg_parser() -> argparse.ArgumentParser:
    parser = argparse.ArgumentParser("next-tag")
    parser.add_argument('--repository_root', default="./", help="Path to root of repository")
    parser.add_argument('--commitish_a', default="", help="Commit-ish used to compare in log to look for triggers")
    parser.add_argument('--commitish_b', default="", help="Commit-ish used to compare in log to look for triggers")

    parser.add_argument('--prerelease', default="", help="If set, then this is a pre-release")
    parser.add_argument("--prerelease_suffix", default="beta", help="Prerelease naming")

    parser.add_argument("--latest_tag", default="", help="Last tag")
    parser.add_argument("--last_release", default="", help="Last release var tag")
    parser.add_argument("--with_v", default="false", help="apply prefix to the new tag")
    parser.add_argument("--default_bump", default="minor", help="If there are no triggers in commits, bump by this")

    return parser




def run(
        test: bool,
        last_release: Version|str|None,
        latest_tag: Version|str|None,
        prerelease: str|bool|None,
        prerelease_suffix: str,
        default_bump: str,
        with_v:bool|str,
        commits:list
) -> dict:

     # get the last release version from a string, default to 0.0.0
    if type(last_release) is str or last_release is None:
        print("converting last_release from string or None")
        last_release:Version = sv.SemverHelper(last_release).parse("0.0.0")
    # convert latest tag as well
    if type(latest_tag) is str:
        print("converting latest_tag from string")
        latest_tag:Version = sv.SemverHelper(latest_tag).parse()


    # allow bool True|False as well as string values
    is_prerelease = st.str_to_bool(prerelease)
    with_v:bool = st.str_to_bool(with_v)
    # find the major, minor, pathc bump counters from the commits
    major, minor, patch = ghm.GitHelper.find_bumps_from_commits(commits, default_bump)

    new_tag = sv.next_tag(
        major_bump=major,
        minor_bump=minor,
        patch_bump=patch,
        is_prerelease=is_prerelease,
        prerelease_suffix=prerelease_suffix,
        latest_tag=latest_tag,
        last_release=last_release
    )

    # generate the string version for output
    new_tag_str = f"{new_tag}"

    # prepend the v if enabled
    if with_v == True:
        new_tag_str = f"v{new_tag_str}"

    return {
        'default_bump': default_bump,
        'majors': major,
        'minors': minor,
        'patches': patch,
        'prerelease_processed': is_prerelease,
        'last_release_processed': last_release,
        'next_tag': f"{new_tag_str}",
    }


def main():
    # get the args
    args = arg_parser().parse_args()

    # get the last release, default to 0.0.0
    lr = sv.SemverHelper(args.last_release)
    last_release = lr.parse("0.0.0")
    # get the latest_tag
    lt = sv.SemverHelper(args.latest_tag)
    latest_tag = lt.parse()

    test = (len(os.getenv("RUN_AS_TEST")) > 0)
    r = ghm.GitHelper(args.repository_root)
    commits:list = r.commits(args.commitish_a, args.commitish_b)

    print(commits, sep="\n")

    config = {
        'test': test,
        'repository_root': args.repository_root,
        'commitish_a': args.commitish_a,
        'commitish_b': args.commitish_b,
        'prerelease': args.prerelease,
        'last_release': args.last_release,
        'latest_tag': args.latest_tag
    }

    res = run(
        test= len(os.getenv("RUN_AS_TEST")) > 0,
        last_release=last_release,
        latest_tag=latest_tag,
        prerelease=args.prerelease,
        prerelease_suffix=args.prerelease_suffix,
        default_bump=args.default_bump,
        with_v=args.with_v,
        commits=commits
    )

    outputs = (config | res)
    print("NEXT TAG DATA")
    o = oh.OutputHelper(('GITHUB_OUTPUT' in os.environ))
    o.out(outputs)


if __name__ == "__main__":
    main()
