import os

def results(outputs:dict, is_github:bool):
    for k,v in outputs.items():
        print(f"{k}={v}")
        if is_github:
            with open(os.environ['GITHUB_OUTPUT'], 'a') as fh:
                print(f"{k}={v}", fh)
