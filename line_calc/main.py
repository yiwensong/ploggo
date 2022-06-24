import argparse
import json
import statistics
import typing

from attr import asdict
from attr import define

@define
class LineSummary:
    bid: float
    ask: float
    mid: float
    spread: float

    def __str__(self):
        self_dict = asdict(self)
        return json.dumps(self_dict, indent=4)


def get_percentage_line(team1dec: float, team2dec: float) -> float:
    """Returns a summary of what the line is"""
    ask = 1.0/team1dec
    bid = 1.0 - 1.0/team2dec
    mid = statistics.mean([bid, ask])
    spread = ask - bid
    return LineSummary(bid=bid, ask=ask, mid=mid, spread=spread)


def parse_arguments() -> argparse.Namespace:
    """Parses arguments and returns the argument namespace."""
    parser = argparse.ArgumentParser(
        description="reads results and spits out some numbers"
    )
    parser.add_argument(
        "decimal1", type=float, help="decimal line of team 1",
    )
    parser.add_argument(
        "decimal2", type=float, help="decimal line of team 2",
    )

    args = parser.parse_args()

    return args


def main() -> None:
    """Runs line_calc"""
    args = parse_arguments()

    result = get_percentage_line(args.decimal1, args.decimal2)
    print(result)


if __name__ == "__main__":
    main()
