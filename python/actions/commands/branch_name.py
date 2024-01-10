#!/usr/bin/env python3
import argparse
import os
import json
from actions.common import githelper as gh, outputhelper as oh, strhelper as st

def arg_parser() -> argparse.ArgumentParser:
    parser = argparse.ArgumentParser("branch-name")
    parser.add_argument('--event_name', default="", required=True, help="Github Action event_name")
    parser.add_argument('--event_data_file', default="", required=True, help="File containing json payload.")
    return parser

def run(
        event_name:str,
        event_data:dict,
        length:int = 12
) -> dict:
    """
    """
    branch_name, source_commitish, destination_commitish = gh.github_branch_data(event_name, event_data)

    full_length = st.safe(branch_name)
    return {
        'event_name': event_name,
        'source_commitish': source_commitish,
        'destination_commitish': destination_commitish,
        'branch_name': branch_name,
        'full_length': full_length,
        'safe': full_length[0:length]
    }


def main():
    # get the args
    args = arg_parser().parse_args()
    event_data:dict = {}
    with open(args.event_data_file, 'r') as f:
        event_data = json.load(f)

    output = run(
        event_name=args.event_name,
        event_data=event_data,
    )

    print("# branch-name outputs:")
    o = oh.OutputHelper(('GITHUB_OUTPUT' in os.environ))
    o.out(output)


if __name__ == "__main__":
    main()
