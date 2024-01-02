import sys
import os
import importlib.util
from git import Repo, Git
from natsort import natsorted
import semver
import argparse
import string
from semver.version import Version

# local imports
parent_dir_name = os.path.dirname(os.path.dirname(os.path.realpath(__file__)))
# load semver helpers
semver_mod = importlib.util.spec_from_file_location("semver_mod", parent_dir_name + '/_shared/python/semver.py')
svh = importlib.util.module_from_spec(semver_mod)  
semver_mod.loader.exec_module(svh)
# load github helper
gh_mod = importlib.util.spec_from_file_location("gh_mod", parent_dir_name + '/_shared/python/github.py')
gh = importlib.util.module_from_spec(gh_mod)  
gh_mod.loader.exec_module(gh)
# load cli helper
cli_mod = importlib.util.spec_from_file_location("cli_mod", parent_dir_name + '/_shared/python/cli.py')
cli = importlib.util.module_from_spec(cli_mod)  
cli_mod.loader.exec_module(cli)
# load rand helper
rand_mod = importlib.util.spec_from_file_location("rand_mod", parent_dir_name + '/_shared/python/rand.py')
rnd = importlib.util.module_from_spec(rand_mod)  
rand_mod.loader.exec_module(rnd)


def arg_parser() -> argparse.ArgumentParser:
    parser = argparse.ArgumentParser("create-tags")
    parser.add_argument('--repository_root', default="./", help="Path to root of repository")    
    parser.add_argument('--commitish', default="", help="Commit-ish ref to create the tag") 
    parser.add_argument("--tag_name", default="", help="Tag to create")    
    return parser


def generate_tag_to_create(tag_name: str, all_tags: list, valid_semver:bool, with_v:bool) -> str:
    rand_length = 3
    original_tag = tag_name
    # if this is semver, then parse and update it
    if valid_semver:
        parsed_tag = Version.parse(svh.trim_v(tag_name) if with_v else tag_name)
        while tag_name in all_tags:
            # if this is a pre-release, then we can adjust that
            if parsed_tag.prerelease is not None:
                parsed_tag = parsed_tag.replace(prerelease=f"{rnd.rand(rand_length)}.0")
                tag_name = f"v{parsed_tag}" if with_v else f"{parsed_tag}"
            # otherwise, bump version as this should be release
            else:
                parsed_tag = parsed_tag.bump_major()
                tag_name = f"v{parsed_tag}" if with_v else f"{parsed_tag}"
    # if its not a semver then tag on a random suffix
    else:
        while tag_name in all_tags:
            tag_name = f"{original_tag}.{rnd.rand(rand_length)}"
    return tag_name


def main():
    print(svh.has_v("true"))
    # get the args
    args = arg_parser().parse_args()
    repo_root = args.repository_root
    commitish = args.commitish
    tag_name = args.tag_name
    with_v = svh.has_v(tag_name)
    test = len( os.getenv("RUN_AS_TEST") ) > 0
    valid_semver = Version.parse(svh.trim_v(tag_name))

    repo = Repo(repo_root)
    # get all tags
    all_tags = gh.tags(repo, "--list")    
    # get all tags that point at this commit
    all_tags_here = gh.tags(repo, f"--points-at={commitish}")    
    
    # looks for clashing tags in the existing set
    tag_to_create = generate_tag_to_create(tag_name, all_tags, valid_semver, with_v)
    # create the tag        
    repo.git.tag(tag_to_create, commitish)
    # if this isnt a test, push to remote
    if test != True:
        print(f"Pushing {tag_to_create} to remote")
        #repo.git.push('origin', tag_to_create)

    all_tags = gh.tags(repo, "--list")    
    all_tags_here = gh.tags(repo, f"--points-at={commitish}")
    latest_tag = all_tags_here[-1]

    outputs={
        'all_tags': ','.join(all_tags),
        'all_tags_here': ','.join(all_tags_here),
        'latest_tag': latest_tag,
        'requested_tag': tag_name,
        'created_tag': tag_to_create
    }

    print("CREATE TAG DATA")    
    cli.results(outputs, 'GITHUB_OUTPUT' in os.environ)

if __name__ == "__main__":
    main()