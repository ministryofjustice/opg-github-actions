import argparse
import importlib.util
import os
## LOCAL IMPORTS
# up 4 levels to root or repo
app_root_dir = os.path.dirname(
    os.path.dirname(
        os.path.dirname(
            os.path.dirname(os.path.realpath(__file__))
        )
    )
)
# str helper
st_mod = importlib.util.spec_from_file_location("strhelper", app_root_dir + '/app/python/strhelper.py')
st = importlib.util.module_from_spec(st_mod)
st_mod.loader.exec_module(st)
# output helper
out_mod = importlib.util.spec_from_file_location("outputhelper", app_root_dir + '/app/python/outputhelper.py')
oh = importlib.util.module_from_spec(out_mod)
out_mod.loader.exec_module(oh)



def arg_parser() -> argparse.ArgumentParser:
    parser = argparse.ArgumentParser("safe-strings")

    parser.add_argument('--string', default="", help="string to clean")
    parser.add_argument('--suffix', default="", help="suffix")
    parser.add_argument('--length', default="", help="Max length of string")
    parser.add_argument('--conditional_match', default="", help="If the original string matches this value, then use the conditional_value directly.")
    parser.add_argument('--conditional_value', default="", help="When original matches conditional_match use this value for all other outputs directly.")

    return parser

def run(
        original:str,
        suffix:str,
        length:int|None,
        conditional_match:str|None,
        conditional_value:str|None,
) -> dict:
    """
    """
    length:int|None = st.int_or_none(length)
    to_clean:str = original
    full_length:str = st.safe(to_clean)
    safe:str = full_length

    # trim the string, adding the suffix
    if suffix and length:
        print("suffix and length")
        l = length - len(suffix)
        safe = f"{safe[0:l]}{suffix}"
    elif length:
        print("length")
        safe = f"{safe[0:length]}"
    elif suffix:
        print("suffix")
        safe = f"{safe}{suffix}"

    # if this exists, replace and overwrite
    if conditional_match == original:
        full_length = conditional_value
        safe = conditional_value

    return {
        'original': original,
        'suffix': suffix,
        'length': length,
        'conditional_match': conditional_match,
        'conditional_value': conditional_value,
        'full_length': full_length,
        'safe': safe
    }


def main():
    # get the args
    args = arg_parser().parse_args()
    output = run(
        original=args.string,
        suffix=args.suffix,
        length=st.int_or_none(args.length),
        conditional_match=args.conditional_match,
        conditional_value=args.conditional_value
    )

    print("CLEAN STRING DATA")
    o = oh.OutputHelper(('GITHUB_OUTPUT' in os.environ))
    o.out(output)


if __name__ == "__main__":
    main()
