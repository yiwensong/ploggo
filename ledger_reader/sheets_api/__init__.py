import pickle

import googleapiclient.discovery
import pandas
from google_auth_oauthlib import flow


def upload_dataframe(
    client,
    sheet_id: str,
    dataframe: pandas.DataFrame,
) -> None:
    """Uploads the dataframe to a sheet."""
    sheet = client.spreadsheets()
    headers = [['Player'] + dataframe.columns.to_list()]
    data = dataframe.to_records().tolist()
    body = {
        "values": headers + data,
    }
    sheet.values().update(
        spreadsheetId=sheet_id,
        range="A1",
        valueInputOption="RAW",
        body=body,
    ).execute()


def get_client(conf_path: str, cred_cache: str):
    """Gets a googleapi client"""
    try:
        with open(cred_cache, "rb") as cache_stream:
            creds = pickle.load(cache_stream)
    except (FileNotFoundError, pickle.UnpicklingError):
        cred_flow = flow.InstalledAppFlow.from_client_secrets_file(
            conf_path,
            ["https://www.googleapis.com/auth/spreadsheets"],
        )
        creds = cred_flow.run_local_server(port=0)
        with open(cred_cache, "wb") as cache_stream:
            pickle.dump(creds, cache_stream)

    return googleapiclient.discovery.build(
        "sheets",
        "v4",
        credentials=creds,
    )
