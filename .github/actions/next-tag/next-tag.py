from git import Repo, Git
import os
import re
import semver
import argparse
from semver.version import Version
import importlib.util

## LOCAL IMPORTS
# up 4 levels to root or repo
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
# output helper
out_mod = importlib.util.spec_from_file_location("outputhelper", app_root_dir + '/app/python/outputhelper.py')
oh = importlib.util.module_from_spec(out_mod)  
out_mod.loader.exec_module(oh)
# semver helper
sv_mod = importlib.util.spec_from_file_location("semverhelper", app_root_dir + '/app/python/semverhelper.py')
sv = importlib.util.module_from_spec(sv_mod)  
sv_mod.loader.exec_module(sv)



def arg_parser() -> argparse.ArgumentParser:
    parser = argparse.ArgumentParser("next-tag")
    parser.add_argument('--test_file', default="", help="trigger the use of a test file for list of tags")
    
    parser.add_argument('--repository_root', default="./", help="Path to root of repository")    
    parser.add_argument('--commitish_a', default="", help="Commit-ish used to compare in log to look for triggers")    
    parser.add_argument('--commitish_b', default="", help="Commit-ish used to compare in log to look for triggers")    
    
    parser.add_argument('--prerelease', default="", help="If set, then this is a pre-release")    
    parser.add_argument("--prerelease_suffix", default="beta", help="Prerelease naming")
    
    parser.add_argument("--latest_tag", default="", help="Last tag")
    parser.add_argument("--last_release", default="", help="Last release var tag")
    parser.add_argument("--with_v", default="false", help="apply prefix to the new tag")
    parser.add_argument("--default_bump", default="minor", help="If there are no triggers in commits, bump by this")
    
    return parser


def split_commits_from_lines(lines):
    split_lines = [[str(i) for i in line.strip().split(" ", 1)] for line in lines]
    return split_lines

def get_commits(repo_root, commitish_a, commitish_b, test, test_file):
    lines = []    
    newline='⇥'
    #use test data
    if test == True and len(test_file) > 0 :
        print("Commits: Using test data")
        with open(test_file) as f:
            line = "".join([l.rstrip("\n") for l in f])            
            lines = list(filter(None, line.split(newline) ))            
    else:
        g = Git(repo_root) 
        r = Repo(repo_root)
        print(f"Commits: Using repository data: {repo_root}")
        print(f"Checking out [{commitish_a}]")
        r.git.checkout(commitish_a)
        print(f"Checking out [{commitish_b}]")
        r.git.checkout(commitish_b)
        print(f"Getting commits between [{commitish_b}]...[{commitish_a}]")
        # add a ~ to the start of each commit for easier splitting 
        # instead of new lines, as commit messages can have many lines
        log_items = g.log(f"--pretty=format:{newline}%h %s%n%b%-", f"{commitish_b}...{commitish_a}")        
        lines = [line for line in log_items.split(newline) if line.strip()]       
    commits = split_commits_from_lines( lines )    
    return commits

def get_increaments(commits:list, default_bump:str) -> tuple:
    majors=1 if default_bump == "major" else 0
    minors=1 if default_bump == "minor" else 0
    patches=1 if default_bump == "patch" else 0
    for c in commits:
        majors = majors + 1 if "#major" in c[1] else majors
        minors = minors + 1 if "#minor" in c[1] else minors
        patches = patches + 1 if "#patch" in c[1] else patches

    return majors, minors, patches

