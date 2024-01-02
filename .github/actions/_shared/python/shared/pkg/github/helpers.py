import string
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


def header(to_github) -> None:
    print("## Test Information")
    print("| A | condition | B | Pass |")
    print("| --- | --- | --- | --- |")
    if to_github:
        with open(os.environ['GITHUB_OUTPUT'], 'a') as fh:
            print("## Test Information", file=fh)
            print("| A | condition | B | Pass |", file=fh)
            print("| --- | --- | --- | --- |", file=fh)

def result(a, condition, b, passed, to_github) -> None:
    if passed:
        passing(a, condition, b, to_github)
    else:
        failing(a, condition, b, to_github)

def failing(a, condition, b, to_github) -> None:
    result_line(a, condition, b, to_github, "❌")

def passing(a, condition, b, to_github) -> None:
    result_line(a, condition, b, to_github, "✅")
    
def result_line(a, condition, b, to_github, char) -> None:
    print (f"| {a} | {condition} | {b} | {char}  |")
    if to_github:
        with open(os.environ['GITHUB_OUTPUT'], 'a') as fh:
            print (f"| {a} | {condition} | {b} | {char}  |", file=fh)
