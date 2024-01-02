import string
from git import Repo, Git
from natsort import natsorted

def tags(repo, param:str) -> list:
    all = list( repo.git.tag(param).split("\n") )
    return natsorted(all)
