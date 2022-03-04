import argparse
import statistics
import typing

def get_percentage_line(team1dec: float, team2dec: float) -> float:
    """Returns the percentage chance that team 1 wins"""
    perc1 = 1.0/team1dec
    perc2 = 1.0/team2dec
    return statistics.mean([perc1, 1.0 - perc2])

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
    """Runs ledger_reader"""
    args = parse_arguments()

    result = get_percentage_line(args.decimal1, args.decimal2)
    print(result)


if __name__ == "__main__":
    main()
