import sys
from unittest import mock

import pandas

import pytest
from ledger_reader import sheets_api


def test_upload_dataframe():
    client = mock.MagicMock()

    sheet_id = "abcdefg"
    dataframe = pandas.DataFrame({
        "2020-06-13": {
            "yiwen": 1,
            "not yiwen": 2,
        },
        "2020-06-14": {
            "yiwen": 3,
            "not yiwen": 1,
        },
    })

    sheets_api.upload_dataframe(
        client,
        sheet_id,
        dataframe,
    )

    client.spreadsheets.assert_called_with()
    client.spreadsheets().values.assert_called_with()
    client.spreadsheets().values().update.assert_called_with(
        spreadsheetId=sheet_id,
        range="A1",
        valueInputOption="RAW",
        body={
            "values": [
                ["Player", "2020-06-13", "2020-06-14"],
                ("yiwen", 1, 3),
                ("not yiwen", 2, 1),
            ],
        },
    )


@mock.patch("google_auth_oauthlib.flow.InstalledAppFlow")
@mock.patch("googleapiclient.discovery.build")
@mock.patch("pickle.dump")
@mock.patch("pickle.load")
def test_get_client(pickle_load, pickle_dump, build, appflow, tmpdir):
    creds = tmpdir.join("cred.json")
    creds.write("FAKE CREDENTIAL")
    cache = tmpdir.join("does-not-exist")

    sheets_api.get_client(str(creds), str(cache))

    pickle_load.assert_not_called()
    pickle_dump.assert_called()
    appflow.from_client_secrets_file.assert_called()
    build.assert_called()


@mock.patch("google_auth_oauthlib.flow.InstalledAppFlow")
@mock.patch("googleapiclient.discovery.build")
@mock.patch("pickle.dump")
@mock.patch("pickle.load")
def test_get_client_with_cache(
    pickle_load,
    pickle_dump,
    build,
    appflow,
    tmpdir,
):
    creds = tmpdir.join("cred.json")
    creds.write("FAKE CREDENTIAL")
    cache = tmpdir.join("exists.pickle")
    cache.write("FAKE CACHE")

    sheets_api.get_client(str(creds), str(cache))

    pickle_load.assert_called()
    pickle_dump.assert_not_called()
    appflow.from_client_secrets_file.assert_not_called()
    build.assert_called()


if __name__ == "__main__":
    sys.exit(pytest.main([__file__]))
