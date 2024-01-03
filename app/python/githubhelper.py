#!/usr/bin/env python3
import os


class GitHubHelper:
    """
    """
    enabled = False

    def __init__(self, enabled: bool):
        self.enabled = enabled
    
    def out(self, outputs:dict) -> None:
        """
        Helper to output a dict of key/value pairs to stdout as 
        well as to github output if the env var exists
        """
        for k,v in outputs.items():
            print(f"{k}={v}")
            if self.enabled:
                with open(os.environ['GITHUB_OUTPUT'], 'a') as fh:
                    print(f"{k}={v}", file=fh)