def run(
        test: bool,        
        last_release: Version,
        latest_tag: Version,
        prerelease: str,
        prerelease_suffix: str,
        default_bump: str,
        with_v:bool,
        commits:list
) -> dict:
    
    is_prerelease = (len(prerelease) > 0) 
    major, minor, patch = get_increaments(commits, default_bump)

    tag = None
    # work out base tag
    if is_prerelease:        
        tag = latest_tag if latest_tag is not None else last_release                
    else:
        tag = last_release
    
    print(f"tag is set: [{tag}]")

    new_tag = tag 
    # work out what to bump
    if major > 0:
        # Last release of v1.4.0
        # is a prerelease
        # has a major flag
        # => v2.0.0-beta.0
        if is_prerelease and tag.major <= last_release.major:
            new_tag = new_tag.bump_major().replace(prerelease = f"{prerelease_suffix}.0")
        # Last release of v1.4.0
        # is a prerelease
        # has a major flag
        # has latest_tag of v2.0.0-beta.1
        # => v2.0.0-beta.2
        elif is_prerelease and latest_tag is not None:
            new_tag = new_tag.bump_prerelease()
        # Last release of v2.0.0
        # is a prerelease
        # has a major flag
        # => v2.0.0-beta.0
        elif is_prerelease:
            new_tag = new_tag.replace(prerelease = f"{prerelease_suffix}.0")
        # Last release of v2.0.0
        # has a major flag
        # => v3.0.0
        else:
            print("major 4")
            new_tag = new_tag.bump_major()

    elif minor > 0:
        # Last release of v1.4.0
        # is a prerelease
        # has a minor flag
        # has latest_tag of v1.5.0-beta.0
        # => v1.5.0-beta.1
        if is_prerelease and latest_tag is not None:            
            new_tag = new_tag.bump_prerelease()
        # Last release of v1.4.0
        # is a prerelease
        # has a minor flag
        # => v1.5.0-beta.0
        elif is_prerelease:
            new_tag = new_tag.bump_minor().replace(prerelease = f"{prerelease_suffix}.0")        
        # Last release of v1.4.0
        # has a minor flag
        # => v1.5.0
        else:
            new_tag = new_tag.bump_minor()
    elif patch > 0:
        # Last release of v1.4.0
        # is a prerelease
        # has a minor flag
        # has latest_tag of v1.4.1-beta.0
        # => v1.4.1-beta.1
        if is_prerelease and latest_tag is not None:
            new_tag = new_tag.bump_prerelease()     
        # Last release of v1.4.0
        # is a prerelease
        # has a minor flag
        # => v1.4.1-beta.0       
        elif is_prerelease:
            new_tag = new_tag.bump_patch().replace(prerelease = f"{prerelease_suffix}.0")
        # Last release of v1.4.0        
        # has a minor flag
        # => v1.4.1
        else:
            new_tag = new_tag.bump_patch()

    # generate the string version for output
    new_tag_str = f"{new_tag}"
    # prepend the v if enabled
    if len(with_v) > 0 and with_v == "true":
        new_tag_str = f"v{new_tag_str}"


    return {
        'default_bump': default_bump,
        'majors': major,
        'minors': minor,
        'patches': patch,
        'prerelease_processed': is_prerelease,
        'last_release_processed': last_release,
        'next_tag': f"{new_tag_str}",
    }


def main():
    # get the args
    args = arg_parser().parse_args()

    # get the last release, default to 0.0.0
    lr = sv.SemverHelper(args.last_release)
    last_release = lr.parse("0.0.0")
    # get the latest_tag
    lt = sv.SemverHelper(args.latest_tag)
    latest_tag = lt.parse()

    test = (len(os.getenv("RUN_AS_TEST")) > 0)
    test_file = args.test_file

    commits = get_commits(args.repository_root, args.commitish_a, args.commitish_b, test, test_file)

    config = {
        'test': test,
        'test_file': args.test_file,
        'repository_root': args.repository_root,
        'commitish_a': args.commitish_a,
        'commitish_b': args.commitish_b,
        'prerelease': args.prerelease,        
        'last_release': args.last_release,
        'latest_tag': args.latest_tag
    }

    res = run( 
        test= len(os.getenv("RUN_AS_TEST")) > 0,
        last_release=last_release,
        latest_tag=latest_tag,
        prerelease=args.prerelease,
        prerelease_suffix=args.prerelease_suffix,        
        default_bump=args.default_bump,
        with_v=args.with_v,
        commits=commits
    )

    outputs = (config | res)
    print("NEXT TAG DATA")    
    o = oh.OutputHelper(('GITHUB_OUTPUT' in os.environ))   
    o.out(outputs)
            
    
if __name__ == "__main__":
    main()