import argparse
import re
import subprocess
import sys

NUM_PROCESSED = 0


def log(*args):
    print(NUM_PROCESSED, *args, file=sys.stderr)


def run_script(script_str):
    """
    Run `script_str` as a Python script and return the string results
    """
    p = subprocess.run(
        ["python3", "-"],
        stdout=subprocess.PIPE,
        stderr=subprocess.STDOUT,
        input=script_str,
        text=True,
        check=True,
    )

    return p.stdout


def process_text(text):
    """
    Find multiline scripts, tagged as python, within a markdown document.
    If a script ends with the line `#eval`, the complete script will be
    executed and the string results of the execution will be appended
    to the markdown script section as a Python comment prefixed with `#`.
    """

    def _replacer(match):
        global NUM_PROCESSED
        NUM_PROCESSED += 1
        script = match.group(1).strip()

        try:
            # Add output of the program
            log("running script")
            res = run_script(script)
        except subprocess.CalledProcessError as e:
            if e.returncode == 1:
                res = e.output

        # Build replacement string
        out = "```python\n"

        # Add the original source
        out += script + "\n"

        # Then print the results - commented
        for line in filter(None, res.split("\n")):
            out += "# %s\n" % (line)

        out += "```"

        return out

    return re.sub(r"```python\n(.+?)^#eval$\n```", _replacer, text, flags=re.M | re.S)


if __name__ == "__main__":
    parser = argparse.ArgumentParser(
        description="Execute inline Python code in a markdown document and add a comment with the execution results."
    )
    parser.add_argument(
        "file", help="a markdown file with inline python scripts ending with #eval"
    )
    args = parser.parse_args()

    if args.file == "-":
        res = process_text(sys.stdin.read())
    else:
        with open(args.file, mode="r") as f:
            text = f.read()
            res = process_text(f.read())

    print(res)
