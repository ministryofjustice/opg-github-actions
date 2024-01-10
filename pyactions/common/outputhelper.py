#!/usr/bin/env python3
import os


class OutputHelper:
    """
    Handles output to stdout and file
    Used for test results and pipelines summaries
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

    def header(self, fh) -> None:
        print("## Test")
        print("| A | condition | B | Pass |")
        print("| --- | --- | --- | --- |")

        fh.writelines(["## Test Information\n", "| A | condition | B | Pass |\n", "| --- | --- | --- | --- |\n"])

    def result(self, a, condition, b, passed, fh) -> None:
        if passed:
            self.passing(a, condition, b, fh)
        else:
            self.failing(a, condition, b, fh)

    def failing(self, a, condition, b, fh) -> None:
        self.result_line(a, condition, b, fh, "❌")

    def passing(self, a, condition, b, fh) -> None:
        self.result_line(a, condition, b, fh, "✅")

    def result_line(self, a, condition, b, fh, char) -> None:
        print (f"| {a} | {condition} | {b} | {char}  |")
        fh.writelines ([f"| {a} | {condition} | {b} | {char}  |\n"])
