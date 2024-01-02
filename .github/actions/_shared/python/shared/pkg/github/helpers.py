import string
import os
from git import Repo, Git
from natsort import natsorted

def tags(repo, param:str) -> list:
    """
    Use the repo object to fetch the tag data, split that
    by new lines and return a list
        - `param` to allow options (such as --list or --points-at)
    """
    all = list( repo.git.tag(param).split("\n") )
    return natsorted(all)


def header(fh) -> None:
    print("## Test")
    print("| A | condition | B | Pass |")
    print("| --- | --- | --- | --- |")

    fh.writelines(["## Test Information\n", "| A | condition | B | Pass |\n", "| --- | --- | --- | --- |\n"])
    

def result(a, condition, b, passed, fh) -> None:
    if passed:
        passing(a, condition, b, fh)
    else:
        failing(a, condition, b, fh)

def failing(a, condition, b, fh) -> None:
    result_line(a, condition, b, fh, "❌")

def passing(a, condition, b, fh) -> None:
    result_line(a, condition, b, fh, "✅")
    
def result_line(a, condition, b, fh, char) -> None:
    print (f"| {a} | {condition} | {b} | {char}  |")
    fh.writelines ([f"| {a} | {condition} | {b} | {char}  |\n"])
    
