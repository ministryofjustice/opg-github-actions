#!/usr/bin/env python3
import os
import argparse
from actions.common import githelper as ghm, outputhelper as oh

def arg_parser() -> argparse.ArgumentParser:
    parser = argparse.ArgumentParser("create-tags")
    parser.add_argument('--repository_root', default="./", required=True, help="Path to root of repository")
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
    print("# create-tag outputs:")
    g = oh.OutputHelper(('GITHUB_OUTPUT' in os.environ))
    g.out(outputs)

if __name__ == "__main__":
    main()
