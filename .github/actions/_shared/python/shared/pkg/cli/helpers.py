import os

def results(outputs:dict, is_github:bool):
    """
    Helper to output a dict of key/value pairs to stdout as 
    well as to github output if the env var exists
    """
    for k,v in outputs.items():
        print(f"{k}={v}")
        if is_github:
            with open(os.environ['GITHUB_OUTPUT'], 'a') as fh:
                print(f"{k}={v}", fh)
