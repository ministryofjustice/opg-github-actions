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


def get_increments(commits:list, default_bump:str) -> tuple:
    """
    Scan all fields in the commits passed looking for triggers of each type.
    Return counter of each.
    The count for default_bump starts at 1 instead of 0 to ensure something is
    always increased.
    """
    majors=1 if default_bump == "major" else 0
    minors=1 if default_bump == "minor" else 0
    patches=1 if default_bump == "patch" else 0
    for c in commits:
        # check each field in the dict
        for k in ['subject', 'notes', 'body']:
            majors = majors + 1 if "#major" in c[k] else majors
            minors = minors + 1 if "#minor" in c[k] else minors
            patches = patches + 1 if "#patch" in c[k] else patches

    return majors, minors, patches

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
        last_release = sv.SemverHelper(last_release).parse("0.0.0")
    # convert latest tag as well
    if type(latest_tag) is str:
        print("converting latest_tag from string")
        latest_tag = sv.SemverHelper(latest_tag).parse()

    # allow bool True|False as well as string values
    if type(prerelease) is bool:
        is_prerelease = prerelease
    else:
        is_prerelease = (len(prerelease) > 0 and prerelease.lower() == "true") if prerelease is not None else False

    major, minor, patch = get_increments(commits, default_bump)

    tag = None
    # work out base tag
    if is_prerelease:
        tag = latest_tag if latest_tag is not None else last_release
    else:
        tag = last_release

    print(f"tag is set: [{tag}]")

    new_tag = tag
    # update the new tag
    if major > 0:
        # Last release of v1.4.0
        # is a prerelease
        # has a major flag
        # => v2.0.0-beta.0
        if is_prerelease and tag.major <= last_release.major:
            new_tag = new_tag.bump_major().replace(prerelease = f"{prerelease_suffix}.0")
        # Last release of v1.4.0
        # is a prerelease
        # has a major flag
        # has latest_tag of v2.0.0-beta.1
        # => v2.0.0-beta.2
        elif is_prerelease and latest_tag is not None:
            new_tag = new_tag.bump_prerelease()
        # Last release of v2.0.0
        # is a prerelease
        # has a major flag
        # => v2.0.0-beta.0
        elif is_prerelease:
            new_tag = new_tag.replace(prerelease = f"{prerelease_suffix}.0")
        # Last release of v2.0.0
        # has a major flag
        # => v3.0.0
        else:
            print("major 4")
            new_tag = new_tag.bump_major()
    elif minor > 0:
        # Last release of v1.4.0
        # is a prerelease
        # has a minor flag
        # has latest_tag of v1.5.0-beta.0
        # => v1.5.0-beta.1
        if is_prerelease and latest_tag is not None:
            new_tag = new_tag.bump_prerelease()
        # Last release of v1.4.0
        # is a prerelease
        # has a minor flag
        # => v1.5.0-beta.0
        elif is_prerelease:
            new_tag = new_tag.bump_minor().replace(prerelease = f"{prerelease_suffix}.0")
        # Last release of v1.4.0
        # has a minor flag
        # => v1.5.0
        else:
            new_tag = new_tag.bump_minor()
    elif patch > 0:
        # Last release of v1.4.0
        # is a prerelease
        # has a minor flag
        # has latest_tag of v1.4.1-beta.0
        # => v1.4.1-beta.1
        if is_prerelease and latest_tag is not None:
            new_tag = new_tag.bump_prerelease()
        # Last release of v1.4.0
        # is a prerelease
        # has a minor flag
        # => v1.4.1-beta.0
        elif is_prerelease:
            new_tag = new_tag.bump_patch().replace(prerelease = f"{prerelease_suffix}.0")
        # Last release of v1.4.0
        # has a minor flag
        # => v1.4.1
        else:
            new_tag = new_tag.bump_patch()

    # generate the string version for output
    new_tag_str = f"{new_tag}"

    # prepend the v if enabled
    if type(with_v) is bool and with_v == True:
        new_tag_str = f"v{new_tag_str}"
    elif type(with_v) is str and len(with_v) > 0 and with_v == "true":
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
    commits = r.commits(args.commitish_a, args.commitish_b)

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
