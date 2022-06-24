import argparse
import typing


def convert_token_line_to_string(tokens: typing.List[str]) -> str:
    """Converts a list of tokens to an emoji string"""
    return "".join([
        f":{token}:"
        for token in tokens
    ])


def exefy(emojis: str) -> str:
    """Exefies the emojis, input as a string"""
    tokens = emojis.split(":")
    nonempty_tokens = [t for t in tokens if t]
    token_lines = [
        nonempty_tokens[i:] + nonempty_tokens[:i]
        for i in range(len(nonempty_tokens))
    ]
    return "\n".join([
        convert_token_line_to_string(token_line)
        for token_line in token_lines
    ])


def parse_arguments() -> argparse.Namespace:
    """Parses arguments and returns the argument namespace."""
    parser = argparse.ArgumentParser(
        description="Turns a big slack emoji into an exe meme"
    )
    parser.add_argument(
        "emojis",
        type=str,
        help="the emoji that needs to be exefied, in a line",
    )

    args = parser.parse_args()

    return args


def main() -> None:
    """Runs line_calc"""
    args = parse_arguments()

    output = exefy(args.emojis)
    print(output)


if __name__ == "__main__":
    main()
