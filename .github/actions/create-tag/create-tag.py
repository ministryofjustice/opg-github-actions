import sys
import os
import importlib.util
from natsort import natsorted
import semver
import argparse
import string
from semver.version import Version

# local imports
parent_dir_name = os.path.dirname(os.path.dirname(os.path.realpath(__file__)))
# load cli helper
cli_mod = importlib.util.spec_from_file_location("clih", parent_dir_name + '/_shared/python/shared/pkg/cli/helpers.py')
cli = importlib.util.module_from_spec(cli_mod)  
cli_mod.loader.exec_module(cli)

# load rand helper
tag_mod = importlib.util.spec_from_file_location("tagh", parent_dir_name + '/_shared/python/shared/pkg/tag/helpers.py')
taghelper = importlib.util.module_from_spec(tag_mod)  
tag_mod.loader.exec_module(taghelper)


def arg_parser() -> argparse.ArgumentParser:
    parser = argparse.ArgumentParser("create-tags")
    parser.add_argument('--repository_root', default="./", help="Path to root of repository")    
    parser.add_argument('--commitish', default="", help="Commit-ish ref to create the tag") 
    parser.add_argument("--tag_name", default="", help="Tag to create")    
    return parser


def main():
    # get the args
    args = arg_parser().parse_args()
    repo_root = args.repository_root
    commitish = args.commitish
    tag_name = args.tag_name
    
    test = len( os.getenv("RUN_AS_TEST") ) > 0

    # get all tags    
    all_tags = taghelper.repo_tags(repo_root, "--list")  
    all_tags_here = taghelper.repo_tags(repo_root, f"--points-at={commitish}") 
    
    # looks for clashing tags in the existing set
    tag_to_create = taghelper.generate_tag_to_create(tag_name, all_tags)
    # create the tag 
    taghelper.create_tag(repo_root, commitish, tag_to_create, (test != True))

    all_tags = taghelper.repo_tags(repo_root, "--list")  
    all_tags_here = taghelper.repo_tags(repo_root, f"--points-at={commitish}") 
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