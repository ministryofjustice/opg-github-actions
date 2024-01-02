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
