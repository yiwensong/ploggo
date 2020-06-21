import pickle
import typing

import googleapiclient.discovery
import pandas
from google_auth_oauthlib import flow


def upload_dataframe(client, sheet_id: str, dataframe: pandas.DataFrame,) -> None:
    """Uploads the dataframe to a sheet."""
    sheet = client.spreadsheets()
    headers = [["Player"] + dataframe.columns.to_list()]
    data = dataframe.to_records().tolist()
    body = {
        "values": headers + data,
    }
    sheet.values().update(
        spreadsheetId=sheet_id, range="A1", valueInputOption="RAW", body=body,
    ).execute()


def download_ledger(client, ledger_doc_id: str) -> typing.IO:
    """Creates a stream with the ledger contents."""
    doc = client.documents()
    ledger_doc = doc.get(documentId=ledger_doc_id).execute()

    return ledger_doc.get("body")


def _make_new_client(conf_path: str, cred_cache: str):
    cred_flow = flow.InstalledAppFlow.from_client_secrets_file(
        conf_path,
        [
            "https://www.googleapis.com/auth/spreadsheets",
            "https://www.googleapis.com/auth/documents",
        ],
    )
    creds = cred_flow.run_local_server(port=0)
    with open(cred_cache, "wb") as cache_stream:
        pickle.dump(creds, cache_stream)
    return creds


def get_client(conf_path: str, cred_cache: str, refresh_cache: bool = False):
    """Gets a googleapi client"""
    if refresh_cache:
        creds = _make_new_client(conf_path, cred_cache)
    else:
        try:
            with open(cred_cache, "rb") as cache_stream:
                creds = pickle.load(cache_stream)
        except (FileNotFoundError, pickle.UnpicklingError):
            creds = _make_new_client(conf_path, cred_cache)

    return googleapiclient.discovery.build("sheets", "v4", credentials=creds,)
