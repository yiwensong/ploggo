import argparse
import re
import typing

import arrow
import pandas

import ledger_reader.sheets_api as sheets_api


# We use this for our dates, since the ledger doesn't have a date
YEAR = 2020

# How do you tell if the entry is a date
DATE_REGEX = re.compile(r"^(?P<month>\d\d?)/(?P<day>\d\d?)$")
# How do you get PL from the lines
PL_REGEX = re.compile(
    r"^(?P<name>[^\(\)]*)( \(.*\))?: (?P<sign>-?)\$(?P<amount>[\d.]+)$"
)


def get_day_to_day_results(
    input_stream: typing.IO,
) -> typing.Dict[str, typing.Dict[str, float]]:
    """Gets the day-to-day results from the ledger"""
    results_dict: typing.Dict[str, typing.Dict[str, float]] = dict()
    date = ""

    line = input_stream.readline()
    while line:
        is_date = DATE_REGEX.match(line)

        if is_date is not None:
            # Next day
            date = arrow.get(
                f"{YEAR}-{is_date.group('month')}-{is_date.group('day')}",
            ).isoformat()[:10]
        else:
            match = PL_REGEX.match(line)
            if match is not None:
                name = match.group("name")
                sign = -1 if match.group("sign") == "-" else 1
                amount = float(match.group("amount")) * sign

                today_result = results_dict.setdefault(date, dict())
                today_result[name] = amount

        line = input_stream.readline()

    return results_dict


def get_results(path: str) -> None:
    with open(path, "r") as result_file:
        day_to_day_results = get_day_to_day_results(result_file)

    results_df = pandas.DataFrame(day_to_day_results).fillna(0)
    return results_df


def parse_arguments() -> argparse.Namespace:
    """Parses arguments and returns the argument namespace."""
    parser = argparse.ArgumentParser(
        description="reads results and spits out some numbers"
    )
    parser.add_argument(
        "result_file", type=str, help="Path to the file containing the results",
    )

    parser.add_argument(
        "--verbose", action="store_true", help="Print the results",
    )

    parser.add_argument(
        "--output", type=str, nargs="?", default=None, help="Put the results here",
    )

    parser.add_argument(
        "--sheet",
        type=str,
        nargs="?",
        default="1KphamNq3KEq1AJe1pPsHKa-kel293mz1RyD_ud7Lg8o",
        help="Put the results in a spreadsheet",
    )

    parser.add_argument(
        "--google-credentials",
        type=str,
        nargs="?",
        default="/Users/yiwen/Downloads/credentials.json",
        help="Where the google credentials live",
    )

    parser.add_argument(
        "--google-credential-cache",
        type=str,
        nargs="?",
        default="/tmp/ledger_reader_cred_cache.pickle",
        help="Where to store the credentials cache",
    )

    parser.add_argument(
        "--disable-sheets-upload",
        action="store_true",
        help="Disables uploading to google sheets",
    )

    args = parser.parse_args()

    return args


def main() -> None:
    """Runs ledger_reader"""
    args = parse_arguments()

    results = get_results(args.result_file)
    if args.verbose:
        print(results)

    if args.output:
        with open(args.output, "w") as output_file:
            results.to_csv(output_file)

    if not args.disable_sheets_upload:
        client = sheets_api.get_client(
            args.google_credentials, args.google_credential_cache,
        )
        sheets_api.upload_dataframe(client, args.sheet, results)


if __name__ == "__main__":
    main()
