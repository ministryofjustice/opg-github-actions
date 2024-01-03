import sys
import os
import importlib.util
from natsort import natsorted
import semver
import argparse
import string
from semver.version import Version

## LOCAL IMPORTS
# up 4 levels to root or repo?
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
# github output helper
gh_mod = importlib.util.spec_from_file_location("outputhelper", app_root_dir + '/app/python/outputhelper.py')
gh = importlib.util.module_from_spec(gh_mod)  
gh_mod.loader.exec_module(gh)


def arg_parser() -> argparse.ArgumentParser:
    parser = argparse.ArgumentParser("create-tags")
    parser.add_argument('--repository_root', default="./", help="Path to root of repository")    
    parser.add_argument('--commitish', default="", help="Commit-ish ref to create the tag") 
    parser.add_argument("--tag_name", default="", help="Tag to create")    
    return parser


def run(repo_root:str, commitish:str, tag_name:str, test:bool) -> dict:
    """Run groups all calls together to allow testing from other file"""
     # get all tags    
    r = ghm.GitHelper(repo_root)
    all_tags = r.tags("--list")
    all_tags_here = r.tags(f"--points-at={commitish}") 
    
    # looks for clashing tags in the existing set
    tag_to_create = r.tag_to_create(tag_name, all_tags)
    # create the tag 
    r.create_tag(tag_to_create, commitish, (test != True))
    # refresh the tags for returning
    all_tags = r.tags("--list")
    all_tags_here = r.tags(f"--points-at={commitish}") 
    latest_tag = all_tags_here[-1]

    outputs={
        'all_tags': ','.join(all_tags),
        'all_tags_here': ','.join(all_tags_here),
        'latest_tag': f"{latest_tag}",
        'requested_tag': f"{tag_name}",
        'created_tag': f"{tag_to_create}"
    }
    return outputs

def main():
    # get the args
    args = arg_parser().parse_args()
    # call the runner directly
    outputs = run(
        args.repository_root, 
        args.commitish, 
        args.tag_name,
        len( os.getenv("RUN_AS_TEST") ) > 0
    )
    print("CREATE TAG DATA") 
    g = gh.OutputHelper(('GITHUB_OUTPUT' in os.environ))   
    g.out(outputs)

if __name__ == "__main__":
    main()